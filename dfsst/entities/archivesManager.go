// Package entities contains database entities and helpers on these entities.
package entities

import (
	cAPI "dfss/dfssc/api"
	"dfss/mgdb"

	"gopkg.in/mgo.v2/bson"
)

// ArchivesManager : handles the structure of a SignatureArchives, with functions suited for the TTP resolve protocol.
type ArchivesManager struct {
	DB       *mgdb.MongoManager
	Archives *SignatureArchives
}

// NewArchivesManager : create a new archivesManager, with the specified mgdb manager, but
// doesn't initialize the signatureArchives. (see function 'InitializeArchives').
func NewArchivesManager(db *mgdb.MongoManager) *ArchivesManager {
	return &ArchivesManager{
		DB: db,
	}
}

// InitializeArchives : if an entry in the database for this signature exists, retrieves it, otherwise creates it.
//
// This function should only be called after function IsRequestValid.
func (manager *ArchivesManager) InitializeArchives(promise *cAPI.Promise, signatureUUID bson.ObjectId, signers *[]Signer) {
	present, archives := manager.ContainsSignature(signatureUUID)

	if !present {
		archives = NewSignatureArchives(signatureUUID, promise.Context.Sequence, *signers, promise.Context.ContractDocumentHash, promise.Context.Seal)
	}

	manager.Archives = archives
}

// ContainsSignature : checks if the specified signatureUUID matches a SignatureArchives in the database.
// If it exists, returns it.
func (manager *ArchivesManager) ContainsSignature(signatureUUID bson.ObjectId) (present bool, archives *SignatureArchives) {
	err := manager.DB.Get("signatures").FindByID(SignatureArchives{ID: signatureUUID}, &archives)
	if err != nil {
		present = false
		archives = &SignatureArchives{}
		return
	}

	present = true
	return
}

// HasReceivedAbortToken : determines if the specified signer has already been sent an abort token in the specified signatureArchives.
func (manager *ArchivesManager) HasReceivedAbortToken(signerIndex uint32) bool {
	for _, s := range manager.Archives.AbortedSigners {
		if signerIndex == s.SignerIndex {
			return true
		}
	}

	return false
}

// WasContractSigned : determines if the ttp already generated a signed contract for this signatureArchives, and returns it if true.
func (manager *ArchivesManager) WasContractSigned() (bool, []byte) {
	signedContract := manager.Archives.SignedContract
	if len(signedContract) != 0 {
		return true, signedContract
	}

	return false, []byte{}
}

// HasSignerPromised : determines if the specified signer has promised to sign to at least one other signer.
func (manager *ArchivesManager) HasSignerPromised(signer uint32) bool {
	for _, p := range manager.Archives.ReceivedPromises {
		if (p.SenderKeyIndex == signer) && (p.RecipientKeyIndex != signer) {
			return true
		}
	}

	return false
}

// AddToAbort : adds the specified signer to the aborted signers of the signatureArchives.
// If the signer is already present, does nothing.
func (manager *ArchivesManager) AddToAbort(signerIndex uint32) {
	for _, s := range manager.Archives.AbortedSigners {
		if s.SignerIndex == signerIndex {
			return
		}
	}

	// TODO
	// This requires the implementation of promises
	var abortIndex uint32
	abortIndex = 0
	abortedSigner := NewAbortedSigner(signerIndex, abortIndex)

	manager.Archives.AbortedSigners = append(manager.Archives.AbortedSigners, *abortedSigner)
}

// AddToDishonest : adds the specified signer to the dishonest signers of the signatureArchives.
// If the signer is already present, does nothing.
func (manager *ArchivesManager) AddToDishonest(signerIndex uint32) {
	for _, s := range manager.Archives.DishonestSigners {
		if s == signerIndex {
			return
		}
	}

	manager.Archives.DishonestSigners = append(manager.Archives.DishonestSigners, signerIndex)
}

// AddPromise : adds the specified promises to the list of received promises of the SignatureArchives.
func (manager *ArchivesManager) AddPromise(promise *Promise) {
	for _, p := range manager.Archives.ReceivedPromises {
		if (&p).Equal(promise) {
			return
		}
	}

	manager.Archives.ReceivedPromises = append(manager.Archives.ReceivedPromises, *promise)
}
