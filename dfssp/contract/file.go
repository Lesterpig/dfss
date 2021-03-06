package contract

import (
	"encoding/json"
	"fmt"
	"time"

	"dfss/dfssp/entities"
)

// FileJSON is the structure used to store file information in JSON format
type FileJSON struct {
	Name   string
	Hash   string
	Hosted bool
}

// SignerJSON is the structure used to store signers information in JSON format
type SignerJSON struct {
	Email string
	Hash  string
}

// JSON is the structure used to store contract information in JSON format
type JSON struct {
	UUID    string
	Date    *time.Time
	Comment string
	File    *FileJSON
	Signers []SignerJSON
}

// GetJSON returns indented json from a contract and some ttp information (nil allowed)
func GetJSON(c *entities.Contract) ([]byte, error) {
	data := JSON{
		UUID:    c.ID.Hex(),
		Date:    &c.Date,
		Comment: c.Comment,
		File: &FileJSON{
			Name:   c.File.Name,
			Hash:   fmt.Sprintf("%x", c.File.Hash),
			Hosted: c.File.Hosted,
		},
		Signers: make([]SignerJSON, len(c.Signers)),
	}

	for i, s := range c.Signers {
		data.Signers[i].Email = s.Email
		data.Signers[i].Hash = fmt.Sprintf("%x", s.Hash)
	}

	return json.MarshalIndent(data, "", "  ")
}
