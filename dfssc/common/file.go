package common

import (
	"encoding/json"
	"errors"

	"dfss/dfssp/contract"
)

// UnmarshalDFSSFile decodes a json-encoded DFSS file
func UnmarshalDFSSFile(data []byte) (*contract.JSON, error) {
	c := &contract.JSON{}
	err := json.Unmarshal(data, c)
	if err != nil {
		return nil, err
	}

	if c.File == nil {
		return nil, errors.New("empty file description")
	}

	return c, nil
}

// RecoverDataJSON : contains all the necessary information to try and recover a previously signed contract from the ttp
type RecoverDataJSON struct {
	SignatureUUID string
	TTPAddrport   string
	TTPHash       []byte
}

// UnmarshalRecoverDataFile decodes a json-encoded Recover dara file
func UnmarshalRecoverDataFile(data []byte) (*RecoverDataJSON, error) {
	r := &RecoverDataJSON{}
	err := json.Unmarshal(data, r)
	if err != nil {
		return nil, err
	}

	if r.SignatureUUID == "" || r.TTPAddrport == "" || r.TTPHash == nil {
		return nil, errors.New("Invalid recover data file")
	}

	return r, nil
}
