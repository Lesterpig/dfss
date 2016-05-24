package entities

import (
	"bytes"
	"errors"

	cAPI "dfss/dfssc/api"
)

// ArePromisesValid : determines if the specified promises contains coherent information wrt the ASSUMED TESTED platform signed information.
func ArePromisesValid(promises []*cAPI.Promise) (bool, []*Promise) {
	var tmpPromises []*Promise

	for _, promise := range promises {
		valid, promiseEntity := IsPromiseValid(promise)
		if !valid {
			return false, nil
		}
		tmpPromises = append(tmpPromises, promiseEntity)
	}

	return true, tmpPromises
}

// IsPromiseValid : determines if the specified promise contains coherent information wrt the ASSUMED TESTED platform signed information.
// ie: the sender and recipient's hashes are correct
//     the index of the promise coresponds to an expected message from the sender in the signed sequence
// If true, returns a new promise entity
func IsPromiseValid(promise *cAPI.Promise) (bool, *Promise) {
	// This checks if the index of the specified promise corresponds to an expected promise from the sender hash of the promise
	sender, recipient, index, err := GetPromiseProfile(promise)
	if err != nil {
		return false, &Promise{}
	}

	entityPromise := NewPromise(sender, recipient, index)

	return true, entityPromise
}

// IsPromiseFromAtoB : determines if the specified promise, supposedly from 'A' to 'B' was indeed created by 'A' for 'B'.
func IsPromiseFromAtoB(promise *cAPI.Promise, from, to []byte, at uint32) bool {
	if !(bytes.Equal(promise.Context.SenderKeyHash, from)) {
		return false
	}
	if !(bytes.Equal(promise.Context.RecipientKeyHash, to)) {
		return false
	}
	return promise.Index == at
}

// GetPromiseProfile : retrieves the indexes of the recipient and sender in the array of signers' hashes, and the index of the promise.
// If the specified promise is not valid wrt to it's index and sender, returns an error.
func GetPromiseProfile(promise *cAPI.Promise) (uint32, uint32, uint32, error) {
	recipient, err := GetIndexOfSigner(promise, promise.Context.RecipientKeyHash)
	if err != nil {
		return 0, 0, 0, err
	}

	sender, err := GetIndexOfSigner(promise, promise.Context.SenderKeyHash)
	if err != nil {
		return 0, 0, 0, err
	}

	index, err := IsIndexValid(promise)
	if err != nil {
		return 0, 0, 0, err
	}

	return recipient, sender, index, nil
}

// IsIndexValid : determines if the index field of the promise is valid wrt the sender hash and the sequence, and returns it.
func IsIndexValid(promise *cAPI.Promise) (uint32, error) {
	index := promise.Index

	if index >= uint32(len(promise.Context.Sequence)) {
		return 0, errors.New("Index out of range in the sequence")
	}

	senderIndex, err := GetIndexOfSigner(promise, promise.Context.SenderKeyHash)
	if err != nil {
		return 0, err
	}

	if senderIndex != promise.Context.Sequence[index] {
		return 0, errors.New("Sequence id at promise index of promise does not match index sender sequence id")
	}

	return index, nil
}
