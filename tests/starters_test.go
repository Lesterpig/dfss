package tests

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const testPort = "9090"

var currentClient = 0

// startPlatform creates root certificate for the platform and the TTP, and starts both modules
func startPlatform(tmpDir string) (platform, ttp, demo *exec.Cmd, stop func(), ca []byte, err error) {
	path := filepath.Join(os.Getenv("GOPATH"), "bin", "dfssp")
	ttpPath := filepath.Join(os.Getenv("GOPATH"), "bin", "dfsst")
	demoPath := filepath.Join(os.Getenv("GOPATH"), "bin", "dfssd")

	// Create temporary directory for platform
	dir, err := ioutil.TempDir(tmpDir, "p_")
	if err != nil {
		return
	}

	// Init
	cmd := exec.Command(path, "--path", dir, "-v", "init")
	ttpsPath := filepath.Join(dir, "ttps")
	err = cmd.Run()
	if err != nil {
		return
	}

	// Create TTP working directory
	cmd = exec.Command(path, "--path", dir, "--ttps", ttpsPath, "--addr", "localhost:9098", "ttp")
	err = cmd.Run()
	if err != nil {
		return
	}

	// Get root certificate
	ca, err = ioutil.ReadFile(filepath.Join(dir, "dfssp_rootCA.pem"))
	if err != nil {
		return
	}

	// Start demonstrator
	demo = exec.Command(demoPath, "--port", "9099", "nogui")
	demo.Stdout = os.Stdout
	demo.Stderr = os.Stderr
	err = demo.Start()
	time.Sleep(time.Second)

	// Start platform
	platform = exec.Command(path, "--db", dbURI, "--path", dir, "-p", testPort, "--ttps", ttpsPath, "-d", "localhost:9099", "start")
	platform.Stdout = os.Stdout
	platform.Stderr = os.Stderr
	err = platform.Start()

	// Start TTP
	ttp = exec.Command(ttpPath, "--db", dbURI, "--port", "9098", "-d", "localhost:9099", "start")
	ttp.Dir = filepath.Join(dir, "ttp")
	ttp.Stdout = os.Stdout
	ttp.Stderr = os.Stderr
	_ = ioutil.WriteFile(filepath.Join(ttp.Dir, "ca.pem"), ca, 0600)
	err = ttp.Start()

	stop = func() {
		_ = platform.Process.Kill()
		_ = ttp.Process.Kill()
		_ = os.RemoveAll(tmpDir)
		time.Sleep(2 * time.Second)
		_ = demo.Process.Kill()
	}

	return
}

// createClient creates a new working directory for a client, creating ca.pem.
// It returns a ready-to-run command, but you probably want to set the last argument of the command.
func createClient(tmpDir string, ca []byte, port int) (*exec.Cmd, error) {
	path := filepath.Join(os.Getenv("GOPATH"), "bin", "dfssc")

	// Create temporary directory for client
	clientID := strconv.Itoa(currentClient)
	currentClient++
	dir, err := ioutil.TempDir(tmpDir, "c"+clientID+"_")
	if err != nil {
		return nil, err
	}

	caPath := filepath.Join(dir, "ca.pem")
	certPath := filepath.Join(dir, "cert.pem")
	keyPath := filepath.Join(dir, "key.pem")

	// Write the CA
	err = ioutil.WriteFile(caPath, ca, os.ModePerm)
	if err != nil {
		return nil, err
	}

	// Prepare the client command.
	// The last argument is up to you!
	cmd := exec.Command(path, "--ca", caPath, "--cert", certPath, "--host", "127.0.0.1:"+testPort, "--key", keyPath, "--port", strconv.Itoa(port), "--timeout", "3s", "-d", "localhost:9099")

	return cmd, nil
}

// newClient clones the current command to another one.
// It's very useful when doing several commands on the same client.
func newClient(old *exec.Cmd) *exec.Cmd {
	cmd := exec.Command(old.Path)
	cmd.Args = old.Args
	cmd.Stdout = old.Stdout
	cmd.Stderr = old.Stderr
	return cmd
}

// setLastArg sets or updates the last argument of a command.
func setLastArg(cmd *exec.Cmd, str string, override bool) {
	if override {
		cmd.Args = cmd.Args[:(len(cmd.Args) - 1)]
	}
	cmd.Args = append(cmd.Args, str)
}

// checkStderr runs the provided command and compares the stderr output with the given one.
// It returns the value of cmd.Wait()
func checkStderr(t *testing.T, cmd *exec.Cmd, value string) error {
	cmd.Stderr = nil
	stderr, _ := cmd.StderrPipe()
	err := cmd.Start()
	assert.Equal(t, nil, err)

	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(stderr)
	s := buf.String()
	assert.Equal(t, value, s)

	return cmd.Wait()
}
