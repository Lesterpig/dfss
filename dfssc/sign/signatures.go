package sign

import (
	"dfss/dfssc/common"
	"fmt"

	cAPI "dfss/dfssc/api"
)

// ExchangeAllSignatures creates and sends signatures to all the signers of the contract
func (m *SignatureManager) ExchangeAllSignatures() error {
	allReceived := make(chan error)
	go m.ReceiveAllSignatures(allReceived)

	myID, err := m.FindID()
	if err != nil {
		return err
	}

	// compute a set of all signers except me
	sendSet := common.GetAllButOne(m.sequence, myID)
	errorChan := make(chan error)
	for _, id := range sendSet {
		go func(id uint32) {
			signature, err2 := m.CreateSignature(myID, id)
			if err2 != nil {
				errorChan <- err2
				return
			}
			err2 = m.SendEvidence(nil, signature, id)
			if err2 != nil {
				errorChan <- err2
				return
			}
			errorChan <- nil
		}(id)
	}

	for range sendSet {
		err = <-errorChan
		if err != nil {
			return err
		}
	}

	return <-allReceived
}

// CreateSignature creates a signature from a sequence ID to another
// provided the specified sequence indexes are valid
func (m *SignatureManager) CreateSignature(from, to uint32) (*cAPI.Signature, error) {
	context, err := m.createContext(from, to)
	if err != nil {
		return nil, err
	}

	return &cAPI.Signature{
		Context: context,
		Payload: []byte{0x42},
	}, nil
}

// ReceiveAllSignatures receive all the signatures
func (m *SignatureManager) ReceiveAllSignatures(out chan error) {
	myID, err := m.FindID()
	if err != nil {
		out <- err
		return
	}

	// compute a set of all signers except me
	pendingSet := common.GetAllButOne(m.sequence, myID)

	// TODO this ctx needs a timeout !
	for len(pendingSet) > 0 {
		signature := (<-m.cServerIface.incomingSignatures).(*cAPI.Signature)
		senderID, exist := m.hashToID[fmt.Sprintf("%x", signature.Context.SenderKeyHash)]
		if exist {
			pendingSet, _ = common.Remove(pendingSet, senderID)
			m.archives.receivedSignatures = append(m.archives.receivedSignatures, signature)
		}
	}

	out <- nil
	return
}
