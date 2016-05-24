// Package resolve provides the resolve protocol.
package resolve

import (
	"errors"
	"fmt"

	cAPI "dfss/dfssc/api"
	"dfss/dfssc/common"
	dAPI "dfss/dfssd/api"
	"dfss/dfsst/entities"
)

// ArePromisesComplete : determines if the set of promises present in the AlertRequest is EQUAL (not just included) to the one expected from the TTP
// for this signer, at this particular step of the signing sequence.
// The provided step is supposed valid (ie. in bound of the sequence) (see call to 'IsRequestValid' in 'Alert' route)
func ArePromisesComplete(promiseEntities []*entities.Promise, promise *cAPI.Promise, step uint32) bool {
	expected, err := generateExpectedPromises(promise, step)
	if err != nil {
		dAPI.DLog("error occured during the generation of expected promises")
		return false
	}

	if len(promiseEntities) != len(expected) {
		dAPI.DLog("promise sets are not equal, nb received: " + fmt.Sprint(len(promiseEntities)) + " ; expected: " + fmt.Sprint(len(expected)))
		return false
	}

	for i, p := range expected {
		if !containsPromise(promiseEntities, p) {
			dAPI.DLog("promise sets are not equal, promise at index: " + fmt.Sprint(i))
			return false
		}
	}

	dAPI.DLog("promise sets are equal")
	return true
}

// generateExpectedPromises : computes the list of expected promises for the recipient of the specified promise, at the specified sequence index
func generateExpectedPromises(promise *cAPI.Promise, step uint32) ([]*entities.Promise, error) {
	var res []*entities.Promise

	seq := promise.Context.Sequence
	recipientH := promise.Context.RecipientKeyHash
	recipientID, err := entities.GetIndexOfSigner(promise, recipientH)
	if err != nil {
		dAPI.DLog(err.Error())
		return nil, err
	}

	// if step is 0, then it means that the resolve call occured during the first round.
	// therefore, there is no way we can generate the signed contract.
	// so we don't check the 0 case, because it has no consequence on the folloing of the resolve algorithm: an abort token will be sent
	// checking the 0 case would cause an error in the consistency between of the resolve index and the promise recipient sequence index.
	if (step != 0) && (seq[int(step)] != recipientID) {
		dAPI.DLog("sequence index at step " + fmt.Sprint(int(step)) + " is " + fmt.Sprint(seq[int(step)]) + ", recipientID is " + fmt.Sprint(recipientID))
		return nil, errors.New("Signer at step is not recipient")
	}

	currentIndex, err := common.FindNextIndex(seq, recipientID, -1)
	if err != nil {
		dAPI.DLog(err.Error())
		return nil, err
	}

	dAPI.DLog("resolve index is: " + fmt.Sprint(step))
	dAPI.DLog("first index is: " + fmt.Sprint(currentIndex))
	for currentIndex <= int(step) {
		dAPI.DLog("started generation round with currentIndex " + fmt.Sprint(currentIndex))
		roundPromises, err := generationRound(seq, recipientID, currentIndex)
		if err != nil {
			dAPI.DLog("error occured during the generation round, currentIndex: " + fmt.Sprint(currentIndex))
			return nil, err
		}

		dAPI.DLog(fmt.Sprint(len(roundPromises)) + " promises were generated for this round")
		for _, p := range roundPromises {
			res = addPromiseToExpected(res, p)
		}

		dAPI.DLog("total number of expected promises this far: " + fmt.Sprint(len(res)))
		currentIndex, err = common.FindNextIndex(seq, recipientID, currentIndex)
		if err != nil {
			dAPI.DLog(err.Error())
		}
		// if it was the last occurence, then we finish
		if currentIndex == -1 {
			break
		}
	}

	dAPI.DLog("promise from rounds have been generated")
	selfPromise := &entities.Promise{
		RecipientKeyIndex: recipientID,
		SenderKeyIndex:    recipientID,
		SequenceIndex:     step,
	}

	dAPI.DLog("adding self promise")
	return append(res, selfPromise), nil
}

// generationRound : generates the list of promises expected bu seqID at index of seq
func generationRound(seq []uint32, seqID uint32, index int) ([]*entities.Promise, error) {
	var res []*entities.Promise

	pendingSet, err := common.GetPendingSet(seq, seqID, index)
	if err != nil {
		return nil, err
	}

	for _, c := range pendingSet {
		p := &entities.Promise{
			RecipientKeyIndex: seqID,
			SenderKeyIndex:    c.Signer,
			SequenceIndex:     c.Index,
		}
		res = append(res, p)
	}

	return res, nil
}

// addPromiseToExpected : adds the specified promise to the list.
// If a previous promise from this sender to this recipient is in the list, replaces it with the new one.
// Otherwise, appends the promise to the list.
func addPromiseToExpected(expected []*entities.Promise, promise *entities.Promise) []*entities.Promise {
	if containsPromise(expected, promise) {
		return expected
	}

	i := containsPreviousPromise(expected, promise)
	if i != -1 {
		expected[i] = promise
	} else {
		expected = append(expected, promise)
	}

	return expected
}

// containsPreviousPromise : determines if the specified promise has a previous occurence in the list (ie. same sender and recipient, but lower index).
// Returns the index at which the previous occurence is in the provided list.
// If there is no previous occurence, returns -1
func containsPreviousPromise(expected []*entities.Promise, promise *entities.Promise) int {
	for i, p := range expected {
		if p.SenderKeyIndex == promise.SenderKeyIndex && p.RecipientKeyIndex == promise.RecipientKeyIndex && p.SequenceIndex < promise.SequenceIndex {
			return i
		}
	}
	return -1
}

// containsPromise : determines if the specified promise is present in the specified list.
func containsPromise(promises []*entities.Promise, promise *entities.Promise) bool {
	for _, p := range promises {
		if p.Equal(promise) {
			return true
		}
	}
	return false
}

// Solve : tries to generate the signed contract from present evidence.
func Solve(manager *entities.ArchivesManager) (bool, []byte) {
	// Test if we can generate the contract
	for i := range manager.Archives.Signers {
		ok := manager.HasSignerPromised(uint32(i))
		if !ok {
			return false, nil
		}
	}

	return true, GenerateSignedContract(manager.Archives)
}

// GenerateSignedContract : generates the signed contract.
// Does not take into account if we have the evidence to do it (see function 'CanGenerate').
//
// XXX : Implementation needs cryptographic promises
func GenerateSignedContract(archives *entities.SignatureArchives) []byte {
	var pseudoContract string
	for _, p := range archives.ReceivedPromises {
		signature := "SIGNATURE FROM SIGNER " + string(archives.Signers[p.SenderKeyIndex].Hash)
		signature += " ON SIGNATURE n° " + string(archives.ID) + "\n"
		pseudoContract += signature
	}

	return []byte(pseudoContract)
}

// ComputeDishonestSigners : computes the dishonest signers from the provided evidence.
// ie: if an already aborted signer sent a message after his call to the ttp, then he is dishonest
func ComputeDishonestSigners(archives *entities.SignatureArchives, evidence []*entities.Promise) []uint32 {
	var res []uint32
	// Used to add only once the signers if several promises point to them as being dishonest
	added := false

	for _, signer := range archives.AbortedSigners {
		added = false
		sIndex := signer.SignerIndex
		for _, p := range evidence {
			if !added && (p.SenderKeyIndex == sIndex) && (p.SequenceIndex > signer.AbortIndex) {
				res = append(res, sIndex)
				added = true
			}
		}
	}

	return res
}
