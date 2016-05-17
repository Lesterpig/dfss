package entities

import (
	"bytes"
	"errors"

	"crypto/sha512"
	"dfss/auth"
	cAPI "dfss/dfssc/api"
	pAPI "dfss/dfssp/api"
	tAPI "dfss/dfsst/api"
	"dfss/net"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2/bson"
)

// IsRequestValid : determines if there are no errors in the received request.
// ie: the information signed by the platform in the received promises is valid and consistent
//     the sender of the request is present amongst the signed signers of the promises
func IsRequestValid(ctx context.Context, request *tAPI.AlertRequest) (valid bool, signatureUUID bson.ObjectId, signers []Signer, senderIndex uint32) {
	// Due to specifications, there should be at least one promise (from the sender to himself)
	if len(request.Promises) == 0 {
		valid = false
		return
	}

	ok, expectedUUID, signers := IsPromiseSignedByPlatform(request.Promises[0])
	if !ok {
		valid = false
		return
	}

	sender, err := GetSenderHashFromContext(ctx)
	if err != nil {
		valid = false
		return
	}

	senderIndex, err = GetIndexOfSigner(request.Promises[0], sender)
	if err != nil {
		valid = false
		return
	}

	// To check that all the promises contain the same signed information, we only need to check that:
	// - it is correctly signed
	// - promises are consistent wrt at least one signed field
	for _, promise := range request.Promises {
		ok, receivedUUID, _ := IsPromiseSignedByPlatform(promise)
		if !ok || (expectedUUID != receivedUUID) {
			valid = false
			return
		}
	}

	return true, expectedUUID, signers, senderIndex
}

// IsPromiseSignedByPlatform : determines if the specified promise contains valid information,
// correctly signed by the platform, and returns the signatureUUID if true.
func IsPromiseSignedByPlatform(promise *cAPI.Promise) (bool, bson.ObjectId, []Signer) {
	ok, signatureUUID := IsSignatureUUIDValid(promise)
	if !ok {
		return false, signatureUUID, nil
	}

	ok, signers := AreSignersHashesValid(promise)
	if !ok {
		return false, signatureUUID, nil
	}

	ok = IsPlatformSealValid(promise)
	if !ok {
		return false, signatureUUID, nil
	}

	return true, signatureUUID, signers
}

// GetSenderHashFromContext : creates the sender's certificate hash from the specified context.
func GetSenderHashFromContext(ctx context.Context) ([]byte, error) {
	state, _, ok := net.GetTLSState(&ctx)
	if !ok || len(state.VerifiedChains) == 0 {
		return nil, errors.New("Empty verified sender certificate")
	}

	cert := state.VerifiedChains[0][0]
	hash := auth.GetCertificateHash(cert)

	return hash, nil
}

// GetIndexOfSigner : determines the index of the specified signer's hash in the array of signers' hashes.
func GetIndexOfSigner(promise *cAPI.Promise, hash []byte) (uint32, error) {
	for i, h := range promise.Context.Signers {
		if bytes.Equal(h, hash) {
			return uint32(i), nil
		}
	}
	return 0, errors.New("Signer's hash couldn't be matched")
}

// IsSignatureUUIDValid : verifies that the specified promise has a valid bons.ObjectId hex, and returns the ObjectId if true.
func IsSignatureUUIDValid(promise *cAPI.Promise) (bool, bson.ObjectId) {
	if bson.IsObjectIdHex(promise.Context.SignatureUUID) {
		return true, bson.ObjectIdHex(promise.Context.SignatureUUID)
	}

	return false, bson.NewObjectId()
}

// AreSignersHashesValid : verifies that all the specified hashes are valid (see function IsSignerHashValid).
// Returns a new array of Signers.
func AreSignersHashesValid(promise *cAPI.Promise) (bool, []Signer) {
	var signers []Signer
	if len(promise.Context.Signers) == 0 {
		return false, nil
	}

	for _, v := range promise.Context.Signers {
		ok, signer := IsSignerHashValid(v)
		if !ok {
			return false, nil
		}

		signers = append(signers, *signer)
	}
	return true, signers
}

// IsSignerHashValid : verifies that the specified array of bytes is a correct SHA-512 hash.
// Returns a new Signer with the specified hash.
func IsSignerHashValid(hash []byte) (bool, *Signer) {
	if sha512.Size != len(hash) {
		return false, nil
	}

	return true, NewSigner(hash)
}

// IsPlatformSealValid : verifies that the specified promise contains the expected information signed by the platform.
func IsPlatformSealValid(promise *cAPI.Promise) bool {
	if AuthContainer == nil {
		return false
	}

	theoric := pAPI.LaunchSignature{
		SignatureUuid: promise.Context.SignatureUUID,
		KeyHash:       promise.Context.Signers,
		Sequence:      promise.Context.Sequence,
	}

	ok, _ := auth.VerifyStructure(AuthContainer.CA, theoric, promise.Context.Seal)
	return ok
}
