package sign

import (
	"errors"
	"fmt"
	"time"

	cAPI "dfss/dfssc/api"
	dAPI "dfss/dfssd/api"
	pAPI "dfss/dfssp/api"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// CreatePromise creates a promise from 'from' to 'to', in the context of the SignatureManager
// provided the specified sequence indexes are valid
func (m *SignatureManager) CreatePromise(from, to uint32) (*cAPI.Promise, error) {
	if int(from) >= len(m.keyHash) || int(to) >= len(m.keyHash) {
		return nil, errors.New("Invalid id for promise creation")
	}
	if m.currentIndex < 0 {
		return nil, errors.New("Invalid currentIndex for promise creation")
	}
	promise := &cAPI.Promise{
		RecipientKeyHash:     m.keyHash[to],
		SenderKeyHash:        m.keyHash[from],
		Index:                uint32(m.currentIndex),
		ContractDocumentHash: m.contract.File.Hash,
		SignatureUuid:        m.uuid,
		ContractUuid:         m.contract.UUID,
	}
	return promise, nil
}

// SendPromise sends the specified promise to the specified peer
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
