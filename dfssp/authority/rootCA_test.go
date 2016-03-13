package authority

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"dfss/auth"
)

var pkey *rsa.PrivateKey

func TestMain(m *testing.M) {
	pkey, _ = auth.GeneratePrivateKey(512)
	os.Exit(m.Run())
}

func TestInitialize(t *testing.T) {
	path, _ := ioutil.TempDir("", "")
	keyPath := filepath.Join(path, PkeyFileName)
	certPath := filepath.Join(path, RootCAFileName)

	err := Initialize(1024, 365, "country", "organization", "unit", "cn", path, nil, nil)

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
	err := Initialize(1024, 365, "UK", "DFSS", "unit", "ROOT", path, nil, nil)
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
	err = Initialize(1024, 10, "FR", "DFSS", "unit", "CHILD", childPath, pid.RootCA, pid.Pkey)
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
	_ = Initialize(1024, 365, "country", "organization", "unit", "cn", path, nil, nil)

	pid, err := Start(path)
	if err != nil {
		t.Fatal(err)
	}
	if pid == nil || pid.Pkey == nil || pid.RootCA == nil {
		t.Fatal("Data was not recovered from saved files")
	}

	_ = os.RemoveAll(path)
}
