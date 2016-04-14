package authority

import (
	"crypto/rsa"
	"crypto/x509"
	"io/ioutil"
	"os"
	"path/filepath"

	"dfss/auth"
)

const (
	// PkeyFileName is the private key file default name
	PkeyFileName = "dfssp_pkey.pem"
	// RootCAFileName is the root certificate file default name
	RootCAFileName = "dfssp_rootCA.pem"
)

// PlatformID contains platform private key and root certificate
type PlatformID struct {
	Pkey   *rsa.PrivateKey
	RootCA *x509.Certificate
}

// Initialize creates and saves the platform's private key and root certificate to a PEM format.
// If ca and rKey are not nil, they will be used as the root certificate and root private key instead of creating a ones.
// The files are saved at the specified path.
func Initialize(bits, days int, country, organization, unit, cn, path string, ca *x509.Certificate, rKey *rsa.PrivateKey) error {
	// Generate the private key.
	key, err := auth.GeneratePrivateKey(bits)

	if err != nil {
		return err
	}

	var cert []byte
	certPath := filepath.Join(path, RootCAFileName)
	keyPath := filepath.Join(path, PkeyFileName)

	if ca == nil {
		// Generate the root certificate, using the private key.
		cert, err = auth.GetSelfSignedCertificate(days, auth.GenerateUID(), country, organization, unit, cn, key)
	} else {
		csr, _ := auth.GetCertificateRequest(country, organization, unit, cn, key)
		request, _ := auth.PEMToCertificateRequest(csr)
		cert, err = auth.GetCertificate(days, auth.GenerateUID(), request, ca, rKey)
		// Override default path values
		certPath = filepath.Join(path, "cert.pem")
		keyPath = filepath.Join(path, "key.pem")
	}

	if err != nil {
		return err
	}

	// Create missing folders, if needed
	err = os.MkdirAll(path, os.ModeDir|0700)
	if err != nil {
		return err
	}

	// Convert the private key to a PEM format, and save it.
	keyPem := auth.PrivateKeyToPEM(key)
	err = ioutil.WriteFile(keyPath, keyPem, 0600)
	if err != nil {
		return err
	}

	// Save the root certificate.
	return ioutil.WriteFile(certPath, cert, 0600)
}

// Start fetches the platform's private rsa key and root certificate, and create a PlatformID accordingly.
//
// The specified path should not end by a separator.
//
// The files are fetched using their default name.
func Start(path string) (*PlatformID, error) {
	keyPath := filepath.Join(path, PkeyFileName)
	certPath := filepath.Join(path, RootCAFileName)

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
