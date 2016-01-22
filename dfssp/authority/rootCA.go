package authority

import (
	"crypto/rsa"
	"crypto/x509"
	"dfss/auth"
	"github.com/pborman/uuid"
	"io/ioutil"
	"math/big"
	"os/user"
	"path/filepath"
)

const (
	// PkeyFileName is the private key file default name
	PkeyFileName = "dfssp_pkey.pem"
	// RootCAFileName is the root certificate file default name
	RootCAFileName = "dfssp_rootCA.pem"
)

// PlatformID contains platform private key and root certificate
type PlatformID struct {
	pkey   *rsa.PrivateKey
	rootCA *x509.Certificate
}

// GetHomeDir determines the home directory of the current user.
func GetHomeDir() string {
	usr, err := user.Current()

	if err != nil {
		panic(err)
	}

	return usr.HomeDir
}

// GenerateRootCA constructs a self-signed certificate, using a unique serial number randomly generated (see UUID)
func GenerateRootCA(days int, country, organization, unit, cn string, key *rsa.PrivateKey) ([]byte, error) {
	// Generating and converting the uuid to fit our needs: an 8 bytes integer.
	uuid := uuid.NewRandom()
	var slice []byte
	slice = uuid[:8]
	// TODO: improve this conversion method/need
	serial := new(big.Int).SetBytes(slice).Uint64()

	cert, err := auth.GetSelfSignedCertificate(days, serial, country, organization, unit, cn, key)

	if err != nil {
		return nil, err
	}

	return cert, nil
}

// Initialize creates and saves the platform's private key and root certificate to a PEM format.
//
// The files are saved at the specified path.
func Initialize(bits, days int, country, organization, unit, cn, path string) error {
	// Generating the private key.
	key, err := auth.GeneratePrivateKey(bits)

	if err != nil {
		return err
	}

	// Generating the root certificate, using the private key.
	cert, err := GenerateRootCA(days, country, organization, unit, cn, key)

	if err != nil {
		return err
	}

	// Converting the private key to a PEM format, and saving it.
	keyPem := auth.PrivateKeyToPEM(key)
	keyPath := filepath.Join(path, PkeyFileName)
	err = ioutil.WriteFile(keyPath, keyPem, 0600)

	if err != nil {
		return err
	}

	// Saving the root certificate.
	certPath := filepath.Join(path, RootCAFileName)
	err = ioutil.WriteFile(certPath, cert, 0600)

	if err != nil {
		return err
	}

	return nil
}

// Start fetches the platform's private rsa key and root certificate, and create a PlatformID accordingly.
//
// The specified path should not end by a separator.
//
// The files are fetched using their default name.
func Start(path string) (*PlatformID, error) {
	keyPath := filepath.Join(path, PkeyFileName)
	certPath := filepath.Join(path, RootCAFileName)

	// Recovering the private rsa key from file.
	keyBytes, err := ioutil.ReadFile(keyPath)

	if err != nil {
		return nil, err
	}

	key, err := auth.PEMToPrivateKey(keyBytes)

	if err != nil {
		return nil, err
	}

	// Recovering the root certificate from file.
	certBytes, err := ioutil.ReadFile(certPath)

	if err != nil {
		return nil, err
	}

	cert, err := auth.PEMToCertificate(certBytes)

	if err != nil {
		return nil, err
	}

	res := &PlatformID{
		pkey:   key,
		rootCA: cert}

	return res, nil
}
