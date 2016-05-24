// Package authority creates and manages platform certificates.
package authority

import (
	"crypto/rsa"
	"crypto/x509"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/viper"

	"dfss/auth"
)

// PlatformID contains platform private key and root certificate
type PlatformID struct {
	Pkey   *rsa.PrivateKey
	RootCA *x509.Certificate
}

// Initialize creates and saves the platform's private key and root certificate to a PEM format.
// If ca and rKey are not nil, they will be used as the root certificate and root private key instead of creating a ones.
// The files are saved at the specified path by viper.
// The returned `hash` is the SHA-512 hash of the generated certificate.
func Initialize(v *viper.Viper, ca *x509.Certificate, rKey *rsa.PrivateKey) (hash []byte, err error) {
	// Generate the private key.
	key, err := auth.GeneratePrivateKey(v.GetInt("key_size"))

	if err != nil {
		return nil, err
	}

	var cert []byte
	path := v.GetString("path")
	// ca_filename and pkey_filename are part of the global conf of
	// the app, that's why we don't fetch them from the local viper
	certPath := filepath.Join(path, viper.GetString("ca_filename"))
	keyPath := filepath.Join(path, viper.GetString("pkey_filename"))

	if ca == nil {
		// Generate the root certificate, using the private key.
		cert, err = auth.GetSelfSignedCertificate(v.GetInt("validity"), auth.GenerateUID(), v.GetString("country"), v.GetString("organization"), v.GetString("unit"), v.GetString("cn"), key)
	} else {
		csr, _ := auth.GetCertificateRequest(v.GetString("country"), v.GetString("organization"), v.GetString("unit"), v.GetString("cn"), key)
		request, _ := auth.PEMToCertificateRequest(csr)
		cert, err = auth.GetCertificate(v.GetInt("validity"), auth.GenerateUID(), request, ca, rKey)
		// Override default path values
		certPath = filepath.Join(path, "cert.pem")
		keyPath = filepath.Join(path, "key.pem")
	}

	if err != nil {
		return nil, err
	}

	// Create missing folders, if needed
	err = os.MkdirAll(path, os.ModeDir|0700)
	if err != nil {
		return nil, err
	}

	// Convert the private key to a PEM format, and save it.
	keyPem := auth.PrivateKeyToPEM(key)
	err = ioutil.WriteFile(keyPath, keyPem, 0600)
	if err != nil {
		return nil, err
	}

	// Save the root certificate.
	rawCert, _ := auth.PEMToCertificate(cert) // TODO optimize this...
	return auth.GetCertificateHash(rawCert), ioutil.WriteFile(certPath, cert, 0600)
}

// Start fetches the platform's private rsa key and root certificate, and create a PlatformID accordingly.
//
// The specified path should not end by a separator.
//
// The files are fetched using their default name.
func Start(path string) (*PlatformID, error) {
	keyPath := filepath.Join(path, viper.GetString("pkey_filename"))
	certPath := filepath.Join(path, viper.GetString("ca_filename"))

	// Recover the private rsa key from file.
	keyBytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	key, err := auth.PEMToPrivateKey(keyBytes)
	if err != nil {
		return nil, err
	}

	// Recover the root certificate from file.
	certBytes, err := ioutil.ReadFile(certPath)
	if err != nil {
		return nil, err
	}

	cert, err := auth.PEMToCertificate(certBytes)
	if err != nil {
		return nil, err
	}

	res := &PlatformID{
		Pkey:   key,
		RootCA: cert,
	}

	return res, nil
}
