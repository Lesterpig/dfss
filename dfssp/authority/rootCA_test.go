package authority

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"

	"dfss/auth"
	"dfss/dfssc/common"
)

var (
	pkey           *rsa.PrivateKey
	PkeyFileName   = "dfssp_pkey.pem"
	RootCAFileName = "dfssp_rootCA.pem"
)

func init() {
	viper.Set("ca_filename", RootCAFileName)
	viper.Set("pkey_filename", PkeyFileName)
}

func TestMain(m *testing.M) {
	pkey, _ = auth.GeneratePrivateKey(512)
	os.Exit(m.Run())
}

func TestInitialize(t *testing.T) {
	path, _ := ioutil.TempDir("", "")
	keyPath := filepath.Join(path, PkeyFileName)
	certPath := filepath.Join(path, RootCAFileName)

	v := common.MockViper("key_size", 1024, "validity", 365, "country", "country", "organization", "organization", "unit", "unit", "cn", "cn", "path", path)
	err := Initialize(v, nil, nil)

	if err != nil {
		t.Fatal(err)
	}

	if _, err = os.Stat(keyPath); os.IsNotExist(err) {
		t.Fatal("Private key file couldn't be found")
	}

	if _, err = os.Stat(certPath); os.IsNotExist(err) {
		t.Fatal("Root certificate file couldn't be found")
	}

	_ = os.RemoveAll(path)
}

func Example() {
	path, _ := ioutil.TempDir("", "")
	keyPath := filepath.Join(path, PkeyFileName)
	certPath := filepath.Join(path, RootCAFileName)

	// Generate root certificate and key
	v := common.MockViper("key_size", 1024, "validity", 365, "country", "UK", "organization", "DFSS", "unit", "unit", "cn", "ROOT", "path", path)
	err := Initialize(v, nil, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	CheckFile(keyPath, "Private key")
	CheckFile(certPath, "Certificate")

	// Fetch files into memory
	pid, err := Start(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Generate child certificate and key
	childPath := filepath.Join(path, "child")
	v = common.MockViper("key_size", 1024, "validity", 10, "country", "FR", "organization", "DFSS", "unit", "unit", "cn", "CHILD", "path", childPath)
	err = Initialize(v, pid.RootCA, pid.Pkey)
	if err != nil {
		fmt.Println(err)
		return
	}

	CheckFile(filepath.Join(childPath, "key.pem"), "Child private key")
	CheckFile(filepath.Join(childPath, "cert.pem"), "Child certificate")

	_ = os.RemoveAll(path)
	// Output:
	// Private key file has been found
	// Certificate file has been found
	// Child private key file has been found
	// Child certificate file has been found
}

func CheckFile(path, name string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println(name + " file couldn't be found")
	} else {
		fmt.Println(name + " file has been found")
	}
}

func TestStart(t *testing.T) {
	path, _ := ioutil.TempDir("", "")
	v := common.MockViper("key_size", 1024, "validity", 365, "country", "country", "organization", "organization", "unit", "unit", "cn", "cn", "path", path)
	_ = Initialize(v, nil, nil)

	pid, err := Start(path)
	if err != nil {
		t.Fatal(err)
	}
	if pid == nil || pid.Pkey == nil || pid.RootCA == nil {
		t.Fatal("Data was not recovered from saved files")
	}

	_ = os.RemoveAll(path)
}
