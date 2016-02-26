// Package auth provides simple ways to handle user authentication
package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

// Cipher used to encrypt private key content in PEM format
var Cipher = x509.PEMCipherAES256

// GeneratePrivateKey builds a private key of given size from default random
func GeneratePrivateKey(bits int) (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, bits)
}

// PrivateKeyToEncryptedPEM builds a PEM-encoded array of bytes from a private key and a password.
// If pwd is empty, then the resulting PEM will not be encrypted.
func PrivateKeyToEncryptedPEM(key *rsa.PrivateKey, pwd string) ([]byte, error) {
	var err error

	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	if pwd != "" {
		block, err = x509.EncryptPEMBlock(rand.Reader, block.Type, block.Bytes, []byte(pwd), Cipher)
		if err != nil {
			return nil, err
		}
	}

	return pem.EncodeToMemory(block), nil
}

// PrivateKeyToPEM produces a unencrypted PEM-encoded array of bytes from a private key.
func PrivateKeyToPEM(key *rsa.PrivateKey) []byte {
	p, _ := PrivateKeyToEncryptedPEM(key, "")
	return p
}

// EncryptedPEMToPrivateKey tries to decrypt and decode a PEM-encoded array of bytes to a private key.
// If pwd is empty, then the function will not try to decrypt the PEM block.
//
// In case of wrong password, the returned error will be equals to x509.IncorrectPasswordError
func EncryptedPEMToPrivateKey(data []byte, pwd string) (*rsa.PrivateKey, error) {
	var err error

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("Data is not a valid pem-encoding")
	}
	decodedData := block.Bytes

	if pwd != "" {
		decodedData, err = x509.DecryptPEMBlock(block, []byte(pwd))

		if err != nil {
			return nil, err
		}
	}

	return x509.ParsePKCS1PrivateKey(decodedData)

}

// PEMToPrivateKey tries to decode a plain PEM-encoded array of bytes to a private key.
func PEMToPrivateKey(data []byte) (*rsa.PrivateKey, error) {
	return EncryptedPEMToPrivateKey(data, "")
}

// IsPEMEncrypted tests whether a PEM-encoded array of bytes is encrypted or not.
func IsPEMEncrypted(data []byte) bool {
	var block, _ = pem.Decode(data)
	return x509.IsEncryptedPEMBlock(block)
}
