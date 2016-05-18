package entities

import (
	"gopkg.in/mgo.v2/bson"
)

// Signer : represents a signer for the TTP
type Signer struct {
	ID   bson.ObjectId `key:"_id" bson:"_id"`   // Internal id of a Signer
	Hash []byte        `key:"hash" bson:"hash"` // The SHA-512 hash of the signer's certificate
}

// NewSigner : creates a new Signer with the specified hash
//
// The specified hash validity is not checked in this function (see function IsValidSignerHash)
func NewSigner(hash []byte) *Signer {
	return &Signer{
		ID:   bson.NewObjectId(),
		Hash: hash,
	}
}

// SignatureArchives : represents the valid archives related to a signature process
type SignatureArchives struct {
	ID bson.ObjectId `key:"_id" bson:"_id"` // Internal id of a SignatureArchives - The unique signature identifier

	Sequence []uint32 `key:"sequence" bson:"sequence"` // Signing sequence
	Signers  []Signer `key:"signers" bson:"signers"`   // List of signers
	TextHash []byte   `key:"textHash" bson:"textHash"` // Small hash of the contract
	Seal     []byte   `key:"seal" bson:"seal"`         // Seal provided by the platform to authentify the context

	ReceivedPromises []Promise       `key:"receivedPromises" bson:"receivedPromises"` // Set of valid received promises (1 by sender)
	AbortedSigners   []AbortedSigner `key:"abortedSigners" bson:"abortedSigners"`     // Signers that were sent an abort token
	DishonestSigners []uint32        `key:"dishonestSigners" bson:"dishonestSigners"` // Indexes of the signers that were evaluated as dishonest

	SignedContract []byte `key:"signedContract" bson:"signedContract"` // Signed contract resulting of the signing process
}

// NewSignatureArchives : creates a new SignatureArchives with the specified parameters
func NewSignatureArchives(signatureUUID bson.ObjectId, sequence []uint32, signers []Signer, textHash, seal []byte) *SignatureArchives {
	return &SignatureArchives{
		ID: signatureUUID,

		Sequence: sequence,
		Signers:  signers,
		TextHash: textHash,
		Seal:     seal,

		ReceivedPromises: make([]Promise, 0),
		AbortedSigners:   make([]AbortedSigner, 0),
		DishonestSigners: make([]uint32, 0),

		SignedContract: make([]byte, 0),
	}
}

// Promise : represents a valid promise
type Promise struct {
	ID bson.ObjectId `key:"_id" bson:"_id"` // Internal id of a Promise

	RecipientKeyIndex uint32 `key:"recipientKeyIndex" bson:"recipientKeyIndex"` // Index of the hash of the recipient's certificate in
	// the `Signers` field of the enclosing SignatureArchives (identical to the one in the signers hashes array of the incoming promises)
	SenderKeyIndex uint32 `key:"senderKeyIndex" bson:"senderKeyIndex"` // Index of the hash of the sender's certificate in
	// the `Signers` field of the enclosing SignatureArchives (identical to the one in the signers hashes array of the incoming promises)
	SequenceIndex uint32 `key:"sequenceIndex" bson:"sequenceIndex"` // Sequence index of the promise
}

// NewPromise : creates a new Promise with the specified fields
func NewPromise(recipientIndex, senderIndex, sequenceIndex uint32) *Promise {
	return &Promise{
		ID:                bson.NewObjectId(),
		RecipientKeyIndex: recipientIndex,
		SenderKeyIndex:    senderIndex,
		SequenceIndex:     sequenceIndex,
	}
}

// Equal : determines if the two specified promises share the same information, without considering the
// bson object id.
func (p1 *Promise) Equal(p2 *Promise) bool {
	if p1.RecipientKeyIndex != p2.RecipientKeyIndex {
		return false
	}
	if p1.SenderKeyIndex != p2.SenderKeyIndex {
		return false
	}
	if p1.SequenceIndex != p2.SequenceIndex {
		return false
	}
	return true
}

// AbortedSigner : represents a signer who was sent an abort token
type AbortedSigner struct {
	ID          bson.ObjectId `key:"_id" bson:"_id"`                 // Internal id of an AbortedSigner
	SignerIndex uint32        `key:"signerIndex" bson:"signerIndex"` // Index of the signer in the signers set
	AbortIndex  uint32        `key:"abortIndex" bson:"abortIndex"`   // Index in the sequence where the signer contacted the ttp and recieveed an abort token
}

// NewAbortedSigner : creates a new AbortedSigner with the specified fields
func NewAbortedSigner(signerIndex, abortIndex uint32) *AbortedSigner {
	return &AbortedSigner{
		ID:          bson.NewObjectId(),
		SignerIndex: signerIndex,
		AbortIndex:  abortIndex,
	}
}
