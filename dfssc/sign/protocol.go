package sign

import (
	"fmt"
	"log"

	cAPI "dfss/dfssc/api"
	"dfss/dfssc/common"
	dAPI "dfss/dfssd/api"
)

var (
	hashToID           map[string]uint32
	incomingPromises   chan *cAPI.Promise
	incomingSignatures chan *cAPI.Signature
)

// Sign performe all the message exchange for the contract to be signed
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

	hashToID = m.makeSignersHashToIDMap()

	incomingPromises = make(chan *cAPI.Promise)
	incomingSignatures = make(chan *cAPI.Signature)

	// Promess rounds
	// Follow the sequence until there is no next occurence of me
	for m.currentIndex >= 0 {

		dAPI.DLog("Starting round at index [" + fmt.Sprintf("%d", m.currentIndex) + "]")
		log.Println("Starting round at index [" + fmt.Sprintf("%d", m.currentIndex) + "]")

		// Set of the promise we are waiting for
		pendingSet, err1 := common.GetPendingSet(m.sequence, myID, m.currentIndex)
		if err1 != nil {
			return err1 // err is renamed to avoid shadowing err on linter check
		}

		// Set of the promises we must send
		sendSet, err1 := common.GetSendSet(m.sequence, myID, m.currentIndex)
		if err1 != nil {
			return err1
		}

		// Exchange messages
		m.promiseRound(pendingSet, sendSet, myID)

		m.currentIndex = nextIndex
		nextIndex, err1 = common.FindNextIndex(m.sequence, myID, m.currentIndex)
		if err1 != nil {
			return err1
		}
	}

	// Signature round
	err = m.SendAllSigns()
	if err != nil {
		return err
	}
	err = m.RecieveAllSigns()
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

// makeEMailMap build an association to reverse a hash to the sequence ID
func (m *SignatureManager) makeSignersHashToIDMap() map[string]uint32 {

	hashToID := make(map[string]uint32)

	signers := m.contract.Signers
	for id, signer := range signers {
		hashToID[signer.Hash] = uint32(id)
	}
	return hashToID
}

func (m *SignatureManager) promiseRound(pendingSet, sendSet []uint32, myID uint32) {

	// Reception of the due promises
	// TODO this ctx needs a timeout !
	for len(pendingSet) > 0 {
		promise := <-incomingPromises
		senderID, exist := hashToID[fmt.Sprintf("%x", promise.SenderKeyHash)]
		if exist {
			var err error
			pendingSet, err = common.Remove(pendingSet, senderID)
			if err != nil {
				// Recieve unexpected promise, ignore ?
			}
			dAPI.DLog("Recieved promise from [" + fmt.Sprintf("%d", senderID) + "]")
			m.archives.recievedPromises = append(m.archives.recievedPromises, promise)
		} else {
			// Wrong sender keyHash
		}
	}

	c := make(chan *cAPI.Promise)
	// Sending of due promises
	for _, id := range sendSet {
		// The signature manager is read only - safe !
		go func(id uint32, m *SignatureManager) {
			promise, err := m.CreatePromise(myID, id)
			if err != nil {
				// Failed to create promise
			}
			_, err = m.SendPromise(promise, id)
			if err != nil {
				// We don't check if the promise has been recieved
			}
			c <- promise
		}(id, m)
	}

	// Verifying we sent all the due promesses
	for _ = range sendSet {
		promise := <-c
		if promise != nil {
			m.archives.sentPromises = append(m.archives.sentPromises, promise)
		} else {
			// something appened during the goroutine
		}
	}
}
