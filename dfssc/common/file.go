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
