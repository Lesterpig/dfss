package sign

import (
	"bytes"
	"crypto/sha512"
	"encoding/hex"
	"io/ioutil"
)

// CheckContractHash computes the hash of the provided file and compares it to the expected one.
func CheckContractHash(filename string, expectedHash string) (ok bool, err error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	expected, err := hex.DecodeString(expectedHash)
	if err != nil {
		return
	}

	hash := sha512.Sum512(data)
	ok = bytes.Equal(expected, hash[:])
	return
}
