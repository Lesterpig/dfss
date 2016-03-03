package contract

import (
	"dfss/dfssp/entities"
)

// GenerateSignSequence for the contract signature
//
// The generated sequence is an array of integers refering to the User array.
func GenerateSignSequence(users []*entities.User) []int {
	return SquaredSignEngine(len(users))
}

// SquaredSignEngine is a basic ^2 engine for sequence generation
func SquaredSignEngine(n int) []int {
	sequence := make([]int, n*n)

	for i := 0; i < n; i++ {
		for k := 0; k < n; k++ {
			sequence[i*n+k] = k
		}
	}

	return sequence
}

// SquaredSignEngineSlice is the same as the above with slicing
func SquaredSignEngineSlice(n int) []int {
	baseSequence := make([]int, n)

	// populate base slice
	for i := 0; i < n; i++ {
		baseSequence[i] = i
	}

	sequence := make([]int, 0, n*n)
	// append n-1 time the slice to itself
	for i := 0; i < n; i++ {
		sequence = append(sequence, baseSequence...)
	}

	return sequence
}