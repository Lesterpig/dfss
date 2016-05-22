package sign

import (
	"fmt"
	"time"

	cAPI "dfss/dfssc/api"
	"dfss/dfssc/common"
	"github.com/spf13/viper"
)

// ExchangeAllSignatures creates and sends signatures to all the signers of the contract
func (m *SignatureManager) ExchangeAllSignatures() error {
	allReceived := make(chan error, chanBufferSize)
	go m.ReceiveAllSignatures(allReceived)

	myID, err := m.FindID()
	if err != nil {
		return err
	}

	// compute a set of all signers except me
	sendSet := common.GetAllButOne(m.sequence, myID)
	errorChan := make(chan error, chanBufferSize)
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

	for len(pendingSet) > 0 {
		select {
		// Waiting for signatures from grpc handler
		case signatureIface := <-m.cServerIface.incomingSignatures:
			signature := (signatureIface).(*cAPI.Signature)
			senderID, exist := m.hashToID[fmt.Sprintf("%x", signature.Context.SenderKeyHash)]
			if exist {
				pendingSet, _ = common.Remove(pendingSet, senderID)
				m.archives.receivedSignatures = append(m.archives.receivedSignatures, signature)
			}

		case <-time.After(viper.GetDuration("timeout")):
			out <- fmt.Errorf("Signature reception timeout!")
			return
		}
	}

	out <- nil
	return
}
