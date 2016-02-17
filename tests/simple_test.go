package tests

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/bmizerany/assert"
)

func TestSimple(t *testing.T) {
	// Start the platform
	workingDir, err := ioutil.TempDir("", "dfss_")
	assert.Equal(t, nil, err)
	cmd, ca, err := StartPlatform(workingDir)
	assert.Equal(t, nil, err)

	time.Sleep(3 * time.Second)

	// Start a client
	_, err = CreateClient(workingDir, ca)
	assert.Equal(t, nil, err)

	// Shutdown
	err = cmd.Process.Kill()
	assert.Equal(t, nil, err)

	// Cleanup
	_ = os.RemoveAll(workingDir)
}
