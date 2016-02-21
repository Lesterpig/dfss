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
			Hash:   []byte{0x01, 0x02},
			Hosted: false,
		},
		Signers: []entities.Signer{
			entities.Signer{Email: "a", Hash: []byte{0xaa}},
			entities.Signer{Email: "b", Hash: []byte{0xbb}},
		},
	}

	expected := `{
  "UUID": "00112233445566778899aabb",
  "Date": "2000-01-02T03:04:05.000000006-05:00",
  "Comment": "A test comment\nallow multiline and accents: éÉ",
  "File": {
    "Name": "filename.pdf",
    "Hash": "0102",
    "Hosted": false
  },
  "Signers": [
    {
      "Email": "a",
      "Hash": "aa"
    },
    {
      "Email": "b",
      "Hash": "bb"
    }
  ],
  "Sequence": null,
  "TTP": null
}`

	j, err := GetJSON(c, nil)
	assert.Equal(t, nil, err)
	assert.Equal(t, expected, string(j))

}
