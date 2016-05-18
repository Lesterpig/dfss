package common

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindNextIndex(t *testing.T) {
	s := []uint32{0, 1, 2, 0, 1, 2}
	id := uint32(2)

	res, err := FindNextIndex(s, id, 0)
	assert.Equal(t, res, 2)
	assert.Equal(t, err, nil)

	res, err = FindNextIndex(s, id, 2)
	assert.Equal(t, res, 5)
	assert.Equal(t, err, nil)

	res, err = FindNextIndex(s, id, 5)
	assert.Equal(t, res, -1)
	assert.Equal(t, err, nil)

	res, err = FindNextIndex(s, id, -2)
	assert.Equal(t, res, -1)
	assert.Equal(t, err.Error(), "Index out of range")

	res, err = FindNextIndex(s, id, len(s))
	assert.Equal(t, res, -1)
	assert.Equal(t, err.Error(), "Index out of range")

	res, err = FindNextIndex(s, 0, -1)
	assert.Equal(t, res, 0)
	assert.Equal(t, err, nil)

	res, err = FindNextIndex(s, 0, 0)
	assert.Equal(t, res, 3)
	assert.Equal(t, err, nil)
}

func TestGetPendingSet(t *testing.T) {
	s := []uint32{0, 1, 2, 0, 1, 2}
	id := uint32(2)

	res, err := GetPendingSet(s, id, 0)
	assert.Equal(t, len(res), 0)
	assert.Equal(t, err, nil)

	res, err = GetPendingSet(s, id, 1)
	assert.Equal(t, len(res), 1)
	assert.Equal(t, res[0].Signer, uint32(0))
	assert.Equal(t, res[0].Index, uint32(0))
	assert.Equal(t, err, nil)

	res, err = GetPendingSet(s, id, 2)
	assert.Equal(t, len(res), 2)
	assert.Equal(t, res[0].Signer, uint32(1))
	assert.Equal(t, res[0].Index, uint32(1))
	assert.Equal(t, res[1].Signer, uint32(0))
	assert.Equal(t, res[1].Index, uint32(0))
	assert.Equal(t, err, nil)

	res, err = GetPendingSet(s, id, 5)
	assert.Equal(t, len(res), 2)
	assert.Equal(t, res[0].Signer, uint32(1))
	assert.Equal(t, res[0].Index, uint32(4))
	assert.Equal(t, res[1].Signer, uint32(0))
	assert.Equal(t, res[1].Index, uint32(3))
	assert.Equal(t, err, nil)

	res, err = GetPendingSet(s, id, -1)
	assert.Equal(t, len(res), 0)
	assert.Equal(t, err.Error(), "Index out of range")

	res, err = GetPendingSet(s, id, len(s))
	assert.Equal(t, len(res), 0)
	assert.Equal(t, err.Error(), "Index out of range")

	s = []uint32{0, 1, 0, 2}
	res, err = GetPendingSet(s, id, 3)
	assert.Equal(t, len(res), 2)
	assert.Equal(t, res[0].Signer, uint32(0))
	assert.Equal(t, res[0].Index, uint32(2))
	assert.Equal(t, res[1].Signer, uint32(1))
	assert.Equal(t, res[1].Index, uint32(1))
	assert.Equal(t, err, nil)
}

func TestContainsCoordinate(t *testing.T) {
	var s []SequenceCoordinate
	s = append(s, SequenceCoordinate{Signer: 0, Index: 0})
	s = append(s, SequenceCoordinate{Signer: 1, Index: 1})

	assert.Equal(t, containsCoordinate(s, 0), true)
	assert.Equal(t, containsCoordinate(s, 1), true)
	assert.Equal(t, containsCoordinate(s, 2), false)
}

func TestContains(t *testing.T) {
	s := []uint32{0, 1}

	assert.Equal(t, contains(s, 0), true)
	assert.Equal(t, contains(s, 1), true)
	assert.Equal(t, contains(s, 2), false)
}

