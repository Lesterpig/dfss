package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArePromisesEqual(t *testing.T) {
	promise0 := NewPromise(uint32(1), uint32(0), uint32(0))
	promise1 := NewPromise(uint32(0), uint32(1), uint32(1))

	equal := ArePromisesEqual(promise0, promise1)
	assert.Equal(t, equal, false)

	promise1.RecipientKeyIndex = uint32(1)
	equal = ArePromisesEqual(promise0, promise1)
	assert.Equal(t, equal, false)

	promise1.SenderKeyIndex = uint32(0)
	equal = ArePromisesEqual(promise0, promise1)
	assert.Equal(t, equal, false)

	promise1.SequenceIndex = uint32(0)
	equal = ArePromisesEqual(promise0, promise1)
	assert.Equal(t, equal, true)
}
