package authority

import (
	"crypto/rsa"
	"dfss/auth"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

var pkey *rsa.PrivateKey

func TestMain(m *testing.M) {
	pkey, _ = auth.GeneratePrivateKey(512)
	os.Exit(m.Run())
}

func TestGetHomeDir(t *testing.T) {
	res := GetHomeDir()

	if res == "" {
		t.Fatal("Result is empty")
	}
}

func TestGenerateRootCA(t *testing.T) {
	res, err := GenerateRootCA(365, "country", "organization", "unit", "cn", pkey)

	if res == nil || err != nil {
		t.Fatal(err)
	}

	if res[0] != '-' {
		t.Fatalf("Bad format\n%s", res)
	}
}

func TestInitialize(t *testing.T) {
	path := os.TempDir()
	keyPath := filepath.Join(path, PkeyFileName)
	certPath := filepath.Join(path, RootCAFileName)

	err := Initialize(1024, 365, "country", "organization", "unit", "cn", path)

	if err != nil {
		t.Fatal(err)
	}

	if _, err = os.Stat(keyPath); os.IsNotExist(err) {
		t.Fatal("Private key file couldn't be found")
	} else {
		_ = os.Remove(keyPath)
	}

	if _, err = os.Stat(certPath); os.IsNotExist(err) {
		t.Fatal("Root certificate file couldn't be found")
	} else {
		_ = os.Remove(certPath)
	}
}

func ExampleInitialize() {
	path := os.TempDir()
	keyPath := filepath.Join(path, PkeyFileName)
	certPath := filepath.Join(path, RootCAFileName)

	err := Initialize(1024, 365, "country", "organization", "unit", "cn", path)

	if err != nil {
		fmt.Println(err)
	}

	checkFile(keyPath, "Private key")
	checkFile(certPath, "Certificate")

	// Output:
	// Private key file has been found
	// Private key file has been deleted
	// Certificate file has been found
	// Certificate file has been deleted
}

func checkFile(path, name string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println(name + " file couldn't be found")
	} else {
		fmt.Println(name + " file has been found")
		err = os.Remove(path)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(name + " file has been deleted")
		}
	}
}

func TestStart(t *testing.T) {
	path := os.TempDir()
	keyPath := filepath.Join(path, PkeyFileName)
	certPath := filepath.Join(path, RootCAFileName)

	_ = Initialize(1024, 365, "country", "organization", "unit", "cn", path)

	pid, err := Start(path)

	if err != nil {
		t.Fatal(err)
	}
	if pid == nil || pid.pkey == nil || pid.rootCA == nil {
		t.Fatal("Data was not recovered from saved files")
	}

	_ = os.Remove(keyPath)
	_ = os.Remove(certPath)
}