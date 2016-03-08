package common

import (
	"fmt"
	"testing"
	"github.com/bmizerany/assert"
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
	assert.Equal(t, res[0], uint32(0))
	assert.Equal(t, err, nil)

	res, err = GetPendingSet(s, id, 2)
	assert.Equal(t, len(res), 2)
	assert.Equal(t, res[0], uint32(1))
	assert.Equal(t, res[1], uint32(0))
	assert.Equal(t, err, nil)

	res, err = GetPendingSet(s, id, 5)
	assert.Equal(t, len(res), 2)
	assert.Equal(t, res[0], uint32(1))
	assert.Equal(t, res[1], uint32(0))
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
	assert.Equal(t, res[0], uint32(0))
	assert.Equal(t, res[1], uint32(1))
	assert.Equal(t, err, nil)
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
	assert.Equal(t, res[0], uint32(1))
	assert.Equal(t, res[1], uint32(2))
	assert.Equal(t, err, nil)

	res, err = GetSendSet(s, id, 0)
	assert.Equal(t, len(res), 2)
	assert.Equal(t, res[0], uint32(1))
	assert.Equal(t, res[1], uint32(2))
	assert.Equal(t, err, nil)

	res, err = GetSendSet(s, 1, 4)
	assert.Equal(t, len(res), 1)
	assert.Equal(t, res[0], uint32(2))
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

func ExampleSequenceAnalysis() {
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
	// Pending Set: [1 0]
	// Send Set: [0 1]
}