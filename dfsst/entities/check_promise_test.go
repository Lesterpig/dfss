package entities

import (
	"testing"

	cAPI "dfss/dfssc/api"
	"github.com/stretchr/testify/assert"
)

func TestArePromisesValid(t *testing.T) {
	promise0 := &cAPI.Promise{
		Context: &cAPI.Context{
			RecipientKeyHash: []byte{},
			SenderKeyHash:    []byte{},
			Sequence:         sequence,
			Signers:          signers,
		},
		Index: uint32(len(sequence)),
	}

	promise1 := &cAPI.Promise{
		Context: &cAPI.Context{
			RecipientKeyHash: []byte{},
			SenderKeyHash:    []byte{},
			Sequence:         sequence,
			Signers:          signers,
		},
		Index: uint32(len(sequence)),
	}

	promises := []*cAPI.Promise{promise0, promise1}
	valid, promiseEntities := ArePromisesValid(promises)
	assert.Equal(t, valid, false)
	assert.Equal(t, promiseEntities, []*Promise(nil))

	promise0.Context.RecipientKeyHash = signers[2]
	promise0.Context.SenderKeyHash = signers[1]
	promise0.Index = 1
	valid, promiseEntities = ArePromisesValid(promises)
	assert.Equal(t, valid, false)
	assert.Equal(t, promiseEntities, []*Promise(nil))

	promise1.Context.RecipientKeyHash = signers[2]
	promise1.Context.SenderKeyHash = signers[0]
	promise1.Index = 0
	valid, promiseEntities = ArePromisesValid(promises)
	assert.Equal(t, valid, true)
	assert.Equal(t, promiseEntities[0].RecipientKeyIndex, uint32(2))
	assert.Equal(t, promiseEntities[0].SenderKeyIndex, uint32(1))
	assert.Equal(t, promiseEntities[0].SequenceIndex, uint32(1))
	assert.Equal(t, promiseEntities[1].RecipientKeyIndex, uint32(2))
	assert.Equal(t, promiseEntities[1].SenderKeyIndex, uint32(0))
	assert.Equal(t, promiseEntities[1].SequenceIndex, uint32(0))
}

func TestIsPromiseValid(t *testing.T) {
	// TODO
	// Test with a promess not from 'A' to 'B'
	promise := &cAPI.Promise{
		Context: &cAPI.Context{
			RecipientKeyHash: []byte{},
			SenderKeyHash:    []byte{},
			Sequence:         sequence,
			Signers:          signers,
		},
		Index: uint32(len(sequence)),
	}

	valid, promiseEntity := IsPromiseValid(promise)
	assert.Equal(t, valid, false)
	assert.Equal(t, promiseEntity, &Promise{})

	promise.Context.RecipientKeyHash = signers[2]
	promise.Context.SenderKeyHash = signers[1]
	promise.Index = 1
	valid, promiseEntity = IsPromiseValid(promise)
	assert.Equal(t, valid, true)
	assert.Equal(t, promiseEntity.RecipientKeyIndex, uint32(2))
	assert.Equal(t, promiseEntity.SenderKeyIndex, uint32(1))
	assert.Equal(t, promiseEntity.SequenceIndex, uint32(1))
}

func TestIsPromiseFromAtoB(t *testing.T) {
	// TODO
	// This requires the implementation of promises
	sender := []byte{1}
	recipient := []byte{2}
	index := uint32(5)
	promise := &cAPI.Promise{
		Context: &cAPI.Context{
			RecipientKeyHash: recipient,
			SenderKeyHash:    sender,
			Signers:          signers,
		},
		Index: index,
	}
	ok := IsPromiseFromAtoB(promise, sender, recipient, index)
	assert.Equal(t, ok, true)
}

func TestGetPromiseProfile(t *testing.T) {
	promise := &cAPI.Promise{
		Context: &cAPI.Context{
			RecipientKeyHash: []byte{},
			SenderKeyHash:    []byte{},
			Sequence:         sequence,
			Signers:          signers,
		},
		Index: uint32(len(sequence)),
	}

	recipient, sender, index, err := GetPromiseProfile(promise)
	assert.Equal(t, recipient, uint32(0))
	assert.Equal(t, sender, uint32(0))
	assert.Equal(t, index, uint32(0))
	assert.Equal(t, err.Error(), "Signer's hash couldn't be matched")

	promise.Context.RecipientKeyHash = signers[2]
	recipient, sender, index, err = GetPromiseProfile(promise)
	assert.Equal(t, recipient, uint32(0))
	assert.Equal(t, sender, uint32(0))
	assert.Equal(t, index, uint32(0))
	assert.Equal(t, err.Error(), "Signer's hash couldn't be matched")

	promise.Context.SenderKeyHash = signers[1]
	recipient, sender, index, err = GetPromiseProfile(promise)
	assert.Equal(t, recipient, uint32(0))
	assert.Equal(t, sender, uint32(0))
	assert.Equal(t, index, uint32(0))
	assert.Equal(t, err.Error(), "Index out of range in the sequence")

	promise.Index = 0
	recipient, sender, index, err = GetPromiseProfile(promise)
	assert.Equal(t, recipient, uint32(0))
	assert.Equal(t, sender, uint32(0))
	assert.Equal(t, index, uint32(0))
	assert.Equal(t, err.Error(), "Sequence id at promise index of promise does not match index sender sequence id")

	promise.Index = 1
	recipient, sender, index, err = GetPromiseProfile(promise)
	assert.Equal(t, recipient, uint32(2))
	assert.Equal(t, sender, uint32(1))
	assert.Equal(t, index, uint32(1))
	assert.Equal(t, err, nil)
}

func TestIsIndexValid(t *testing.T) {
	promise := &cAPI.Promise{
		Context: &cAPI.Context{
			SenderKeyHash: []byte{},
			Sequence:      sequence,
			Signers:       signers,
		},
		Index: uint32(len(sequence)),
	}

	index, err := IsIndexValid(promise)
	assert.Equal(t, index, uint32(0))
	assert.Equal(t, err.Error(), "Index out of range in the sequence")

	promise.Index = 1
	index, err = IsIndexValid(promise)
	assert.Equal(t, index, uint32(0))
	assert.Equal(t, err.Error(), "Signer's hash couldn't be matched")

	promise.Context.SenderKeyHash = signers[0]
	index, err = IsIndexValid(promise)
	assert.Equal(t, index, uint32(0))
	assert.Equal(t, err.Error(), "Sequence id at promise index of promise does not match index sender sequence id")

	promise.Context.SenderKeyHash = signers[1]
	index, err = IsIndexValid(promise)
	assert.Equal(t, index, uint32(1))
	assert.Equal(t, err, nil)
}
