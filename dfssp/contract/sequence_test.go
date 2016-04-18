package contract_test

import (
	"dfss/dfssp/contract"
	"github.com/stretchr/testify/assert"
	"testing"
)

var refSeq = []uint32{0, 1, 2, 0, 1, 2, 0, 1, 2} // for n = 3

func TestGenerateSignSequence(t *testing.T) {
	assert.Equal(t, refSeq, contract.GenerateSignSequence(3))
}

// Perform sequence generation in a loop
func BenchmarkSquaredEngine(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = contract.SquaredSignEngine(10)
	}
}

// Perform sequence generation in a loop with slicing
func BenchmarkSquaredEngineSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = contract.SquaredSignEngineSlice(10)
	}
}
