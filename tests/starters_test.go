package tests

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

const testPort = "9090"

var currentClient = 0

// StartPlatform creates root certificate for the platform and starts the platform.
func startPlatform(tmpDir string) (*exec.Cmd, []byte, error) {
	path := filepath.Join(os.Getenv("GOPATH"), "bin", "dfssp")

	// Create temporary directory for platform
	dir, err := ioutil.TempDir(tmpDir, "p_")
	if err != nil {
		return nil, nil, err
	}

	// Init
	cmd := exec.Command(path, "-cn", "localhost", "-path", dir, "-v", "init")
	err = cmd.Run()
	if err != nil {
		return nil, nil, err
	}

	// Get root certificate
	ca, err := ioutil.ReadFile(filepath.Join(dir, "dfssp_rootCA.pem"))
	if err != nil {
		return nil, nil, err
	}

	// Start
	cmd = exec.Command(path, "-db", dbURI, "-path", dir, "-p", testPort, "-v", "start")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()

	return cmd, ca, err
}

// CreateClient creates a new working directory for a client, creating ca.pem.
// It returns a ready-to-run command, but you probably want to set the last argument of the command.
func createClient(tmpDir string, ca []byte) (*exec.Cmd, error) {
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
	cmd := exec.Command(path, "-ca", caPath, "-cert", certPath, "-host", "localhost:"+testPort, "-key", keyPath, "-v")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd, nil
}

func newClient(old *exec.Cmd) *exec.Cmd {
	cmd := exec.Command(old.Path)
	cmd.Args = old.Args
	cmd.Stdout = old.Stdout
	cmd.Stderr = old.Stderr
	return cmd
}

// SetLastArg sets or updates the last argument of a command.
func setLastArg(cmd *exec.Cmd, str string, override bool) {
	if override {
		cmd.Args = cmd.Args[:(len(cmd.Args) - 1)]
	}
	cmd.Args = append(cmd.Args, str)
}