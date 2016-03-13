package tests

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/bmizerany/assert"
)

// TestRegisterAuth tries to register and auth several users.
//
// GOOD CASES
// client1 : test@example.com with password
// client2 : test@example.com without password
//
// BAD CASES
// client3 : test@example.com
// client4 : wrong mail
// client5 : wrong key size
// client6 : wrong mail during auth
// client7 : wrong token during auth
//
// TODO Add expired accounts test
// TODO Add Stderr test
func TestRegisterAuth(t *testing.T) {
	// Cleanup
	eraseDatabase()

	// Start the platform
	workingDir, err := ioutil.TempDir("", "dfss_")
	assert.Equal(t, nil, err)
	platform, ttp, ca, err := startPlatform(workingDir)
	assert.Equal(t, nil, err)
	defer func() {
		_ = platform.Process.Kill()
		_ = ttp.Process.Kill()
		_ = os.RemoveAll(workingDir)
	}()

	time.Sleep(2 * time.Second)

	// Register client1
	client1, err := createClient(workingDir, ca, 0)
	assert.Equal(t, nil, err)
	err = registerAndAuth(client1, "test@example.com", "password", "", true, true)
	assert.Equal(t, nil, err)

	// Register client2
	client2, err := createClient(workingDir, ca, 0)
	assert.Equal(t, nil, err)
	err = registerAndAuth(client2, "test2@example.com", "", "2048", true, true)
	assert.Equal(t, nil, err)

	// Register client3
	client3, err := createClient(workingDir, ca, 0)
	assert.Equal(t, nil, err)
	err = registerAndAuth(client3, "test@example.com", "", "", true, true)
	assert.NotEqual(t, nil, err)

	// Register client4
	client4, err := createClient(workingDir, ca, 0)
	assert.Equal(t, nil, err)
	err = registerAndAuth(client4, "test wrong mail", "", "", true, true)
	assert.NotEqual(t, nil, err)

	// Register client5
	client5, err := createClient(workingDir, ca, 0)
	assert.Equal(t, nil, err)
	err = registerAndAuth(client5, "wrong@key.fr", "", "1024", true, true)
	assert.NotEqual(t, nil, err)

	// Register client6
	client6, err := createClient(workingDir, ca, 0)
	assert.Equal(t, nil, err)
	err = registerAndAuth(client6, "bad@auth.com", "", "", false, true)
	assert.NotEqual(t, nil, err)

	// Register client7
	client7, err := createClient(workingDir, ca, 0)
	assert.Equal(t, nil, err)
	err = registerAndAuth(client7, "bad@auth2.com", "", "", true, false)
	assert.NotEqual(t, nil, err)
}

func registerAndAuth(client *exec.Cmd, mail, password, keySize string, authMail, authToken bool) error {

	setLastArg(client, "register", false)
	client.Stdin = strings.NewReader(
		mail + "\n" +
			"FR\n" +
			"TEST\n" +
			"TEST\n" +
			keySize + "\n" +
			password + "\n" +
			password + "\n",
	)
	err := client.Run()
	if err != nil {
		return err
	}

	if !authMail { // simulates wrong mail
		mail = "very@badmail.com"
	}

	token := "badToken" // simulates wrong token
	if authToken {
		token = getRegistrationToken(mail)
	}

	// Auth client
	client = newClient(client)
	setLastArg(client, "auth", true)
	client.Stdin = strings.NewReader(
		mail + "\n" +
			token + "\n",
	)
	return client.Run()
}
