// Package sign handles contract and signature operations.
package sign

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	cAPI "dfss/dfssc/api"
	"dfss/dfssc/common"
	dAPI "dfss/dfssd/api"
	tAPI "dfss/dfsst/api"
	"dfss/net"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Sign performs all the message exchanges for the contract to be signed
//
// * Initialize the SignatureManager from starter.go
// * Compute the reversed map [mail -> ID] of signers
// * Make channels for handlers
// * Promises rounds
// * Signature round
func (m *SignatureManager) Sign() error {

	defer func() {
		m.finished = true
		m.closeConnections()
	}()

	nextIndex, err := m.Initialize()
	if err != nil {
		return err
	}

	m.makeSignersHashToIDMap()
	m.cServerIface.incomingPromises = make(chan interface{}, chanBufferSize)
	m.cServerIface.incomingSignatures = make(chan interface{}, chanBufferSize)

	// Cooldown delay, let other clients wake-up their channels
	time.Sleep(time.Second)

	seqLen := len(m.sequence)

	// Promess rounds
	// Follow the sequence until there is no next occurence of me
	round := 0
	for m.currentIndex >= 0 {
		round = round + 1
		stopIfNeeded(round)
		m.OnProgressUpdate(m.currentIndex, seqLen+1)
		time.Sleep(viper.GetDuration("slowdown"))
		dAPI.DLog("starting round at index [" + fmt.Sprintf("%d", m.currentIndex) + "] with nextIndex=" + fmt.Sprintf("%d", nextIndex))

		// Set of promises we are waiting for
		var pendingSet []common.SequenceCoordinate
		pendingSet, err = common.GetPendingSet(m.sequence, m.myID, m.currentIndex)
		if err != nil {
			return err
		}

		// Set of the promises we must send
		var sendSet []common.SequenceCoordinate
		sendSet, err = common.GetSendSet(m.sequence, m.myID, m.currentIndex)
		if err != nil {
			return err
		}

		// Exchange messages
		var stop bool
		stop, err = m.promiseRound(pendingSet, sendSet)
		if err != nil || stop {
			dAPI.DLog("stopping protocol execution")
			return err
		}

		m.currentIndex = nextIndex
		nextIndex, err = common.FindNextIndex(m.sequence, m.myID, m.currentIndex)
		if err != nil {
			return err
		}
	}

	// Signature round
	stopIfNeeded(-1)
	m.OnProgressUpdate(seqLen, seqLen+1)
	dAPI.DLog("entering signature round")
	err = m.ExchangeAllSignatures()
	if err != nil {
		return err
	}

	dAPI.DLog("exiting signature round")
	m.OnProgressUpdate(seqLen+1, seqLen+1)
	return m.PersistSignaturesToFile()
}

// GetClient retrieves the Client to the specified sequence id provided it exists
func (m *SignatureManager) GetClient(to uint32) (client *cAPI.ClientClient, mail string) {
	mail = m.contract.Signers[to].Email
	client = m.peers[mail]
	return
}

// makeSignersHashToIDMap builds an association to reverse a hash to the sequence ID
func (m *SignatureManager) makeSignersHashToIDMap() {
	m.hashToID = make(map[string]uint32)
	signers := m.contract.Signers
	for id, signer := range signers {
		m.hashToID[signer.Hash] = uint32(id)
	}
}

// promiseRound describes a promise round: reception and sending
// returns true if the client has to stop the protocol, false otherwise.
// returns an error if any occured.
func (m *SignatureManager) promiseRound(pendingSet, sendSet []common.SequenceCoordinate) (bool, error) {
	// Reception of the due promises
	var promises []*cAPI.Promise
	for len(pendingSet) > 0 {
		select {
		case promiseIface := <-m.cServerIface.incomingPromises:
			promise := (promiseIface).(*cAPI.Promise)
			valid, senderID := m.checkPromise(pendingSet, promise)
			if valid {
				var err error
				pendingSet, err = common.RemoveCoordinate(pendingSet, senderID)
				if err != nil {
					continue
				}
				promises = append(promises, promise)
			} else {
				return true, m.resolve()
			}

		case <-time.After(net.DefaultTimeout):
			return true, m.resolve()
		}
	}

	// Now that we received everything, we update the evidence we will give to the ttp
	m.updateReceivedPromises(promises)
	m.lastValidIndex = m.currentIndex

	c := make(chan error, chanBufferSize)
	// Sending of due promises
	for _, coord := range sendSet {
		go func(coord common.SequenceCoordinate, m *SignatureManager) {
			promise, err := m.CreatePromise(m.myID, coord.Signer, uint32(m.currentIndex))
			if err == nil {
				err = m.SendEvidence(promise, nil, coord.Signer)
			}
			c <- err
		}(coord, m)
	}

	// Verifying we sent all the due promises
	for range sendSet {
		v := <-c
		if v != nil {
			// We couldn't send a due promise
			dAPI.DLog("Couldn't send promise: " + v.Error())
			return true, m.resolve()
		}
	}

	return false, nil
}

