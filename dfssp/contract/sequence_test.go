package contract_test

import (
	"dfss/dfssp/contract"
	"dfss/dfssp/entities"
	"github.com/bmizerany/assert"
	"testing"
)

var refSeq = []int{0, 1, 2, 0, 1, 2, 0, 1, 2} // for n = 3

func TestGenerateSignSequence(t *testing.T) {

	// initialise fixtures
	var users = make([]*entities.User, 3)
	users[0] = entities.NewUser()
	users[1] = entities.NewUser()
	users[2] = entities.NewUser()
	users[0].Email = "user1@example.com"
	users[1].Email = "user2@example.com"
	users[2].Email = "user3@example.com"

	assert.Equal(t, refSeq, contract.GenerateSignSequence(users))
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
