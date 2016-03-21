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

// SendAllSigns creates and sends signatures to all the signers of the contract
// TODO Use goroutines to send in parallel
func (m *SignatureManager) SendAllSigns() error {
	myID, err := m.FindID()
	if err != nil {
		return err
	}

	// compute a set of all signers exept me
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
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

// RecieveAllSigns is not done yet
func (m *SignatureManager) RecieveAllSigns() error {
	myID, err := m.FindID()
	if err != nil {
		return err
	}

	// compute a set of all signers exept me
	pendingSet := common.GetAllButOne(m.sequence, myID)

	// TODO this ctx needs a timeout !
	for len(pendingSet) > 0 {
		signature := <-incomingSignatures
		senderID, exist := hashToID[string(signature.SenderKeyHash)]
		if exist {
			var err error
			pendingSet, err = common.Remove(pendingSet, senderID)
			if err != nil {
				// Recieve unexpected signature, ignore ?
			}
			m.archives.recievedSignatures = append(m.archives.recievedSignatures, signature)
		} else {
			// Wrong sender keyHash
		}
	}

	return nil
}
