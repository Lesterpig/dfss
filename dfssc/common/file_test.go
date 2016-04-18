package common

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalDFSSFile(t *testing.T) {
	data, err := ioutil.ReadFile("../testdata/file.dfss")
	assert.Equal(t, nil, err)

	c, err := UnmarshalDFSSFile(data)
	assert.Equal(t, nil, err)
	assert.Equal(t, "filename.pdf", c.File.Name)
	// a tiny test, full unmarshal is tested in dfss/dfssp/contract package
}
