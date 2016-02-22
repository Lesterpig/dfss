package templates

import (
	"testing"

	"dfss/dfssp/entities"
	"github.com/bmizerany/assert"
)

func TestInit(t *testing.T) {
	Init() // will panic if any error found in templates
}

func TestGet(t *testing.T) {

	contract := entities.NewContract()
	contract.File.Hash = []byte{0x01, 0x02, 0x11, 0xaa}
	contract.File.Name = "name.pdf"
	contract.Comment = "comment"
	contract.AddSigner(nil, "mail@example.com", nil)
	contract.AddSigner(nil, "mail2@example.com", nil)

	s, err := Get("contract", contract)

	expected := `Dear Sir or Madam,

Someone asked you to sign a contract on the DFSS platform.
Please download the attached file and open it with the DFSS client.

Signers :
  - mail@example.com
  - mail2@example.com

Contract name : name.pdf
SHA-512 hash  : 010211aa
Comment       : comment

Yours faithfully,

The DFSS Platform
`

	assert.Equal(t, nil, err)
	assert.Equal(t, expected, s)
}
