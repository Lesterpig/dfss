package common

import (
	"errors"
)

// FindNextIndex analyses the specified sequence and tries to find the next occurence of id after the specified index (excluded)
// Therefore to find the first occurence of the id in the sequence, use -1 as the index
//
// If there is no occurence of id, then -1 is returned
//
// The sequence is supposed to be correct (in regards to its mathematical definition), and id is supposed to be a valid id for the sequence
func FindNextIndex(s []uint32, id uint32, index int) (int, error) {
	if index >= len(s) || index < -1 {
		return -1, errors.New("Index out of range")
	}

	for i := index + 1; i < len(s); i++ {
		if s[i] == id {
			return i, nil
		}
	}

	return -1, nil
}

// GetPendingSet analyses the specified sequence and computes the set of ids occuring between index (excluded) and the previous occurence of id
//
// The sequence is supposed to be correct (in regards to its mathematical definition), and id is supposed to be a valid id for the sequence
//
// If the index is not the one of the specified id in the sequence, the result still holds, but may be incomplete for your needs
// If the id is not valid for the specified sequence, the result will be the set of ids of the sequence
func GetPendingSet(s []uint32, id uint32, index int) ([]uint32, error) {
	res := []uint32{}

	if index >= len(s) || index < 0 {
		return res, errors.New("Index out of range")
	}

	curIndex := index - 1
	if curIndex < 0 {
		return res, nil
	}

	for curIndex > -1 && s[curIndex] != id {
		curID := s[curIndex]
		if !contains(res, curID) {
			res = append(res, curID)
		}
		curIndex--
	}

	return res, nil
}

// GetSendSet analyses the specified sequence and computes the set of ids occuring between index (excluded) and the next occurence of id
//
// The sequence is supposed to be correct (in regards to its mathematical definition), and id is supposed to be a valid id for the sequence
//
// If the index is not the one of the specified id in the sequence, the result still holds, but may be incomplete for your needs
// If the id is not valid for the specified sequence, the result will be the set of ids of the sequence
func GetSendSet(s []uint32, id uint32, index int) ([]uint32, error) {
	res := []uint32{}

	if index >= len(s) || index < 0 {
		return res, errors.New("Index out of range")
	}

	curIndex := index + 1
	if curIndex >= len(s) {
		return res, nil
	}

	for curIndex < len(s) && s[curIndex] != id {
		curID := s[curIndex]
		if !contains(res, curID) {
			res = append(res, curID)
		}
		curIndex++
	}

	return res, nil
}

// contains determines if s contains e
func contains(s []uint32, e uint32) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// GetAllButOne creates the slice of all sequence ids, except the one specified
func GetAllButOne(s []uint32, e uint32) []uint32 {
	var res = make([]uint32, 0)

	for i := 0; i < len(s); i++ {
		curID := s[i]
		if !contains(res, curID) && curID != e {
			res = append(res, curID)
		}
	}

	return res
}

// Remove an ID from the sequence
func Remove(s []uint32, e uint32) ([]uint32, error) {
	for i, a := range s {
		if a == e {
			return append(s[:i], s[i+1:]...), nil
		}
	}
	return s, errors.New("ID not in sequence")
}
