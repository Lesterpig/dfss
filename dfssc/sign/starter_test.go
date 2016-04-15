package sign

import (
	"testing"

	"dfss/dfssp/contract"

	"github.com/bmizerany/assert"
)

func TestFindId(t *testing.T) {
	mail := "signer2@foo.foo"

	s1 := contract.SignerJSON{Email: "signer1@foo.foo"}
	s2 := contract.SignerJSON{Email: mail}
	s3 := contract.SignerJSON{Email: "signer3@foo.foo"}

	contract := &contract.JSON{
		Signers: []contract.SignerJSON{s1, s2, s3},
	}

	sequence := []uint32{0, 1, 2, 0, 1, 2, 0, 1, 2}

	sm := &SignatureManager{
		contract: contract,
		sequence: sequence,
		mail:     mail,
	}

	id, err := sm.FindID()
	assert.Equal(t, err, nil)
	assert.Equal(t, id, uint32(1))

	sm.mail = ""
	id, err = sm.FindID()
	assert.Equal(t, err.Error(), "Mail couldn't be found amongst signers")
	assert.Equal(t, id, uint32(0))
}
