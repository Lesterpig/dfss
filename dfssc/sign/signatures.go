package sign

import (
	"dfss/dfssc/common"
	"errors"
	"fmt"
	"time"

	cAPI "dfss/dfssc/api"
	dAPI "dfss/dfssd/api"
	pAPI "dfss/dfssp/api"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// ExchangeAllSignatures creates and sends signatures to all the signers of the contract
func (m *SignatureManager) ExchangeAllSignatures() error {

	allReceived := make(chan error)
	go m.ReceiveAllSigns(allReceived)

	myID, err := m.FindID()
	if err != nil {
		return err
	}

	// compute a set of all signers exept me
	sendSet := common.GetAllButOne(m.sequence, myID)

	errorChan := make(chan error)

	for _, id := range sendSet {
		go func(id uint32) {
			signature, err2 := m.CreateSignature(myID, id)
			if err2 != nil {
				errorChan <- err2
				return
			}

			dAPI.DLog("{" + fmt.Sprintf("%d", myID) + "} Send sign to " + fmt.Sprintf("%d", id))
			_, err2 = m.SendSignature(signature, id)
			if err2 != nil {
				errorChan <- err2
				return
			}

			errorChan <- nil
			return
		}(id)
	}

	for range sendSet {
		err = <-errorChan
		if err != nil {
			return err
		}
	}

	err = <-allReceived
	if err != nil {
		return err
	}

	return nil
}

// CreateSignature creates a signature from a sequence ID to another
// provided the specified sequence indexes are valid
// TODO Implement a true cryptographic signature
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
func (m *SignatureManager) SendSignature(signature *cAPI.Signature, to uint32) (*pAPI.ErrorCode, error) {
	connection, err := m.GetClient(to)
	if err != nil {
		return &pAPI.ErrorCode{}, err
	}

	// Handle the timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	errCode, err := (*connection).TreatSignature(ctx, signature)
	if err == grpc.ErrClientConnTimeout {
		dAPI.DLog("Signature timeout for [" + fmt.Sprintf("%d", to) + "]")
		return &pAPI.ErrorCode{Code: pAPI.ErrorCode_TIMEOUT, Message: "signature timeout"}, err
	} else if err != nil {
		return &pAPI.ErrorCode{Code: pAPI.ErrorCode_INTERR, Message: "internal server error"}, err
	}

	m.archives.sentSignatures = append(m.archives.sentSignatures, signature)

	return errCode, nil
}

// ReceiveAllSigns receive all the signatures
func (m *SignatureManager) ReceiveAllSigns(out chan error) {
	myID, err := m.FindID()
	if err != nil {
		out <- err
		return
	}

	// compute a set of all signers except me
	pendingSet := common.GetAllButOne(m.sequence, myID)

	// TODO this ctx needs a timeout !
	for len(pendingSet) > 0 {
		signature := <-m.cServerIface.incomingSignatures
		senderID, exist := m.hashToID[fmt.Sprintf("%x", signature.SenderKeyHash)]
		dAPI.DLog("{" + fmt.Sprintf("%d", myID) + "} Receive sign from " + fmt.Sprintf("%d", senderID))
		if exist {
			var err error
			pendingSet, err = common.Remove(pendingSet, senderID)
			if err != nil {
				// Receive unexpected signature, ignore ?
			}
			m.archives.receivedSignatures = append(m.archives.receivedSignatures, signature)
		} else {
			// Wrong sender keyHash
		}
	}

	out <- nil
	return
}
