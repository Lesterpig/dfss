package sign

import (
	"fmt"
	"time"

	cAPI "dfss/dfssc/api"
	"dfss/dfssc/common"
	dAPI "dfss/dfssd/api"
)

// Sign performs all the message exchanges for the contract to be signed
//
// * Initialize the SignatureManager from starter.go
// * Compute the reversed map [mail -> ID] of signers
// * Make channels for handlers
// * Promises rounds
// * Signature round
func (m *SignatureManager) Sign() error {
	myID, nextIndex, err := m.Initialize()
	if err != nil {
		return err
	}

	m.makeSignersHashToIDMap()
	m.cServerIface.incomingPromises = make(chan interface{})
	m.cServerIface.incomingSignatures = make(chan interface{})

	// Cooldown delay, let other clients wake-up their channels
	time.Sleep(time.Second)

	// Promess rounds
	// Follow the sequence until there is no next occurence of me
	for m.currentIndex >= 0 {
		dAPI.DLog("starting round at index [" + fmt.Sprintf("%d", m.currentIndex) + "] with nextIndex=" + fmt.Sprintf("%d", nextIndex))

		// Set of the promise we are waiting for
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

	dAPI.DLog("entering signature round")
	// Signature round
	err = m.ExchangeAllSignatures()
	if err != nil {
		return err
	}
	dAPI.DLog("exiting signature round")

	// Network's job is done, cleaning time
	// Shutdown and platform client and TODO peer server & connections
	return m.platformConn.Close()
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
// TODO better error management - this function should return `error` !
func (m *SignatureManager) promiseRound(pendingSet, sendSet []uint32, myID uint32) {

	// Reception of the due promises
	// TODO this ctx needs a timeout !
	for len(pendingSet) > 0 {
		promise := (<-m.cServerIface.incomingPromises).(*cAPI.Promise)
		senderID, exist := m.hashToID[fmt.Sprintf("%x", promise.Context.SenderKeyHash)]
		if exist {
			var err error
			pendingSet, err = common.Remove(pendingSet, senderID)
			if err != nil {
				continue
			}
			m.archives.receivedPromises = append(m.archives.receivedPromises, promise)
		}
	}

	c := make(chan *cAPI.Promise)
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
		_ = <-c
	}
}

// closeAllPeerClient tries to close all established connection with other peers
func (m *SignatureManager) closeAllPeerClient() {
	for k, client := range m.peersConn {
		_ = client.Close()
		// Remove associated grpc client
		delete(m.peers, k)
	}
}