func TestGetSendSet(t *testing.T) {
	s := []uint32{0, 1, 2, 0, 1, 2}
	id := uint32(0)

	res, err := GetSendSet(s, id, 3)
	assert.Equal(t, len(res), 2)
	assert.Equal(t, res[0].Signer, uint32(1))
	assert.Equal(t, res[0].Index, uint32(4))
	assert.Equal(t, res[1].Signer, uint32(2))
	assert.Equal(t, res[1].Index, uint32(5))
	assert.Equal(t, err, nil)

	res, err = GetSendSet(s, id, 0)
	assert.Equal(t, len(res), 2)
	assert.Equal(t, res[0].Signer, uint32(1))
	assert.Equal(t, res[0].Index, uint32(1))
	assert.Equal(t, res[1].Signer, uint32(2))
	assert.Equal(t, res[1].Index, uint32(2))
	assert.Equal(t, err, nil)

	res, err = GetSendSet(s, 1, 4)
	assert.Equal(t, len(res), 1)
	assert.Equal(t, res[0].Signer, uint32(2))
	assert.Equal(t, res[0].Index, uint32(5))
	assert.Equal(t, err, nil)

	res, err = GetSendSet(s, 2, 5)
	assert.Equal(t, len(res), 0)
	assert.Equal(t, err, nil)

	res, err = GetSendSet(s, id, -1)
	assert.Equal(t, len(res), 0)
	assert.Equal(t, err.Error(), "Index out of range")

	res, err = GetSendSet(s, id, len(s))
	assert.Equal(t, len(res), 0)
	assert.Equal(t, err.Error(), "Index out of range")
}

func TestGetAllButOne(t *testing.T) {
	s := []uint32{0, 1, 2, 0, 1, 2}
	id := uint32(2)

	res := GetAllButOne(s, id)
	assert.Equal(t, len(res), 2)
	assert.Equal(t, res[0], uint32(0))
	assert.Equal(t, res[1], uint32(1))

	res = GetAllButOne(s, uint32(42))
	assert.Equal(t, len(res), 3)
	assert.Equal(t, res[0], uint32(0))
	assert.Equal(t, res[1], uint32(1))
	assert.Equal(t, res[2], uint32(2))
}

func TestRemoveCoordinate(t *testing.T) {
	var s []SequenceCoordinate
	s = append(s, SequenceCoordinate{Signer: 0, Index: 0})
	s = append(s, SequenceCoordinate{Signer: 1, Index: 1})

	assert.Equal(t, len(s), 2)
	assert.Equal(t, containsCoordinate(s, 0), true)
	assert.Equal(t, containsCoordinate(s, 1), true)

	s2, err := RemoveCoordinate(s, 3)
	assert.Equal(t, err.Error(), "ID not in sequence")
	assert.Equal(t, len(s2), 2)
	assert.Equal(t, containsCoordinate(s2, 0), true)
	assert.Equal(t, containsCoordinate(s2, 1), true)

	s3, err := RemoveCoordinate(s, 0)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(s3), 1)
	assert.Equal(t, containsCoordinate(s3, 0), false)
	assert.Equal(t, containsCoordinate(s3, 1), true)
}

func TestRemove(t *testing.T) {
	s := []uint32{0, 1}

	assert.Equal(t, len(s), 2)

	s2, err := Remove(s, 3)
	assert.Equal(t, err.Error(), "ID not in sequence")
	assert.Equal(t, len(s2), 2)
	assert.Equal(t, contains(s2, 0), true)
	assert.Equal(t, contains(s2, 1), true)

	s3, err := Remove(s, 0)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(s3), 1)
	assert.Equal(t, contains(s3, 0), false)
	assert.Equal(t, contains(s3, 1), true)
}

func ExampleFindNextIndex() {
	s := []uint32{0, 1, 2, 0, 1, 2}
	id := uint32(2)

	index, _ := FindNextIndex(s, id, -1)
	fmt.Println("First index:", index)

	pSet, _ := GetPendingSet(s, id, index)
	fmt.Println("Pending Set:", pSet)

	sSet, _ := GetSendSet(s, id, index)
	fmt.Println("Send Set:", sSet)
	// Output:
	// First index: 2
	// Pending Set: [{1 1} {0 0}]
	// Send Set: [{0 3} {1 4}]
}
