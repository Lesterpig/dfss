package contract

import (
	"testing"
	"time"

	"dfss/dfssp/entities"
	"github.com/bmizerany/assert"
	"gopkg.in/mgo.v2/bson"
)

func TestGetJSON(t *testing.T) {

	location, _ := time.LoadLocation("EST")

	c := &entities.Contract{
		ID:   bson.ObjectIdHex("00112233445566778899aabb"),
		Date: time.Date(2000, 1, 2, 3, 4, 5, 6, location),
		Comment: `A test comment
allow multiline and accents: éÉ`,
		File: &entities.File{
			Name:   "filename.pdf",
			Hash:   "hash",
			Hosted: false,
		},
		Signers: []entities.Signer{
			entities.Signer{Email: "a", Hash: "ha"},
			entities.Signer{Email: "b", Hash: "hb"},
		},
	}

	expected := `{
  "UUID": "00112233445566778899aabb",
  "Date": "2000-01-02T03:04:05.000000006-05:00",
  "Comment": "A test comment\nallow multiline and accents: éÉ",
  "File": {
    "Name": "filename.pdf",
    "Hash": "hash",
    "Hosted": false
  },
  "Signers": [
    {
      "Email": "a",
      "Hash": "ha"
    },
    {
      "Email": "b",
      "Hash": "hb"
    }
  ],
  "Sequence": null,
  "TTP": null
}`

	j, err := GetJSON(c, nil)
	assert.Equal(t, nil, err)
	assert.Equal(t, expected, string(j))

}
