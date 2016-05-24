package tests

import (
	"bufio"
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

// setupSignature prepares required servers and clients to sign a contract.
// - Start platform, ttp, demonstrator
// - Register client1, client2 and client3
// - Create contract `contract.txt`
func setupSignature(t *testing.T) (stop func(), clients []*exec.Cmd, contractPath, contractFilePath string) {
	// Cleanup
	eraseDatabase()

	// Start the platform
	workingDir, err := ioutil.TempDir("", "dfss_")
	assert.Equal(t, nil, err)
	_, _, _, stop, ca, err := startPlatform(workingDir)
	assert.Equal(t, nil, err)

	time.Sleep(2 * time.Second)

	// Register clients
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
	contractFilePath = filepath.Join("testdata", "contract.txt")
	client1.Stdin = strings.NewReader(
		"password\n" +
			contractFilePath + "\n" +
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
	contractData, err := contract.GetJSON(contractEntity)
	assert.Equal(t, nil, err)
	contractPath = filepath.Join(workingDir, "c.dfss")
	err = ioutil.WriteFile(contractPath, contractData, os.ModePerm)
	assert.Equal(t, nil, err)

	// Test with wrong file, should abort
	wrongFileClient := newClient(client1)
	setLastArg(wrongFileClient, "sign", true)
	setLastArg(wrongFileClient, contractPath, false)
	wrongFileClient.Stdin = strings.NewReader("wrongfile.txt\npassword\nyes\n")
	_, err = wrongFileClient.Output()
	assert.NotNil(t, err)

	clients = make([]*exec.Cmd, 3)
	clients[0] = newClient(client1)
	clients[1] = newClient(client2)
	clients[2] = newClient(client3)
	return
}

// TestSignContract unroll the whole signature process.
// In this test, everything should work fine without any ttp call.
func TestSignContract(t *testing.T) {
	// Setup
	stop, clients, contractPath, contractFilePath := setupSignature(t)
	defer stop()

	// Sign!
	closeChannel := make(chan []byte, 3)
	for i := 0; i < 3; i++ {
		setLastArg(clients[i], "sign", true)
		setLastArg(clients[i], contractPath, false)
		go func(c *exec.Cmd, i int) {
			time.Sleep(time.Duration(i*2) * time.Second)
			c.Stdin = strings.NewReader(contractFilePath + "\npassword\nyes\n")
			c.Stderr = os.Stderr
			output, err := c.Output()
			if err != nil {
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

	checkProofFile(t, 3)
	time.Sleep(time.Second)
}

// checkProofFile counts the number of proof file contained in the current directory, and compares it to the nb parameter.
func checkProofFile(t *testing.T, nb int) {
	// Ensure that all the files are present
	proofFile := regexp.MustCompile(`client[0-9]+@example.com.*\.proof`)
	files, _ := ioutil.ReadDir("./")

	matches := 0
	for _, file := range files {
		if proofFile.Match([]byte(file.Name())) {
			matches++
			err := os.Remove("./" + file.Name())
			assert.True(t, err == nil, "Cannot remove .proof matching file")
		}
	}
	assert.Equal(t, nb, matches, "Invalid number of proof file(s)")
}

// TestSignContractFailure tests the signature with a faulty client, when contract can't be generated.
// In this test, everything should not work fine, because client3 shutdowns way too early.
func TestSignContractFailure(t *testing.T) {
	signatureHelper(t, "1", 0)
}

// TestSignContractSuccess tests the signature with a faulty client, when contract can be generated.
// In this test, everything should not work fine, because client3 shutdowns way too early.
func TestSignContractSuccess(t *testing.T) {
	signatureHelper(t, "2", 2)
}

// signatureHelper : launches a parametrized signature, with the number of rounds a client will accomplish before shutting down,
// and the number of proof files expected to be generated.
func signatureHelper(t *testing.T, round string, nbFiles int) {
	// Setup
	stop, clients, contractPath, contractFilePath := setupSignature(t)
	defer stop()

	// Configure client3 to be faulty
	setLastArg(clients[2], "--stopbefore", true)
	setLastArg(clients[2], round, false)
	setLastArg(clients[2], "sign", false)

	// Sign!
	closeChannel := make(chan []byte, 3)
	for i := 0; i < 3; i++ {
		setLastArg(clients[i], "sign", true)
		setLastArg(clients[i], contractPath, false)
		go func(c *exec.Cmd, i int) {
			c.Stdin = strings.NewReader(contractFilePath + "\npassword\nyes\n")
			c.Stderr = bufio.NewWriter(os.Stdout)
			output, _ := c.Output()
			closeChannel <- output
		}(clients[i], i)
	}

	for i := 0; i < 3; i++ {
		// TODO check stderr?
		<-closeChannel
	}

	checkProofFile(t, nbFiles)
	time.Sleep(time.Second)
}
