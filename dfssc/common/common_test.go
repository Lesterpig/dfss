package common

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/bmizerany/assert"
)

var path = os.TempDir()
var ffoo = filepath.Join(path, "foo.txt")
var fbar = filepath.Join(path, "bar.txt")
var fbaz = filepath.Join(path, "baz.txt")
var fqux = filepath.Join(path, "qux.txt")

func TestMain(m *testing.M) {

	// Setup

	// Run tests
	code := m.Run()

	// Teardown
	DeleteQuietly(ffoo)
	DeleteQuietly(fbar)
	DeleteQuietly(fbaz)
	DeleteQuietly(fqux)

	os.Exit(code)
}

// Save an array of bytes on the disk
func TestSaveToDisk(t *testing.T) {
	data := make([]byte, 3)
	data[0] = 'f'
	data[1] = 'o'
	data[2] = 'o'
	err := SaveToDisk(data, ffoo)
	assert.T(t, err == nil, "An error has been raised")

	assert.T(t, FileExists(ffoo), "foo.txt should be present")
}

// Save a string on the disk
func TestSaveStringToDisk(t *testing.T) {
	s := "bar"
	err := SaveStringToDisk(s, fbar)
	assert.T(t, err == nil, "An error has been raised")

	assert.T(t, FileExists(fbar), "bar.txt should be present")
}

// DeleteQuietly should never raise an error, even with non-existant file
func TestDeleteQuietly(t *testing.T) {
	s := "baz"
	_ = SaveStringToDisk(s, fbaz)
	assert.T(t, FileExists(fbaz), "baz.txt should be present")
	DeleteQuietly(fbaz)
	assert.T(t, !FileExists(fbaz), "baz.txt should not be present")

	// Assert it does not panic when deleting an inexistant file
	DeleteQuietly("dummy")
}

// Test the reading of a file
func TestReadFile(t *testing.T) {
	s := "qux"
	_ = SaveStringToDisk(s, fqux)
	assert.T(t, FileExists(fqux), "qux.txt should be present")

	data, err := ReadFile(fqux)
	if err != nil {
		fmt.Println(err.Error())
		assert.T(t, err == nil, "An error has been raised while reading the file")
	}
	assert.Equal(t, s, fmt.Sprintf("%s", data), "Expected qux, received ", fmt.Sprintf("%s", data))

}
