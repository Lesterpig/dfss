package sign

import (
	"errors"
	"fmt"
	"time"

	cAPI "dfss/dfssc/api"
	"dfss/dfssc/common"
	dAPI "dfss/dfssd/api"
	pAPI "dfss/dfssp/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Sign makes the SignatureManager perform its specified signature
func (m *SignatureManager) Sign() error {
	myID, currentIndex, nextIndex, err := m.Initialize()
	if err != nil {
		return err
	}

	// Promess rounds
	for nextIndex != -1 {
		pendingSet, err1 := common.GetPendingSet(m.sequence, myID, currentIndex)
		if err1 != nil {
			return err1 // err is renamed to avoid shadowing err on linter check
		}

		sendSet, err1 := common.GetSendSet(m.sequence, myID, currentIndex)
		if err1 != nil {
			return err1
		}

		// Reception of the due promesses
		for len(pendingSet) != 0 {
			i := 0
			// TODO
			// Improve, because potential memory leak
			// see https://github.com/golang/go/wiki/SliceTricks
			pendingSet = append(pendingSet[:i], pendingSet[i+1:]...)
		}

		c := make(chan int)
		// Sending of the due promesses
		/*
		   for _, id := range sendSet {
		       go func(id) {
		           promise, err := m.CreatePromise(id)
		           recpt := m.SendPromise(promise, id)
		           c <- id
		       }(id)
		   }
		*/

		// Verifying we sent all the due promesses
		for _ = range sendSet {
			<-c
		}

		currentIndex = nextIndex
		nextIndex, err1 = common.FindNextIndex(m.sequence, myID, currentIndex)
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

// CreatePromise creates a promise from 'from' to 'to', in the context of the SignatureManager
// provided the specified sequence indexes are valid
func (m *SignatureManager) CreatePromise(from, to uint32) (*cAPI.Promise, error) {
	if int(from) >= len(m.keyHash) || int(to) >= len(m.keyHash) {
		return &cAPI.Promise{}, errors.New("Invalid id for promise creation")
	}
	promise := &cAPI.Promise{
		RecipientKeyHash: m.keyHash[to],
		SenderKeyHash:    m.keyHash[from],
		SignatureUuid:    m.uuid,
		ContractUuid:     m.contract.UUID,
	}
	return promise, nil
}

// SendPromise sends the specified promise to the specified peer
// TODO
func (m *SignatureManager) SendPromise(promise *cAPI.Promise, to uint32) (*pAPI.ErrorCode, error) {
	connection, err := m.GetClient(to)
	if err != nil {
		return &pAPI.ErrorCode{}, err
	}

	// Handle the timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	errCode, err := (*connection).TreatPromise(ctx, promise)
	if err == grpc.ErrClientConnTimeout {
		dAPI.DLog("Promise timeout for [" + fmt.Sprintf("%d", to) + "]")
		return &pAPI.ErrorCode{Code: pAPI.ErrorCode_TIMEOUT, Message: "promise timeout"}, err
	} else if err != nil {
		return &pAPI.ErrorCode{Code: pAPI.ErrorCode_INTERR, Message: "internal server error"}, err
	}

	m.archives.sentPromises = append(m.archives.sentPromises, promise)

	return errCode, nil
}

// GetClient retrieves the Client to the specified sequence id provided it exists
func (m *SignatureManager) GetClient(to uint32) (*cAPI.ClientClient, error) {
	mailto := m.contract.Signers[to].Email

	if _, ok := m.peers[mailto]; !ok {
		return nil, fmt.Errorf("No connection to user %s", mailto)
	}

	return m.peers[mailto], nil
}

// SendAllSigns creates and sends signatures to all the signers of the contract
// TODO
// Use goroutines to send in parallel
func (m *SignatureManager) SendAllSigns() error {
	myID, err := m.FindID()
	if err != nil {
		return err
	}

	sendSet := common.GetAllButOne(m.sequence, myID)

	for _, id := range sendSet {
		signature, err := m.CreateSignature(myID, id)
		if err != nil {
			return err
		}

		_, err = m.SendSignature(signature, id)
		if err != nil {
			return err
		}
	}

	return nil
}

// CreateSignature creates a signature from from to to, in the context of the SignatureManager
// provided the specified sequence indexes are valid
// TODO
// Implement a true cryptographic signature
func (m *SignatureManager) CreateSignature(from, to uint32) (*cAPI.Signature, error) {
	if int(from) >= len(m.keyHash) || int(to) >= len(m.keyHash) {
		return &cAPI.Signature{}, errors.New("Invalid id for signature creation")
	}
	signature := &cAPI.Signature{
		RecipientKeyHash: m.keyHash[to],
		SenderKeyHash:    m.keyHash[from],
		Signature:        "Signature",
		SignatureUuid:    m.uuid,
		ContractUuid:     m.contract.UUID,
	}
	return signature, nil
}

// SendSignature sends the specified signature to the specified peer
// TODO
func (m *SignatureManager) SendSignature(signature *cAPI.Signature, to uint32) (*pAPI.ErrorCode, error) {
	connection, err := m.GetClient(to)
	if err != nil {
		return &pAPI.ErrorCode{}, err
	}

	// TODO
	// Handle the timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	errCode, err := (*connection).TreatSignature(ctx, signature)
	if err != nil {
		return &pAPI.ErrorCode{}, err
	}

	m.archives.sentSignatures = append(m.archives.sentSignatures, signature)

	return errCode, nil
}

// RecieveAllSigns is not done yet
// TODO
func (m *SignatureManager) RecieveAllSigns() error {
	return nil
}
