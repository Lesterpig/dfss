package authority

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTTPHolderNewFile(t *testing.T) {
	// Load from a missing file
	holder, err := NewTTPHolder("holder")
	assert.NotNil(t, holder)
	assert.Nil(t, err)

	ttp := holder.Get()
	assert.Nil(t, ttp)
	assert.Equal(t, 0, holder.Nb())

	holder.Add("1.2.3.4", []byte{0x01, 0xff})
	holder.Add("localhost:9000", []byte{0xaa})
	assert.Equal(t, 2, holder.Nb())

	err = holder.Save("holder")
	assert.Nil(t, err)

	res, _ := ioutil.ReadFile("holder")
	assert.Equal(t, "1.2.3.4 01ff\nlocalhost:9000 aa\n", fmt.Sprintf("%s", res))

	_ = os.Remove("holder")
}

func TestTTPHolderRetrieveFile(t *testing.T) {
	holder, err := NewTTPHolder(filepath.Join("testdata", "ttps"))
	assert.NotNil(t, holder)
	assert.Nil(t, err)
	assert.Equal(t, 2, holder.Nb())

	ttp := holder.Get()
	assert.NotNil(t, ttp)
	assert.Equal(t, "1.2.3.4:9005", ttp.Addrport)
	assert.Equal(t, "aabbcc", fmt.Sprintf("%x", ttp.Hash))

	ttp = holder.Get()
	assert.NotNil(t, ttp)
	assert.Equal(t, "192.168.1.1:36500", ttp.Addrport)
	assert.Equal(t, "0123456789", fmt.Sprintf("%x", ttp.Hash))

	ttp = holder.Get()
	assert.NotNil(t, ttp)
	assert.Equal(t, "1.2.3.4:9005", ttp.Addrport)
	assert.Equal(t, "aabbcc", fmt.Sprintf("%x", ttp.Hash))
}

func TestTTPHolderRetrieveCorruptedFile(t *testing.T) {
	holder, err := NewTTPHolder(filepath.Join("testdata", "corrupted_ttps"))
	assert.Nil(t, holder)
	assert.NotNil(t, err)
}
