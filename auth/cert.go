package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"
)

// GetCertificateRequest creates a request to be sent to any authoritative signer, as a PEM-encoded array of bytes.
//
// It can be safely sent via the network.
func GetCertificateRequest(country, organization, unit, mail string, key *rsa.PrivateKey) ([]byte, error) {
	template := &x509.CertificateRequest{
		Subject: pkix.Name{
			Country:            []string{country},
			Organization:       []string{organization},
			OrganizationalUnit: []string{unit},
			CommonName:         mail,
		},
	}

	der, err := x509.CreateCertificateRequest(rand.Reader, template, key)

	if err != nil {
		return nil, err
	}

	return pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: der,
	}), nil
}

// PEMToCertificateRequest tries to decode a PEM-encoded array of bytes to a certificate request
func PEMToCertificateRequest(data []byte) (*x509.CertificateRequest, error) {
	block, _ := pem.Decode(data)
	return x509.ParseCertificateRequest(block.Bytes)
}

// GetCertificate builds a certificate from a certificate request and an authoritative certificate (CA), as a PEM-encoded array of bytes.
// This function assumes that the identity of the signee is valid.
//
// The serial has to be unique and positive.
//
// The generated certificate can safely be distributed to unknown actors.
func GetCertificate(days int, serial uint64, req *x509.CertificateRequest, parent *x509.Certificate, key *rsa.PrivateKey) ([]byte, error) {

	template := &x509.Certificate{
		SerialNumber: new(big.Int).SetUint64(serial),
		Subject:      req.Subject,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(0, 0, days),
		IsCA:         false,
	}

	der, err := x509.CreateCertificate(rand.Reader, template, parent, req.PublicKey, key)

	if err != nil {
		return nil, err
	}

	return pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: der,
	}), nil
}

// GetSelfSignedCertificate builds a CA certificate from a private key, as a PEM-encoded array of bytes.
//
// The serial has to be unique and positive.
//
// The generated certificate should be distributed to any other actor in the network under this CA.
func GetSelfSignedCertificate(days int, serial uint64, country, organization, unit, cn string, key *rsa.PrivateKey) ([]byte, error) {

	template := &x509.Certificate{
		SerialNumber: new(big.Int).SetUint64(serial),
		Subject: pkix.Name{
			Country:            []string{country},
			Organization:       []string{organization},
			OrganizationalUnit: []string{unit},
			CommonName:         cn,
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(0, 0, days),
		IsCA:      true,
	}

	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)

	if err != nil {
		return nil, err
	}

	return pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: der,
	}), nil

}

// PEMToCertificate tries to decode a PEM-encoded array of bytes to a certificate
func PEMToCertificate(data []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(data)
	return x509.ParseCertificate(block.Bytes)
}
