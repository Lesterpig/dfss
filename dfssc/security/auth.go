package security

import (
	"crypto/rsa"
	"crypto/x509"
)

// AuthContainer contains common information for TLS authentication.
// Files are not loaded from the beginning, call LoadFiles to load them.
type AuthContainer struct {
	FileCA     string
	FileCert   string
	FileKey    string
	AddrPort   string
	Passphrase string

	CA   *x509.Certificate
	Cert *x509.Certificate
	Key  *rsa.PrivateKey
}

// NewAuthContainer is a shortcut to build an AuthContainer
func NewAuthContainer(fileCA, fileCert, fileKey, addrPort, passphrase string) *AuthContainer {
	return &AuthContainer{
		FileCA:     fileCA,
		FileCert:   fileCert,
		FileKey:    fileKey,
		AddrPort:   addrPort,
		Passphrase: passphrase,
	}
}

// LoadFiles tries to load the required certificates and key for TLS authentication
func (a *AuthContainer) LoadFiles() (ca *x509.Certificate, cert *x509.Certificate, key *rsa.PrivateKey, err error) {
	ca, err = GetCertificate(a.FileCA)
	if err != nil {
		return
	}
	cert, err = GetCertificate(a.FileCert)
	if err != nil {
		return
	}
	key, err = GetPrivateKey(a.FileKey, a.Passphrase)

	a.CA = ca
	a.Cert = cert
	a.Key = key

	return
}
