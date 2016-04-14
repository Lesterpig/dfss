package sign

import (
	"fmt"
	"log"

	cAPI "dfss/dfssc/api"
	"dfss/dfssc/common"
	dAPI "dfss/dfssd/api"
)

// Sign perform all the message exchanges for the contract to be signed
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

	m.cServerIface.incomingPromises = make(chan *cAPI.Promise)
	m.cServerIface.incomingSignatures = make(chan *cAPI.Signature)

	// Promess rounds
	// Follow the sequence until there is no next occurence of me
	for m.currentIndex >= 0 {

		dAPI.DLog(fmt.Sprintf("{%d} ", myID) + "Starting round at index [" + fmt.Sprintf("%d", m.currentIndex) + "] with nextIndex=" + fmt.Sprintf("%d", nextIndex))

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

	dAPI.DLog("{" + fmt.Sprintf("%d", myID) + "} Enter signature round")

	// Signature round
	err = m.ExchangeAllSignatures()
	if err != nil {
		return err
	}

	dAPI.DLog("{" + fmt.Sprintf("%d", myID) + "} Exit signature round")

	// Network's job is done, cleaning time
	// Shutdown and platform client and TODO peer server & connections
	err = m.platformConn.Close()
	if err != nil {
		return err
	}

	return nil
}

// GetClient retrieves the Client to the specified sequence id provided it exists
func (m *SignatureManager) GetClient(to uint32) (*cAPI.ClientClient, error) {
	mailto := m.contract.Signers[to].Email

	if _, ok := m.peers[mailto]; !ok {
		return nil, fmt.Errorf("No connection to user %s", mailto)
	}

	return m.peers[mailto], nil
}

// makeSignersHashToIDMap build an association to reverse a hash to the sequence ID
func (m *SignatureManager) makeSignersHashToIDMap() {

	m.hashToID = make(map[string]uint32)

	signers := m.contract.Signers
	for id, signer := range signers {
		m.hashToID[signer.Hash] = uint32(id)
	}
}

// promiseRound describe a promise round: reception and sending
//
// TODO better error management - this function should return `error` !
func (m *SignatureManager) promiseRound(pendingSet, sendSet []uint32, myID uint32) {

	// Reception of the due promises
	// TODO this ctx needs a timeout !
	for len(pendingSet) > 0 {
		promise := <-m.cServerIface.incomingPromises
		senderID, exist := m.hashToID[fmt.Sprintf("%x", promise.SenderKeyHash)]
		if exist {
			var err error
			pendingSet, err = common.Remove(pendingSet, senderID)
			if err != nil {
				_ = fmt.Errorf("Receive unexpected promise")
			}
			m.archives.receivedPromises = append(m.archives.receivedPromises, promise)
			dAPI.DLog("{" + fmt.Sprintf("%d", myID) + "} Received promise from [" + fmt.Sprintf("%d", senderID) + "] for index " + fmt.Sprintf("%d", promise.Index))
		} else {
			// Wrong sender keyHash
			log.Println("{" + fmt.Sprintf("%d", myID) + "} Wrong sender keyhash !")
		}
	}

	c := make(chan *cAPI.Promise)
	// Sending of due promises
	for _, id := range sendSet {
		// The signature manager is read only - safe !
		go func(id uint32, m *SignatureManager) {
			promise, err := m.CreatePromise(myID, id)
			if err != nil {
				_ = fmt.Errorf("Failed to create promise from %d to %d", myID, id)
			}
			dAPI.DLog("{" + fmt.Sprintf("%d", myID) + "} Send promise to " + fmt.Sprintf("%d", id))
			_, err = m.SendPromise(promise, id)
			if err != nil {
				dAPI.DLog("{" + fmt.Sprintf("%d", myID) + "} Promise have not been received !")
				_ = fmt.Errorf("Failed to deliver promise from %d to %d", myID, id)
			}
			c <- promise
		}(id, m)
	}

	// Verifying we sent all the due promises
	for _ = range sendSet {
		promise := <-c
		if promise != nil {
			m.archives.sentPromises = append(m.archives.sentPromises, promise)
		} else {
			// something appened during the goroutine
		}
	}
}

// closeAllPeerClient tries to close all established connection with other peers
func (m *SignatureManager) closeAllPeerClient() {
	for k, client := range m.peersConn {
		_ = client.Close()
		// Remove associated grpc client
		delete(m.peers, k)
		fmt.Println("- Close connection to " + k)
	}
}
