package tests

import (
	"dfss/dfssc/common"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bmizerany/assert"
)

// TestExport tries to export the certificate and pricate key of the client
//
// CLASSIC SCENARIO
// 1. Export is node successfully
//
// BAD CASES
// 1. Wrong passphrase for unlocking private key
// 2. Missing certificate for the client
func TestExport(t *testing.T) {
	// Initialize the directory and files
	workingDir, err := ioutil.TempDir("", "dfss_")
	assert.Equal(t, nil, err)

	certPath := filepath.Join(workingDir, "cert.pem")
	certFixture, _ := common.ReadFile(filepath.Join(os.Getenv("GOPATH"), "src", "dfss", "tests", "testdata", "cert.pem"))
	err = common.SaveToDisk(certFixture, certPath)
	assert.Equal(t, nil, err)

	keyPath := filepath.Join(workingDir, "key.pem")
	keyFixture, _ := common.ReadFile(filepath.Join(os.Getenv("GOPATH"), "src", "dfss", "tests", "testdata", "key.pem"))
	err = common.SaveToDisk(keyFixture, keyPath)

	confPath := filepath.Join(workingDir, "dfssc.conf")
	common.DeleteQuietly(confPath)
	assert.Equal(t, nil, err)
	path := filepath.Join(os.Getenv("GOPATH"), "bin", "dfssc")

	// Basic command
	cmd := exec.Command(path, "-cert", certPath, "-key", keyPath, "export", confPath)

	// Export the configuration
	cmd.Stdin = strings.NewReader(
		"pass\n" +
			"password\n",
	)
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	assert.Equal(t, nil, err)
	assert.T(t, common.FileExists(confPath))

	// Bad case 1 : Wrong passphrase for private key
	badCmd1 := exec.Command(path, "-cert", certPath, "-key", keyPath, "export", confPath)
	common.DeleteQuietly(confPath)
	badCmd1.Stdin = strings.NewReader(
		"passphrase\n" +
			"wrong passphrase\n",
	)
	_ = badCmd1.Run()
	// assert.Equalf(t, nil, err, "%x", err)
	assert.T(t, !common.FileExists(confPath))

	// Bad case 2 : Missing certificate
	badCmd2 := exec.Command(path, "-cert", certPath, "-key", keyPath, "export", confPath)
	common.DeleteQuietly(certPath)
	badCmd2.Stdin = strings.NewReader(
		"passphrase\n" +
			"password\n",
	)
	_ = badCmd2.Run()
	// assert.Equal(t, nil, err)
	assert.T(t, !common.FileExists(confPath))
}

// TestImport tries to import the certificate and private key of a user
//
// CLASSIC SCENARIO
// Import is done without problem
//
// BAD CASES
// 1. There is already a private key
// 2. Wrong passphrase for the configuration
// 3. Wrong passphrase for the private key
func TestImport(t *testing.T) {
	// Initialize the directory and files
	workingDir, err := ioutil.TempDir("", "dfss_")
	assert.Equal(t, nil, err)

	certPath := filepath.Join(workingDir, "cert.pem")
	certFixture, _ := common.ReadFile(filepath.Join(os.Getenv("GOPATH"), "src", "dfss", "tests", "testdata", "cert.pem"))
	err = common.SaveToDisk(certFixture, certPath)
	assert.Equal(t, nil, err)

	keyPath := filepath.Join(workingDir, "key.pem")
	keyFixture, _ := common.ReadFile(filepath.Join(os.Getenv("GOPATH"), "src", "dfss", "tests", "testdata", "key.pem"))
	err = common.SaveToDisk(keyFixture, keyPath)
	assert.Equal(t, nil, err)

	confPath := filepath.Join(workingDir, "dfssc.conf")
	common.DeleteQuietly(confPath)
	assert.Equal(t, nil, err)
	path := filepath.Join(os.Getenv("GOPATH"), "bin", "dfssc")

	// Create the config file
	cmd := exec.Command(path, "-cert", certPath, "-key", keyPath, "export", confPath)

	cmd.Stdin = strings.NewReader(
		"pass\n" +
			"password\n",
	)
	err = cmd.Run()
	assert.Equal(t, nil, err)
	assert.T(t, common.FileExists(confPath))

	// Nominal case
	common.DeleteQuietly(certPath)
	common.DeleteQuietly(keyPath)

	cmd = exec.Command(path, "import", confPath)
	cmd.Stdin = strings.NewReader(
		"pass\n" +
			"password\n",
	)
	err = cmd.Run()
	assert.Equal(t, nil, err)
	assert.T(t, common.FileExists(certPath))
	assert.T(t, common.FileExists(keyPath))

	// Bad case 1 : There is already the key file
	common.DeleteQuietly(certPath)

	badCmd1 := exec.Command(path, "import", confPath)
	badCmd1.Stdin = strings.NewReader(
		"pass\n" +
			"password\n",
	)
	_ = badCmd1.Run()
	// assert.Equalf(t, nil, err, "%x", err)
	assert.T(t, !common.FileExists(certPath))

	// Bad case 2 : Wrong passphrase of the configuration
	common.DeleteQuietly(keyPath)

	badCmd2 := exec.Command(path, "import", confPath)
	badCmd2.Stdin = strings.NewReader(
		"I am a wrong passphrase\n" +
			"password\n",
	)
	_ = badCmd2.Run()
	// assert.Equal(t, nil, err)
	assert.T(t, !common.FileExists(certPath))
	assert.T(t, !common.FileExists(keyPath))

	// Bad case 3 : Wrong passphrase for the private key
	badCmd3 := exec.Command(path, "import", confPath)
	badCmd3.Stdin = strings.NewReader(
		"\n" +
			"I am a wrong passphrase\n",
	)
	_ = badCmd3.Run()
	// assert.Equal(t, nil, err)
	assert.T(t, !common.FileExists(certPath))
	assert.T(t, !common.FileExists(keyPath))
}
