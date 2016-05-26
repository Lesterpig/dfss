// Package security is responsible for generating keys and certificate requests
package security

import (
	"crypto/rsa"
	"fmt"

	"dfss/auth"
	"dfss/dfssc/common"
	"github.com/spf13/viper"
)

// GenerateKeys generate a pair of keys and save it to the disk
func GenerateKeys(bits int, passphrase string) (*rsa.PrivateKey, error) {
	key, err := auth.GeneratePrivateKey(bits)
	if err != nil {
		return nil, err
	}

	pem, err := auth.PrivateKeyToEncryptedPEM(key, passphrase)
	if err != nil {
		return nil, err
	}

	err = common.SaveToDisk(pem, viper.GetString("file_key"))
	if err != nil {
		return nil, err
	}

	return key, nil
}

// GenerateCertificateRequest generate a certificate request from data, and
// return a PEM-encoded certificate as a string
func GenerateCertificateRequest(country, organization, unit, mail string, key *rsa.PrivateKey) (string, error) {
	data, err := auth.GetCertificateRequest(country, organization, unit, mail, key)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s", data), nil

}

// SaveCertificate saves a PEM-encoded certificate on disk
func SaveCertificate(cert, filename string) error {
	return common.SaveStringToDisk(cert, filename)
}
