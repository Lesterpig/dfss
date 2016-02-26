package tests

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"dfss/dfssp/contract"
	"github.com/bmizerany/assert"
)

func TestSignContract(t *testing.T) {
	// Cleanup
	eraseDatabase()

	// Start the platform
	workingDir, err := ioutil.TempDir("", "dfss_")
	assert.Equal(t, nil, err)
	platform, ca, err := startPlatform(workingDir)
	assert.Equal(t, nil, err)
	defer func() {
		_ = platform.Process.Kill()
		_ = os.RemoveAll(workingDir)
	}()

	time.Sleep(2 * time.Second)

	// Register clients
	clients := make([]*exec.Cmd, 3)
	client1, err := createClient(workingDir, ca)
	assert.Equal(t, nil, err)
	err = registerAndAuth(client1, "client1@example.com", "password", "", true, true)
	assert.Equal(t, nil, err)
	client2, err := createClient(workingDir, ca)
	assert.Equal(t, nil, err)
	err = registerAndAuth(client2, "client2@example.com", "password", "", true, true)
	assert.Equal(t, nil, err)
	client3, err := createClient(workingDir, ca)
	assert.Equal(t, nil, err)
	err = registerAndAuth(client3, "client3@example.com", "password", "", true, true)
	assert.Equal(t, nil, err)

	// Create contract
	client1 = newClient(client1)
	setLastArg(client1, "new", true)
	client1.Stdin = strings.NewReader(
		"password\n" +
			filepath.Join("testdata", "contract.txt") + "\n" +
			"\n" +
			"client1@example.com\n" +
			"client2@example.com\n" +
			"client3@example.com\n" +
			"\n",
	)
	err = checkStderr(t, client1, "")
	assert.Equal(t, nil, err)

	// Get contract file
	contractEntity := getContract("contract.txt", 0)
	contractData, err := contract.GetJSON(contractEntity, nil)
	assert.Equal(t, nil, err)
	contractPath := filepath.Join(workingDir, "c.dfss")
	err = ioutil.WriteFile(contractPath, contractData, os.ModePerm)
	assert.Equal(t, nil, err)

	// Sign!
	clients[0] = newClient(client1)
	clients[1] = newClient(client2)
	clients[2] = newClient(client3)
	closeChannel := make(chan error, 3)
	for i := 0; i < 3; i++ {
		setLastArg(clients[i], "sign", true)
		setLastArg(clients[i], contractPath, false)
		go func(c *exec.Cmd, i int) {
			time.Sleep(time.Duration(i*2) * time.Second)
			c.Stdin = strings.NewReader("password\nyes\n")
			closeChannel <- c.Run()
		}(clients[i], i)
	}

	for i := 0; i < 3; i++ {
		assert.Equal(t, nil, <-closeChannel)
	}
}
