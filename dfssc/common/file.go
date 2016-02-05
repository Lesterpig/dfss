package common

import (
	"encoding/json"

	"dfss/dfssp/contract"
)

// UnmarshalDFSSFile decodes a json-encoded DFSS file
func UnmarshalDFSSFile(data []byte) (*contract.JSON, error) {

	c := &contract.JSON{}
	err := json.Unmarshal(data, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