// closeConnections tries to close all established connection with other peers and platform.
// It also stops the local server.
func (m *SignatureManager) closeConnections() {
	_ = m.platformConn.Close()
	for k, peer := range m.peersConn {
		_ = peer.Close()
		delete(m.peers, k)
	}
	m.cServer.Stop()
}

// updateReceivedPromises : updates the RecievedPromises field of the SignatureManager with the provided promises:
// if we don't yet have a promise from this signer, we add it to the array.
// otherwise we replace the one we have by the provided promise.
func (m *SignatureManager) updateReceivedPromises(promises []*cAPI.Promise) {
	for _, p := range promises {
		present, index := m.containsPromiseFrom(p.Context.SenderKeyHash, p.Index)

		if present {
			// it's present, so there is no index error
			_ = m.removeReceivedPromise(index)
		}
		m.archives.receivedPromises = append(m.archives.receivedPromises, p)
	}
}

// containsPromiseFrom : determines if the SignatureManager has already archived a promise from the specified signer, previous to the specified index.
func (m *SignatureManager) containsPromiseFrom(signer []byte, index uint32) (bool, int) {
	for i, p := range m.archives.receivedPromises {
		if bytes.Equal(p.Context.SenderKeyHash, signer) {
			return p.Index < index, i
		}
	}
	return false, 0
}

// removeReceivedPromise : removes the promise at the specified index from the archived received promises.
// If the index is invalid, return an error.
// If the promise is not there, does nothing.
func (m *SignatureManager) removeReceivedPromise(index int) error {
	promises := m.archives.receivedPromises
	if index < 0 || index >= len(promises) {
		return errors.New("Index out of range")
	}

	m.archives.receivedPromises = append(promises[:index], promises[index+1:]...)

	return nil
}

// callForResolve : calls the ttp for resolution.
func (m *SignatureManager) callForResolve() (*tAPI.TTPResponse, error) {
	selfPromise, err := m.CreatePromise(m.myID, m.myID, uint32(m.lastValidIndex))
	if err != nil {
		return nil, err
	}

	toSend := append(m.archives.receivedPromises, selfPromise)

	request := &tAPI.AlertRequest{Promises: toSend, Index: uint32(m.lastValidIndex)}

	ctx, cancel := context.WithTimeout(context.Background(), net.DefaultTimeout)
	defer cancel()
	response, err := m.ttp.Alert(ctx, request)
	if err != nil {
		return nil, errors.New(grpc.ErrorDesc(err))
	}

	return response, nil
}

// resolve : calls for the resolution, and persists the contract if obtained.
func (m *SignatureManager) resolve() error {
	if m.ttp == nil {
		dAPI.DLog("unable to contact TTP")
		return errors.New("No connection to TTP, aborting!")
	}

	dAPI.DLog("contacting TTP with resolve index " + fmt.Sprint(m.lastValidIndex))
	response, err := m.callForResolve()
	if err != nil {
		dAPI.DLog("Resolve call generated an error: " + err.Error())
		return err
	}
	if response.Abort {
		dAPI.DLog("contacted TTP, received abort token")
		return nil
	}
	dAPI.DLog("contacted TTP, received signed contract")
	return ioutil.WriteFile(m.mail+"-"+m.contract.UUID+".proof", response.Contract, 0600)
}

// checkPromise : verifies that the promise is valid wrt the expected promises.
// We assume that the promise data is consistent wrt the platform seal.
func (m *SignatureManager) checkPromise(expected []common.SequenceCoordinate, promise *cAPI.Promise) (bool, uint32) {
	// the promise is consistent, but not for the expected signature
	// this should not happen
	if promise.Context.SignatureUUID != m.uuid {
		return false, 0
	}

	// the promise is not for us
	recipientID, exist := m.hashToID[fmt.Sprintf("%x", promise.Context.RecipientKeyHash)]
	if !exist || recipientID != m.myID {
		return false, 0
	}

	// we didn't expect a promise from this client
	senderID, exist := m.hashToID[fmt.Sprintf("%x", promise.Context.SenderKeyHash)]
	if !exist {
		return false, 0
	}
	for _, c := range expected {
		if c.Signer == senderID && c.Index == promise.Index {
			return true, senderID
		}
	}

	return false, 0
}

func stopIfNeeded(index int) {
	s := viper.GetInt("stopbefore")
	if s == 0 {
		return
	}

	if index == -1 && s == -1 || index == s {
		os.Exit(0)
	}
}
