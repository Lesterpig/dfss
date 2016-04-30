// Package common holds the common functions to be used by other packages.
package common

import (
	"bytes"
	"io/ioutil"
	"os"
)

// SaveToDisk saves the given array of bytes to disk with the given filename
func SaveToDisk(bytes []byte, filename string) error {
	return ioutil.WriteFile(filename, bytes, 0644)
}

// SaveStringToDisk saves the given string to disk with the given filename
func SaveStringToDisk(str, filename string) error {
	buffer := bytes.NewBufferString(str)
	return ioutil.WriteFile(filename, buffer.Bytes(), 0644)
}

// DeleteQuietly try to delete a file, do not fail if an error is raised
func DeleteQuietly(filename string) {
	_ = os.Remove(filename)
}

// ReadFile try to Read a file on the disk
func ReadFile(filename string) ([]byte, error) {
	data, err := ioutil.ReadFile(filename)
	return data, err
}

// FileExists check if a file exists or not
func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}
