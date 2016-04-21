package user

import (
	"dfss/dfssc/common"
	"dfss/dfssc/security"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const certFixture = `-----BEGIN CERTIFICATE-----
MIIB5TCCAY+gAwIBAgIJAKId2y6Lo9T8MA0GCSqGSIb3DQEBCwUAME0xCzAJBgNV
BAYTAkZSMQ0wCwYDVQQKDARERlNTMRswGQYDVQQLDBJERlNTIFBsYXRmb3JtIHYw
LjExEjAQBgNVBAMMCWxvY2FsaG9zdDAgFw0xNjAxMjYxNTM2NTNaGA80NDgwMDMw
ODE1MzY1M1owTTELMAkGA1UEBhMCRlIxDTALBgNVBAoMBERGU1MxGzAZBgNVBAsM
EkRGU1MgUGxhdGZvcm0gdjAuMTESMBAGA1UEAwwJbG9jYWxob3N0MFwwDQYJKoZI
hvcNAQEBBQADSwAwSAJBAMGAgCtkRLePYFRTUN0V/0v/6phm0guHGS6f0TkSEas4
CGZTKFJVTBksMGIBtfyYw3XQx2bO8myeypDN5nV05DcCAwEAAaNQME4wHQYDVR0O
BBYEFO09nxx5/qeLK5Wig1+3kg66gn/mMB8GA1UdIwQYMBaAFO09nxx5/qeLK5Wi
g1+3kg66gn/mMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQELBQADQQCqNSH+rt/Z
ru2rkabLiHOGjI+AenSOvqWZ2dWAlLksYcyuQHKwjGWgpmqkiQCnkIDwIxZvu69Y
OBz0ASFn7eym
-----END CERTIFICATE-----
`

const keyFixture = `-----BEGIN RSA PRIVATE KEY-----
MIIBOwIBAAJBAMGAgCtkRLePYFRTUN0V/0v/6phm0guHGS6f0TkSEas4CGZTKFJV
TBksMGIBtfyYw3XQx2bO8myeypDN5nV05DcCAwEAAQJAHSdRKDh5KfbOGqZa3pR7
3GV4YPHM37PBFYc6rJCOXO9W8L4Q1kvEhjKXp7ke18Cge7bVmlKspvxvC62gxSQm
QQIhAPMYwpp29ZREdk8yU65Sp6w+EbZS9TjZkC+pk3syYjaxAiEAy8XWnnDMsUxb
6vp1SaaIfxI441AYzh3+8c56CAvt02cCIQDQ2jfvHz7zyDHg7rsILMkTaSwseW9n
DTwcRtOHZ40LsQIgDWEVAVwopG9+DYSaVNahWa6Jm6szpbzkc136NzMJT3sCIQDv
T2KSQQIYEvPYZmE+1b9f3rs/w7setrGtqVFkm/fTWQ==
-----END RSA PRIVATE KEY-----
`

var cPath = os.TempDir()

// TestInvalidFiles assert an error is raised when provided with wrong files
func TestInvalidFiles(t *testing.T) {
	keyFile := filepath.Join(cPath, "invalidKey.pem")
	certFile := filepath.Join(cPath, "invalidCert.pem")
	defer deleteFiles(keyFile, certFile, "invalidConf.pem")

	_, err := NewConfig(common.MockViper("file_key", "inexistantKey", "file_cert", "inexistantCert"))
	assert.True(t, err != nil, "No key file nor cert file, expected error")

	_ = common.SaveStringToDisk(fmt.Sprintf("%s", []byte(keyFixture)), keyFile)
	_, err = NewConfig(common.MockViper("file_key", keyFile, "file_cert", "inexistantCert"))
	assert.True(t, err != nil, "No cert file, expected error")

	_ = security.SaveCertificate(fmt.Sprintf("%s", []byte(certFixture)), certFile)
	_, err = NewConfig(common.MockViper("file_key", keyFile, "file_cert", certFile))
	assert.True(t, err == nil, "Expected no error, files are present and valid")
}

// TestErrorDumpingConfig checks the error that may be raised while dumping the configuration to the disk
func TestErrorDumpingConfig(t *testing.T) {
	keyFile := filepath.Join(cPath, "privKey.pem")
	certFile := filepath.Join(cPath, "cert.pem")
	mockViper := common.MockViper("file_key", keyFile, "file_cert", certFile)
	defer deleteFiles(keyFile, certFile, "invalidConf.pem")

	err := common.SaveStringToDisk(fmt.Sprintf("%s", []byte(keyFixture)), keyFile)
	assert.True(t, err == nil, "Expected no error, file is present")

	err = security.SaveCertificate(fmt.Sprintf("%s", []byte(certFixture)), certFile)
	assert.True(t, err == nil, "Expected no error, cert is valid")

	config, err := NewConfig(mockViper)
	assert.True(t, err == nil, "Expected no error, files are present and valid")

	err = config.SaveConfigToFile("file", "abc", "")
	assert.True(t, err != nil, "Expected an error, passphrase is too short (< 4 char)")

	err = config.SaveConfigToFile(keyFile, "passphrase", "")
	assert.True(t, err != nil, "Expected an error, file is already there")

	common.DeleteQuietly(keyFile)
	_ = common.SaveStringToDisk("Invalid key", keyFile)
	config, _ = NewConfig(mockViper)
	err = config.SaveConfigToFile("file", "passphrase", "passphrase")
	assert.True(t, err != nil, "Expected an error, private key is invalid")

	common.DeleteQuietly(certFile)
	common.DeleteQuietly(keyFile)
	_ = common.SaveStringToDisk(fmt.Sprintf("%s", []byte(keyFixture)), keyFile)
	_ = common.SaveStringToDisk("Invalid certificate", certFile)
	config, _ = NewConfig(mockViper)
	err = config.SaveConfigToFile("file", "passphrase", "passphrase")
	assert.True(t, err != nil, "Expected an error, certificate is invalid")

}

// TestDumpingFile tries to save the configuration and checks there are no problems
func TestDumpingFile(t *testing.T) {
	keyFile := filepath.Join(cPath, "privKey2.pem")
	certFile := filepath.Join(cPath, "cert2.pem")
	configPath := filepath.Join(os.TempDir(), "dfss.conf")
	common.DeleteQuietly(configPath)
	defer deleteFiles(keyFile, certFile, configPath)

	err := common.SaveStringToDisk(fmt.Sprintf("%s", []byte(keyFixture)), keyFile)
	assert.True(t, err == nil, "Expected no error, file is present")

	err = security.SaveCertificate(fmt.Sprintf("%s", []byte(certFixture)), certFile)
	assert.True(t, err == nil, "Expected no error, cert is valid")

	config, err := NewConfig(common.MockViper("file_key", keyFile, "file_cert", certFile))
	assert.True(t, err == nil, "Expected no error, files are present and valid")

	err = config.SaveConfigToFile(configPath, "passphrase", "")
	assert.True(t, err == nil, "Expected no error, config is valid")

	assert.True(t, common.FileExists(configPath), "Expected a config file present")
	_, err = common.ReadFile(configPath)

	assert.True(t, err == nil, "Expected no error, config is present")
}

// TestErrorDecodeFile tries to decode a configuration file and checks the errors raised
func TestErrorDecodeFile(t *testing.T) {
	keyFile := filepath.Join(cPath, "privKey2.pem")
	certFile := filepath.Join(cPath, "cert2.pem")
	configPath := filepath.Join(os.TempDir(), "dfss2.conf")
	common.DeleteQuietly(configPath)
	defer deleteFiles(keyFile, certFile, configPath)

	_ = common.SaveStringToDisk(fmt.Sprintf("%s", []byte(keyFixture)), keyFile)
	_ = security.SaveCertificate(fmt.Sprintf("%s", []byte(certFixture)), certFile)

	_, err := DecodeConfiguration("inexistantFile", "passphrase", "")
	assert.True(t, err != nil, "File is invalid, impossible to decode configuration")

	_, err = DecodeConfiguration(keyFile, "pas", "")
	assert.True(t, err != nil, "Passphrase is invalid, should be at least 4 char long")

	config, err := NewConfig(common.MockViper("file_key", keyFile, "file_cert", certFile))
	assert.True(t, err == nil, "Expected no error, files are present and valid")

	err = config.SaveConfigToFile(configPath, "passphrase", "")
	assert.True(t, err == nil, "Expected no error, config is valid")

	config, err = DecodeConfiguration(configPath, "pass", "")
	assert.True(t, err != nil, "Expected error, wrong passphrase")

}

// TestDecodeConfig tries to decode the configuration and checks it is right
func TestDecodeConfig(t *testing.T) {
	keyFile := filepath.Join(cPath, "privKey4.pem")
	certFile := filepath.Join(cPath, "cert4.pem")
	configPath := filepath.Join(os.TempDir(), "dfss4.conf")
	common.DeleteQuietly(configPath)
	defer deleteFiles(keyFile, certFile, configPath)

	_ = common.SaveStringToDisk(fmt.Sprintf("%s", []byte(keyFixture)), keyFile)
	_ = security.SaveCertificate(fmt.Sprintf("%s", []byte(certFixture)), certFile)

	config, err := NewConfig(common.MockViper("file_key", keyFile, "file_cert", certFile))
	assert.True(t, err == nil, "Expected no error, files are present and valid")

	err = config.SaveConfigToFile(configPath, "passphrase", "")
	assert.True(t, err == nil, "Expected no error, config is valid")

	decoded, err := DecodeConfiguration(configPath, "", "passphrase")
	assert.True(t, err == nil, "Expected no error, config should have been decodedi")

	assert.Equal(t, config.KeyFile, decoded.KeyFile, "Wrong keyFile parameter")
	assert.Equal(t, config.KeyData, decoded.KeyData, "Wrong keyData parameter")
	assert.Equal(t, config.CertFile, decoded.CertFile, "Wrong certFile parameter")
	assert.Equal(t, config.CertData, decoded.CertData, "Wrong certData parameter")
}

// TestSaveFilesToDisk tries to create the certificate and private key on the disk from
// the config file
func TestSaveFilesToDisk(t *testing.T) {
	keyFile := filepath.Join(cPath, "privKey5.pem")
	certFile := filepath.Join(cPath, "cert5.pem")
	configPath := filepath.Join(os.TempDir(), "dfss5.conf")
	deleteFiles(keyFile, certFile, configPath)
	defer deleteFiles(keyFile, certFile, configPath)

	_ = common.SaveStringToDisk(fmt.Sprintf("%s", []byte(keyFixture)), keyFile)
	_ = security.SaveCertificate(fmt.Sprintf("%s", []byte(certFixture)), certFile)

	config, err := NewConfig(common.MockViper("file_key", keyFile, "file_cert", certFile))
	assert.True(t, err == nil, "Expected no error, files are present and valid")

	err = config.SaveUserInformations()
	assert.True(t, err != nil, "Expected an error, files are already present")
	common.DeleteQuietly(keyFile)

	err = config.SaveUserInformations()
	assert.True(t, err != nil, "Expected an error, certificate file is already present")

	common.DeleteQuietly(certFile)

	err = config.SaveUserInformations()
	assert.True(t, err == nil, "No error expected, files are not here")

	assert.True(t, common.FileExists(keyFile), "Expected private key file")
	assert.True(t, common.FileExists(certFile), "Expected certificate file")
}

// Helper function to delete all the files
func deleteFiles(keyFile, certFile, confFile string) {
	common.DeleteQuietly(keyFile)
	common.DeleteQuietly(certFile)
	common.DeleteQuietly(confFile)
}
