package tests

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bmizerany/assert"
)

// TestNewContract tries to creates new contracts on the platform.
//
// CLASSIC SCENARIO
// 1. client1 registers on the platform
// 2. client1 sends a new contract on the platform, but client2 is not here yet
// 3. client2 registers on the platform
// 4. client2 sends a new contract on the platform, and everyone is here
//
// BAD CASES
// 1. client2 sends a new contract with a wrong password
// 2. client3 sends a new contract without authentication
// 3. client1 sends a new contract with an invalid filepath
func TestNewContract(t *testing.T) {
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

	// Register client1
	client1, err := createClient(workingDir, ca)
	assert.Equal(t, nil, err)
	err = registerAndAuth(client1, "client1@example.com", "password", "", true, true)
	assert.Equal(t, nil, err)

	// Create contract
	client1 = newClient(client1)
	setLastArg(client1, "new", true)
	client1.Stdin = strings.NewReader(
		"password\n" +
			filepath.Join("testdata", "contract.txt") + "\n" +
			"A very nice comment\n" +
			"client1@example.com\n" +
			"client2@example.com\n" +
			"\n",
	)
	err = checkStderr(t, client1, "Operation succeeded with a warning message: Some users are not ready yet\n")
	assert.NotEqual(t, nil, err)

	// Check database
	contract := getContract("contract.txt", 0)
	assert.Equal(t, false, contract.Ready)
	assert.Equal(t, "A very nice comment", contract.Comment)
	assert.Equal(t, "6a95f6bcd6282186a7b1175fbaab4809ca5f665f7c4d55675de2399c83e67252069d741a88c766b1a79206d6dfbd5552cd7f9bc69b43bee161d1337228b4a4a8", fmt.Sprintf("%x", contract.File.Hash))
	assert.Equal(t, 2, len(contract.Signers))
	assert.Equal(t, "client1@example.com", contract.Signers[0].Email)
	assert.Equal(t, "client2@example.com", contract.Signers[1].Email)
	assert.T(t, len(contract.Signers[0].Hash) > 0)
	assert.T(t, len(contract.Signers[1].Hash) == 0)

	// Register second signer
	client2, err := createClient(workingDir, ca)
	assert.Equal(t, nil, err)
	err = registerAndAuth(client2, "client2@example.com", "password2", "", true, true)
	assert.Equal(t, nil, err)

	// Check database²
	contract = getContract("contract.txt", 0)
	assert.Equal(t, true, contract.Ready)
	assert.T(t, len(contract.Signers[0].Hash) > 0)
	assert.T(t, len(contract.Signers[1].Hash) > 0)

	// Create a second contract
	client2 = newClient(client2)
	setLastArg(client2, "new", true)
	client2.Stdin = strings.NewReader(
		"password2\n" +
			filepath.Join("testdata", "contract.txt") + "\n" +
			"Another comment with some accents héhé\n" +
			"client1@example.com\n" +
			"client2@example.com\n" +
			"\n",
	)
	err = checkStderr(t, client2, "")
	assert.Equal(t, nil, err)

	// Check database³
	contract = getContract("contract.txt", 1)
	assert.Equal(t, true, contract.Ready)
	assert.Equal(t, "Another comment with some accents héhé", contract.Comment)
	assert.T(t, len(contract.Signers[0].Hash) > 0)
	assert.T(t, len(contract.Signers[1].Hash) > 0)

	// Bad case: wrong password
	client2 = newClient(client2)
	setLastArg(client2, "new", true)
	client2.Stdin = strings.NewReader(
		"wrongPwd\n" +
			filepath.Join("testdata", "contract.txt") + "\n" +
			"\n" +
			"client1@example.com\n" +
			"client2@example.com\n" +
			"\n",
	)
	err = checkStderr(t, client2, "x509: decryption password incorrect\n")
	assert.NotEqual(t, nil, err)

	// Bad case: no authentication
	client3, err := createClient(workingDir, ca)
	setLastArg(client3, "new", false)
	client3.Stdin = strings.NewReader(
		"\n" +
			filepath.Join("testdata", "contract.txt") + "\n" +
			"\n" +
			"client1@example.com\n" +
			"\n",
	)
	err = client3.Run()
	assert.NotEqual(t, nil, err)

	// Bad case: bad filepath
	client1 = newClient(client1)
	setLastArg(client1, "new", true)
	client1.Stdin = strings.NewReader(
		"password\n" +
			"invalidFile\n" +
			"client1@example.com\n" +
			"\n",
	)
	err = checkStderr(t, client1, "open invalidFile: no such file or directory\n")
	assert.NotEqual(t, nil, err)

	// Check number of stored contracts
	assert.Equal(t, 2, dbManager.Get("contracts").Count())
}
