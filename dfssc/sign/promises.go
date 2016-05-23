package sign

import (
	"encoding/hex"
	"errors"
	"time"

	cAPI "dfss/dfssc/api"
	dAPI "dfss/dfssd/api"
	pAPI "dfss/dfssp/api"
	"golang.org/x/net/context"
)

func (m *SignatureManager) createContext(from, to uint32) (*cAPI.Context, error) {
	if int(from) >= len(m.keyHash) || int(to) >= len(m.keyHash) {
		return nil, errors.New("Invalid id for context creation")
	}

	h, _ := hex.DecodeString(m.contract.File.Hash)

	return &cAPI.Context{
		RecipientKeyHash:     m.keyHash[to],
		SenderKeyHash:        m.keyHash[from],
		Sequence:             m.sequence,
		Signers:              m.keyHash,
		ContractDocumentHash: h,
		SignatureUUID:        m.uuid,
		TtpAddrPort:          m.ttpData.Addrport,
		TtpHash:              m.ttpData.Hash,
		Seal:                 m.seal,
	}, nil
}

// CreatePromise creates a promise from 'from' to 'to', in the context of the SignatureManager
// provided that the specified sequence indexes are valid
func (m *SignatureManager) CreatePromise(from, to, at uint32) (*cAPI.Promise, error) {
	context, err := m.createContext(from, to)
	if err != nil {
		return nil, err
	}

	if m.currentIndex < 0 {
		return nil, errors.New("Invalid currentIndex for promise creation")
	}

	return &cAPI.Promise{
		Index:   at,
		Context: context,
		Payload: []byte{0x41},
	}, nil
}

// SendEvidence factorizes the send code between promises and signatures.
// You can use it by setting either promise or signature to `nil`.
// The successfully sent evidence is then added to the archives.
func (m *SignatureManager) SendEvidence(promise *cAPI.Promise, signature *cAPI.Signature, to uint32) (err error) {
	connection, mail := m.GetClient(to)
	if connection == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	var result *pAPI.ErrorCode
	if promise != nil {
		result, err = (*connection).TreatPromise(ctx, promise)
	} else if signature != nil {
		result, err = (*connection).TreatSignature(ctx, signature)
	} else {
		err = errors.New("both promise and signature are nil, cannot send anything")
	}

	if err == nil && result != nil && result.Code == pAPI.ErrorCode_SUCCESS {
		m.archives.mutex.Lock()
		if promise != nil {
			dAPI.DLog("successfully sent promise to " + mail)
		} else {
			dAPI.DLog("successfully sent signature to " + mail)
			m.archives.sentSignatures = append(m.archives.sentSignatures, signature)
		}
		m.archives.mutex.Unlock()
	} else {
		dAPI.DLog("was unable to send evidence to " + mail)
		if err != nil {
			return
		}
		err = errors.New("received wrong error code")
	}

	return
}
