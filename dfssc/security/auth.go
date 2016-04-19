package security

import (
	"crypto/rsa"
	"crypto/x509"

	"github.com/spf13/viper"
)

// AuthContainer contains common information for TLS authentication.
// Files are not loaded from the beginning, call LoadFiles to load them.
type AuthContainer struct {
	Passphrase string

	CA   *x509.Certificate
	Cert *x509.Certificate
	Key  *rsa.PrivateKey
}

// NewAuthContainer is a shortcut to build an AuthContainer
func NewAuthContainer(passphrase string) *AuthContainer {
	return &AuthContainer{
		Passphrase: passphrase,
	}
}

// LoadFiles tries to load the required certificates and key for TLS authentication
func (a *AuthContainer) LoadFiles() (ca *x509.Certificate, cert *x509.Certificate, key *rsa.PrivateKey, err error) {
	ca, err = GetCertificate(viper.GetString("file_ca"))
	if err != nil {
		return
	}
	cert, err = GetCertificate(viper.GetString("file_cert"))
	if err != nil {
		return
	}
	key, err = GetPrivateKey(viper.GetString("file_key"), a.Passphrase)

	a.CA = ca
	a.Cert = cert
	a.Key = key

	return
}
