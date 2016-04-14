package sign

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/bmizerany/assert"
)

func TestFetchContract(t *testing.T) {
	checkFetchResult(t, "01", false, "Hello")
}

func TestFetchContractWrongUUID(t *testing.T) {
	checkFetchResult(t, "02", true, "")
}

func checkFetchResult(t *testing.T, uuid string, errExpected bool, content string) {
	file, _ := ioutil.TempFile("", "")
	defer func() { _ = os.Remove(file.Name()) }()
	err := FetchContract(fca, fcert, fkey, addrPort, "password", uuid, file.Name())
	if errExpected {
		assert.NotEqual(t, nil, err)
	} else {
		assert.Equal(t, nil, err)
	}

	data, _ := ioutil.ReadFile(file.Name())
	assert.Equal(t, content, fmt.Sprintf("%s", data))
}
