package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStructure struct {
	FieldA int64
	FieldB []byte
	FieldC *TestStructure
}

func TestSignStructure(t *testing.T) {
	key, err := GeneratePrivateKey(1024)
	assert.Nil(t, err)

	res, err := SignStructure(key, TestStructure{})
	assert.Nil(t, err)
	assert.True(t, len(res) > 0)
}

func TestVerifyStructure(t *testing.T) {
	key, err := GeneratePrivateKey(1024)
	assert.Nil(t, err)

	selfSigned, err := GetSelfSignedCertificate(1, 0, "", "", "", "test", key)
	assert.Nil(t, err)
	cert, err := PEMToCertificate(selfSigned)
	assert.Nil(t, err)

	s := TestStructure{
		FieldA: 5,
		FieldB: []byte{0x01, 0x02},
		FieldC: &TestStructure{},
	}

	res, _ := SignStructure(key, s)
	valid, err := VerifyStructure(cert, s, res)
	assert.Nil(t, err)
	assert.True(t, valid)

	s.FieldB[1] = 0x42
	valid, _ = VerifyStructure(cert, s, res)
	assert.False(t, valid)
}
