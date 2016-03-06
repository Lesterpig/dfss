package contract

// GenerateSignSequence for the contract signature
//
// The generated sequence is an array of integers refering to the User array.
func GenerateSignSequence(n int) []uint32 {
	return SquaredSignEngine(uint32(n))
}

// SquaredSignEngine is a basic ^2 engine for sequence generation
func SquaredSignEngine(n uint32) []uint32 {
	sequence := make([]uint32, n*n)

	var i, k uint32
	for i = 0; i < n; i++ {
		for k = 0; k < n; k++ {
			sequence[i*n+k] = k
		}
	}

	return sequence
}

// SquaredSignEngineSlice is the same as the above with slicing
func SquaredSignEngineSlice(n uint32) []uint32 {
	baseSequence := make([]uint32, n)

	// populate base slice
	var i uint32
	for i = 0; i < n; i++ {
		baseSequence[i] = i
	}

	sequence := make([]uint32, 0, n*n)
	// append n-1 time the slice to itself
	for i = 0; i < n; i++ {
		sequence = append(sequence, baseSequence...)
	}

	return sequence
}
