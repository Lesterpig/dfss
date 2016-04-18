package tests

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"dfss/dfssp/contract"
	"github.com/stretchr/testify/assert"
)

// TestSignContract unroll the whole signature process.
//
// GOOD CASE
// - Start platform
// - Register client1, client2 and client3
// - Create contract `contract.txt`
// - Sign it
// - Check if all proof files are present
//
// TODO BAD CASES

func TestSignContract(t *testing.T) {
	// Cleanup
	eraseDatabase()

	// Start the platform
	workingDir, err := ioutil.TempDir("", "dfss_")
	assert.Equal(t, nil, err)
	_, _, _, stop, ca, err := startPlatform(workingDir)
	assert.Equal(t, nil, err)
	defer stop()

	time.Sleep(2 * time.Second)

	// Register clients
	clients := make([]*exec.Cmd, 3)
	client1, err := createClient(workingDir, ca, 9091)
	assert.Equal(t, nil, err)
	err = registerAndAuth(client1, "client1@example.com", "password", "", true, true)
	assert.Equal(t, nil, err)
	client2, err := createClient(workingDir, ca, 9092)
	assert.Equal(t, nil, err)
	err = registerAndAuth(client2, "client2@example.com", "password", "", true, true)
	assert.Equal(t, nil, err)
	client3, err := createClient(workingDir, ca, 9093)
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

	closeChannel := make(chan []byte, 3)
	for i := 0; i < 3; i++ {
		setLastArg(clients[i], "sign", true)
		setLastArg(clients[i], contractPath, false)
		go func(c *exec.Cmd, i int) {
			time.Sleep(time.Duration(i*2) * time.Second)
			c.Stdin = strings.NewReader("password\nyes\n")
			c.Stderr = os.Stderr
			output, err1 := c.Output()
			if err1 != nil {
				output = nil
			}
			closeChannel <- output
		}(clients[i], i)
	}

	regexes := []*regexp.Regexp{
		regexp.MustCompile(`Everybody is ready, starting the signature [a-f0-9]+`),
		regexp.MustCompile(`Do you REALLY want to sign contract\.txt\? Type 'yes' to confirm:`),
	}
	for i := 0; i < 3; i++ {
		output := <-closeChannel
		assert.NotEqual(t, nil, output, "The return error should be null")
		for _, r := range regexes {
			assert.True(t, r.Match(output), "Regex is not satisfied: ", r.String())
		}
	}

	// Ensure that all the files are present
	proofFile := regexp.MustCompile(`client[0-9]+@example.com.*\.proof`)
	files, _ := ioutil.ReadDir("./")

	matches := 0
	for _, file := range files {
		if proofFile.Match([]byte(file.Name())) {
			matches++
			err = os.Remove("./" + file.Name())
			assert.True(t, err == nil, "Cannot remove .proof matching file")
		}
	}
	assert.True(t, matches == 3, "Missing proof file ?")

	time.Sleep(time.Second)
}
