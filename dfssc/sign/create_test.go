package sign

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"dfss/auth"
	"dfss/mockp/server"
	"github.com/bmizerany/assert"
)

var path = filepath.Join("..", "testdata")
var fca = filepath.Join(path, "ca.pem")
var fcert = filepath.Join(path, "cert.pem")
var fkey = filepath.Join(path, "key.pem")
var fcontract = filepath.Join(path, "contract.txt")
var addrPort = "localhost:9090"

func TestMain(m *testing.M) {

	// Load ca and key for platform
	caData, err := ioutil.ReadFile(fca)
	if err != nil {
		os.Exit(1)
	}

	keyData, err := ioutil.ReadFile(filepath.Join(path, "cakey.pem"))
	if err != nil {
		os.Exit(1)
	}

	ca, _ := auth.PEMToCertificate(caData)
	key, _ := auth.PEMToPrivateKey(keyData)

	// Start the platform mock
	go server.Run(ca, key, addrPort)
	time.Sleep(2 * time.Second)

	os.Exit(m.Run())
}

func TestNewCreateManager(t *testing.T) {
	err := NewCreateManager(fca, fcert, fkey, addrPort, "password", fcontract, "success", []string{"a@example.com", "b@example.com"})
	assert.Equal(t, nil, err)

	err = NewCreateManager(fca, fcert, fkey, addrPort, "password", fcontract, "warning", []string{"a@example.com", "b@example.com"})
	assert.Equal(t, "Operation succeeded with a warning message: Some users are not ready yet", err.Error())
}

func TestComputeFile(t *testing.T) {

	m := &CreateManager{filepath: fcontract}
	err := m.computeFile()
	assert.Equal(t, nil, err)
	assert.Equal(t, "37fd29decfb2d689439478b1f64b60441534c1e373a7023676c94ac6772639edab46f80139d167a2741f159e62b3064eca58bb331d32cd10770f29064af2a9de", fmt.Sprintf("%x", m.hash))
	assert.Equal(t, "contract.txt", m.filename)

}
