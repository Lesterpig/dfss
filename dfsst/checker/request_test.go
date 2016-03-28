package checker

import (
	"bytes"
	"crypto/sha512"
	cAPI "dfss/dfssc/api"
	"github.com/bmizerany/assert"
	"gopkg.in/mgo.v2/bson"
	"testing"
)

var (
	sequence             []uint32
	signers              [][]byte
	contractDocumentHash []byte
	signatureUUID        string
	signatureUUIDBson    bson.ObjectId

	signedHash []byte
)

func init() {
	sequence = []uint32{0, 1, 2, 0, 1, 2, 0, 1, 2}

	for i := 0; i < 3; i++ {
		h := sha512.Sum512([]byte{byte(i)})
		signer := h[:]
		signers = append(signers, signer)
	}

	contractDocumentHash = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}
	signatureUUIDBson = bson.NewObjectId()
	signatureUUID = signatureUUIDBson.Hex()

	signedHash = []byte{}
}

func TestIsRequestValid(t *testing.T) {
	// TODO
	// This requires the use of a real Alert message
}

func TestIsPromiseSignedByPlatform(t *testing.T) {
	promise := &cAPI.Promise{
		Context: &cAPI.Context{
			SignatureUUID: "toto",
			Signers:       [][]byte{},
		},
	}

	valid, _, _ := IsPromiseSignedByPlatform(promise)
	assert.Equal(t, valid, false)

	promise.Context.SignatureUUID = signatureUUID
	valid, sigID, _ := IsPromiseSignedByPlatform(promise)
	assert.Equal(t, valid, false)

	promise.Context.Signers = append(promise.Context.Signers, []byte{0})
	valid, sigID, _ = IsPromiseSignedByPlatform(promise)
	assert.Equal(t, valid, false)

	// TODO
	// when 'IsPlatformSignedHashValid' is implemented
	promise.Context.Signers = signers
	valid, sigID, signerss := IsPromiseSignedByPlatform(promise)
	assert.Equal(t, valid, true)
	assert.Equal(t, sigID, signatureUUIDBson)
	assert.Equal(t, len(signerss), len(signers))
}

func TestGetSenderHashFromContext(t *testing.T) {
	// TODO
	// This requires the use of a real Alert message
}

func TestGetIndexOfSigner(t *testing.T) {
	promise := &cAPI.Promise{
		Context: &cAPI.Context{
			Signers: signers,
		},
	}

	hash := []byte{}
	i, err := GetIndexOfSigner(promise, hash)
	assert.Equal(t, err.Error(), "Signer's hash couldn't be matched")
	assert.Equal(t, i, uint32(0))

	hash = signers[1]
	i, err = GetIndexOfSigner(promise, hash)
	assert.Equal(t, err, nil)
	assert.Equal(t, i, uint32(1))
}

func TestIsSignatureUUIDValid(t *testing.T) {
	promise := &cAPI.Promise{
		Context: &cAPI.Context{
			SignatureUUID: "toto",
		},
	}
	b, _ := IsSignatureUUIDValid(promise)
	assert.Equal(t, b, false)

	promise.Context.SignatureUUID = signatureUUID
	b, id := IsSignatureUUIDValid(promise)
	assert.Equal(t, id, signatureUUIDBson)
	assert.Equal(t, b, true)
}

func TestAreSignersHashesValid(t *testing.T) {
	promise := &cAPI.Promise{
		Context: &cAPI.Context{
			Signers: [][]byte{},
		},
	}

	b, signerss := AreSignersHashesValid(promise)
	assert.Equal(t, b, false)
	assert.Equal(t, len(signerss), 0)

	promise.Context.Signers = append(promise.Context.Signers, []byte{0})
	b, signerss = AreSignersHashesValid(promise)
	assert.Equal(t, b, false)
	assert.Equal(t, len(signerss), 0)

	promise.Context.Signers = signers
	b, signerss = AreSignersHashesValid(promise)
	assert.Equal(t, b, true)
	assert.Equal(t, len(signerss), len(signers))
	for i, v := range signerss {
		assert.Equal(t, v.Hash, signers[i])
	}
}

func TestIsSignerHashValid(t *testing.T) {
	hash := []byte{}
	b, signer := IsSignerHashValid(hash)
	assert.Equal(t, b, false)

	b, signer = IsSignerHashValid(signers[0])
	assert.Equal(t, b, true)
	assert.Equal(t, bytes.Equal(signer.Hash, signers[0]), true)
}

// TO MODIFY WHEN SOURCE FUNCTION WILL BE UPDATED
func TestIsPlatformSignedHashValid(t *testing.T) {
	promise := &cAPI.Promise{
		Context: &cAPI.Context{
			ContractDocumentHash: contractDocumentHash,
			Sequence:             sequence,
			Signers:              signers,
			SignatureUUID:        signatureUUID,
			SignedHash:           signedHash,
		},
	}

	b := IsPlatformSignedHashValid(promise)
	assert.Equal(t, b, true)
}
