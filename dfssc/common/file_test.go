package common

import (
	"crypto/sha512"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

func TestUnmarshalDFSSFile(t *testing.T) {
	data, err := ioutil.ReadFile("../testdata/file.dfss")
	assert.Equal(t, nil, err)

	c, err := UnmarshalDFSSFile(data)
	assert.Equal(t, nil, err)
	assert.Equal(t, "filename.pdf", c.File.Name)
	// a tiny test, full unmarshal is tested in dfss/dfssp/contract package
}

func TestUnmarshalRecoverDataFile(t *testing.T) {
	bsonUUID := bson.NewObjectId()
	uuid := bsonUUID.Hex()
	ttpAddrport := "127.0.0.1:0000"

	h := sha512.Sum512([]byte{7})
	ttpHash := h[:]

	recData := RecoverDataJSON{
		SignatureUUID: uuid,
		TTPAddrport:   ttpAddrport,
		TTPHash:       ttpHash,
	}

	file, err := json.MarshalIndent(recData, "", "  ")
	assert.Nil(t, err)

	unmarshal, err := UnmarshalRecoverDataFile(file)
	assert.Nil(t, err)

	assert.Equal(t, uuid, unmarshal.SignatureUUID)
	assert.Equal(t, ttpAddrport, unmarshal.TTPAddrport)
	assert.Equal(t, ttpHash, unmarshal.TTPHash)
}
