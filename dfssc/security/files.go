package security

import (
	"crypto/rsa"
	"crypto/x509"
	"dfss/auth"
	"dfss/dfssc/common"
)

// GetCertificate return the Certificate stored on the disk
func GetCertificate(filename string) (*x509.Certificate, error) {
	data, err := common.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	cert, err := auth.PEMToCertificate(data)
	if err != nil {
		return nil, err
	}

	return cert, nil
}

// GetPrivateKey return the private key stored on the disk
func GetPrivateKey(filename, passphrase string) (*rsa.PrivateKey, error) {

	data, err := common.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	key, err := auth.EncryptedPEMToPrivateKey(data, passphrase)
	if err != nil {
		return nil, err
	}

	return key, nil
}
