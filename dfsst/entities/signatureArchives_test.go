package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArePromisesEqual(t *testing.T) {
	promise0 := NewPromise(uint32(1), uint32(0), uint32(0))
	promise1 := NewPromise(uint32(0), uint32(1), uint32(1))

	equal := promise0.Equal(promise1)
	assert.Equal(t, equal, false)

	promise1.RecipientKeyIndex = uint32(1)
	equal = promise0.Equal(promise1)
	assert.Equal(t, equal, false)

	promise1.SenderKeyIndex = uint32(0)
	equal = promise0.Equal(promise1)
	assert.Equal(t, equal, false)

	promise1.SequenceIndex = uint32(0)
	equal = promise0.Equal(promise1)
	assert.Equal(t, equal, true)
}

func TestContainsSigner(t *testing.T) {
	archives := NewSignatureArchives(signatureUUIDBson, sequence, signersEntities, contractDocumentHash, seal)
	badHash := []byte{}
	present, index := archives.ContainsSigner(badHash)
	assert.False(t, present)
	assert.Equal(t, uint32(0), index)

	badHash = []byte{0, 0, 7}
	present, index = archives.ContainsSigner(badHash)
	assert.False(t, present)
	assert.Equal(t, uint32(0), index)

	goodHash := signersEntities[0].Hash
	present, index = archives.ContainsSigner(goodHash)
	assert.True(t, present)
	assert.Equal(t, uint32(0), index)

	goodHash = signersEntities[1].Hash
	present, index = archives.ContainsSigner(goodHash)
	assert.True(t, present)
	assert.Equal(t, uint32(1), index)
}
