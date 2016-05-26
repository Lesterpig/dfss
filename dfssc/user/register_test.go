package user

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"dfss/auth"
	"dfss/dfssc/common"
	"dfss/dfssc/security"
	"dfss/mockp/server"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

const caFixture = `-----BEGIN CERTIFICATE-----
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

const serverKeyFixture = `-----BEGIN RSA PRIVATE KEY-----
MIIBOwIBAAJBAMGAgCtkRLePYFRTUN0V/0v/6phm0guHGS6f0TkSEas4CGZTKFJV
TBksMGIBtfyYw3XQx2bO8myeypDN5nV05DcCAwEAAQJAHSdRKDh5KfbOGqZa3pR7
3GV4YPHM37PBFYc6rJCOXO9W8L4Q1kvEhjKXp7ke18Cge7bVmlKspvxvC62gxSQm
QQIhAPMYwpp29ZREdk8yU65Sp6w+EbZS9TjZkC+pk3syYjaxAiEAy8XWnnDMsUxb
6vp1SaaIfxI441AYzh3+8c56CAvt02cCIQDQ2jfvHz7zyDHg7rsILMkTaSwseW9n
DTwcRtOHZ40LsQIgDWEVAVwopG9+DYSaVNahWa6Jm6szpbzkc136NzMJT3sCIQDv
T2KSQQIYEvPYZmE+1b9f3rs/w7setrGtqVFkm/fTWQ==
-----END RSA PRIVATE KEY-----
`

// Use temporary files
var path = os.TempDir()
var fca = filepath.Join(path, "root_ca.pem")
var fcert = filepath.Join(path, "cert.pem")
var fkey = filepath.Join(path, "key.pem")
var addrPort = "localhost:9000"

func mock(fca, fcert, fkey string) *viper.Viper {
	return common.MockViper("file_key", fkey, "file_cert", fcert, "file_ca", fca)
}

// Main test function
func TestMain(m *testing.M) {
	// Generate a certificate and save it on the disk
	_ = security.SaveCertificate(fmt.Sprintf("%s", []byte(caFixture)), fca)

	ca, _ := auth.PEMToCertificate([]byte(caFixture))
	skey, _ := auth.PEMToPrivateKey([]byte(serverKeyFixture))

	// Init config
	viper.Set("file_ca", fca)
	viper.Set("file_cert", fcert)
	viper.Set("file_key", fkey)
	viper.Set("platform_addrport", addrPort)

	// Start the platform mock
	go server.Run(ca, skey, addrPort)
	time.Sleep(2 * time.Second)

	// Run tests
	code := m.Run()

	os.Exit(code)
}

// Test the validation of the fields
func TestRegisterValidation(t *testing.T) {
	_, err := NewRegisterManager("password", "FR", "organization", "unit", "dummy", 2048, mock(fca, fcert, fkey))
	assert.True(t, err != nil, "Email is invalid")

	_, err = NewRegisterManager("password", "FR", "organization", "unit", "mpcs@dfss.io", 2048, mock(fca, fkey, fkey))
	assert.True(t, err != nil, "Cert file is the same as key file")

	_, err = NewRegisterManager("password", "FR", "organization", "unit", "mpcs@dfss.io", 2048, mock("inexistant.pem", fcert, fkey))
	assert.True(t, err != nil, "CA file is invalid")

	f, _ := os.Create(fcert)
	_ = f.Close()
	_, err = NewRegisterManager("password", "FR", "organization", "unit", "mpcs@dfss.io", 2048, mock(fca, fcert, fkey))
	assert.True(t, err != nil, "Cert file already exist")

	k, _ := os.Create(fkey)
	_ = k.Close()
	_, err = NewRegisterManager("password", "FR", "organization", "unit", "mpcs@dfss.io", 2048, mock(fca, fcert, fkey))
	assert.True(t, err != nil, "Key file already exist")

	_ = os.Remove(fcert)
	_ = os.Remove(fkey)
}

// Test the error codes received from the mock
// Only the SUCCESS code should not raise an error
func TestGetCertificate(t *testing.T) {
	manager, err := NewRegisterManager("password", "FR", "organization", "unit", "dfss@success.io", 2048, mock(fca, fcert, fkey))
	assert.True(t, err == nil, "An error occurred while processing")
	err = manager.GetCertificate()
	assert.True(t, err == nil, "An error occurred while getting the certificate")

	go testRegisterInvalidResponse(t, "dfss@invarg.io")
	go testRegisterInvalidResponse(t, "dfss@badauth.io")
	go testRegisterInvalidResponse(t, "dfss@warning.io")
	go testRegisterInvalidResponse(t, "dfss@interr.io")
	go testRegisterInvalidResponse(t, "dfss@inexistant.io")
}

// Test an invalid error code and check we get an error
func testRegisterInvalidResponse(t *testing.T, mail string) {
	manager, err := NewRegisterManager("password", "FR", "organization", "unit", mail, 2048, mock(fca, fcert+mail, fkey+mail))

	assert.True(t, err == nil, "An error occurred while processing")
	err = manager.GetCertificate()
	assert.True(t, err != nil, "An error should have occurred while getting the certificate")

	_ = os.Remove(fcert + mail)
	_ = os.Remove(fkey + mail)
}
