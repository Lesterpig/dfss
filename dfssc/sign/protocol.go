// Package sign handles contract and signature operations.
package sign

import (
	"fmt"
	"time"

	cAPI "dfss/dfssc/api"
	"dfss/dfssc/common"
	dAPI "dfss/dfssd/api"
	"github.com/spf13/viper"
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

	myID, nextIndex, err := m.Initialize()
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
	for m.currentIndex >= 0 {
		m.OnProgressUpdate(m.currentIndex, seqLen+1)
		time.Sleep(viper.GetDuration("slowdown"))
		dAPI.DLog("starting round at index [" + fmt.Sprintf("%d", m.currentIndex) + "] with nextIndex=" + fmt.Sprintf("%d", nextIndex))

		// Set of promises we are waiting for
		var pendingSet []uint32
		pendingSet, err = common.GetPendingSet(m.sequence, myID, m.currentIndex)
		if err != nil {
			return err
		}

		// Set of the promises we must send
		var sendSet []uint32
		sendSet, err = common.GetSendSet(m.sequence, myID, m.currentIndex)
		if err != nil {
			return err
		}

		// Exchange messages
		m.promiseRound(pendingSet, sendSet, myID)

		m.currentIndex = nextIndex
		nextIndex, err = common.FindNextIndex(m.sequence, myID, m.currentIndex)
		if err != nil {
			return err
		}
	}

	m.OnProgressUpdate(seqLen, seqLen+1)
	dAPI.DLog("entering signature round")
	// Signature round
	err = m.ExchangeAllSignatures()
	if err != nil {
		return err
	}
	dAPI.DLog("exiting signature round")
	m.OnProgressUpdate(seqLen+1, seqLen+1)

	return nil
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
func (m *SignatureManager) promiseRound(pendingSet, sendSet []uint32, myID uint32) {

	// Reception of the due promises
	for len(pendingSet) > 0 {
		select {
		case promiseIface := <-m.cServerIface.incomingPromises:
			promise := (promiseIface).(*cAPI.Promise)
			senderID, exist := m.hashToID[fmt.Sprintf("%x", promise.Context.SenderKeyHash)]
			if exist {
				var err error
				pendingSet, err = common.Remove(pendingSet, senderID)
				if err != nil {
					continue
				}
				m.archives.receivedPromises = append(m.archives.receivedPromises, promise)
			}

		case <-time.After(time.Minute):
			// TODO contact TTP
			return
		}
	}

	c := make(chan *cAPI.Promise, chanBufferSize)
	// Sending of due promises
	for _, id := range sendSet {
		go func(id uint32, m *SignatureManager) {
			promise, err := m.CreatePromise(myID, id)
			if err == nil {
				_ = m.SendEvidence(promise, nil, id)
			}
			c <- promise
		}(id, m)
	}

	// Verifying we sent all the due promises
	for range sendSet {
		<-c
	}
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
