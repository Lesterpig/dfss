package security

import (
	"dfss/dfssc/common"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var crtFixture = `-----BEGIN CERTIFICATE-----
MIICRDCCAa2gAwIBAgIJAIf+q5v9t+rTMA0GCSqGSIb3DQEBCwUAMDoxCzAJBgNV
BAYTAkZSMQwwCgYDVQQKDANPUkcxEDAOBgNVBAsMB09SR1VOSVQxCzAJBgNVBAMM
AkNBMCAXDTE1MTEyMDE1NDMwOVoYDzQ3NTMxMDE2MTU0MzA5WjA6MQswCQYDVQQG
EwJGUjEMMAoGA1UECgwDT1JHMRAwDgYDVQQLDAdPUkdVTklUMQswCQYDVQQDDAJD
QTCBnzANBgkqhkiG9w0BAQEFAAOBjQAwgYkCgYEAqNbQWl2UZiOmJcZDA5x2H2U5
m2qC0D3NNPv0jNOm6shGTKhcLH1W8DrDtv5NjRWN1XTpfy0VkrsoyPpsU6PFFzZC
GmkCoXKBD/dvDNrid2MgbURzx+0a+EmUfFh0+tVP2Dzy0zgb/FZWkM6HT0VQ8KAb
SmlRBctiujDV1RBUOm8CAwEAAaNQME4wHQYDVR0OBBYEFNSIxFzdlyGUGBnqsKpd
bS4te57xMB8GA1UdIwQYMBaAFNSIxFzdlyGUGBnqsKpdbS4te57xMAwGA1UdEwQF
MAMBAf8wDQYJKoZIhvcNAQELBQADgYEAECSnGMqnlgyBdqTC02/Uo1jiPqJjLZV1
TRJFHxs4JPAsff+rdAQ1TQVfaNnvAkAoXVzM34xGPkJserMUBc7aQ61WrByGImai
RqEe6wUqHGuH2blNt+2LSSuFWuR02+LxsJARDVSLViAS3lNgXlgGnzOaRs31iwwU
czHnSiYoCog=
-----END CERTIFICATE-----
`

var path = os.TempDir()

// Test the generation of keys
func TestGenerateKeys(t *testing.T) {
	fkey := filepath.Join(path, "genKey.pem")
	viper.Set("file_key", fkey)

	rsa, err := GenerateKeys(512, "pwd")
	assert.True(t, err == nil, "An error has been raised during generation")
	assert.True(t, rsa != nil, "RSA key should not be nil")
	assert.True(t, common.FileExists(fkey), "File is missing")
	common.DeleteQuietly(fkey)
}

// Test the generation of a certificate request
func TestCertificateRequest(t *testing.T) {
	fkey := filepath.Join(path, "genCsr.pem")
	viper.Set("file_key", fkey)

	rsa, err := GenerateKeys(512, "pwd")
	defer common.DeleteQuietly(fkey)
	assert.True(t, err == nil, "An error has been raised during generation")
	assert.True(t, rsa != nil, "RSA key should not be nil")
	assert.True(t, common.FileExists(fkey), "File is missing")

	csr, err := GenerateCertificateRequest("France", "DFSS", "DFSS_C", "dfssc@dfss.org", rsa)
	assert.True(t, err == nil, "An error has been raised during generation of certificate request")
	assert.True(t, csr != "", "Certificate request should not be nil")
}

// Test saving rsa key on the disk
func TestDumpingKey(t *testing.T) {
	fkey := filepath.Join(path, "dumpKey.pem")
	viper.Set("file_key", fkey)

	rsa, err := GenerateKeys(512, "pwd")
	defer common.DeleteQuietly(fkey)

	assert.True(t, err == nil, "An error has been raised during generation")
	assert.True(t, rsa != nil, "RSA key should not be nil")
	assert.True(t, common.FileExists(fkey), "File is missing")

	k, err := GetPrivateKey(fkey, "")
	assert.True(t, err != nil, "An error should have been raised")

	k, err = GetPrivateKey(fkey, "dummypwd")
	assert.True(t, err != nil, "An error should have been raised")

	k, err = GetPrivateKey(fkey, "pwd")
	assert.True(t, err == nil, "No error should have been raised")
	assert.Equal(t, *rsa, *k, "Keys should be equal")

}

// Test the saving of a certificate in a file
func TestDumpCrt(t *testing.T) {
	fcert := filepath.Join(path, "dumpCert.pem")

	err := SaveCertificate(crtFixture, fcert)
	defer common.DeleteQuietly(fcert)

	assert.True(t, err == nil, "An error has been raised during saving")
	assert.True(t, common.FileExists(fcert), "File is missing")

	data, err := common.ReadFile(fcert)
	assert.True(t, err == nil, "An error has been raised while reading file")
	assert.Equal(t, crtFixture, fmt.Sprintf("%s", data), "Certificates are not equal")

	crt, err := GetCertificate(fcert)
	assert.True(t, err == nil, "An error has been raised while parsing certificate")
	assert.True(t, crt != nil, "Certificate is nil")
}
