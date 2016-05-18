package auth

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"fmt"
)

// SignStructure signs the provided structure with the private key.
// The used protocol is RSA PKCS#1 v1.5 with SHA-512 hash.
// The structure is serialized to a string representation using the fmt package.
func SignStructure(key *rsa.PrivateKey, structure interface{}) ([]byte, error) {
	hash, err := hashStruct(structure)
	if err != nil {
		return nil, err
	}

	return rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA512, hash)
}

// VerifyStructure verifies the signed message according to the provided structure and certificate.
// See SignStructure for protocol definition.
func VerifyStructure(cert *x509.Certificate, structure interface{}, signed []byte) (bool, error) {
	hash, err := hashStruct(structure)
	if err != nil {
		return false, err
	}

	err = rsa.VerifyPKCS1v15(cert.PublicKey.(*rsa.PublicKey), crypto.SHA512, hash, signed)
	return err == nil, err
}

func hashStruct(structure interface{}) (hash []byte, err error) {
	data := []byte(fmt.Sprintf("%v", structure))
	rawHash := sha512.Sum512(data)
	hash = rawHash[:]
	return
}
