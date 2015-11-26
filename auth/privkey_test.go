package auth

import (
	"crypto/x509"
	"fmt"
	"reflect"
	"testing"
)

func TestGeneratePrivateKey(t *testing.T) {
	_, err := GeneratePrivateKey(1024)

	if err != nil {
		t.Fatal(err)
	}

	if !testing.Short() {
		_, err = GeneratePrivateKey(4096)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestPrivateKeyToPEM(t *testing.T) {
	key, _ := GeneratePrivateKey(512)
	res := PrivateKeyToPEM(key)

	if res[0] != '-' {
		t.Fatalf("Bad format\n%s", res)
	}

	if IsPEMEncrypted(res) {
		t.Fatal("Result is encrypted")
	}
}

func TestPrivateKeyToEncryptedPEM(t *testing.T) {
	key, _ := GeneratePrivateKey(512)
	res, err := PrivateKeyToEncryptedPEM(key, "password")

	if res[0] != '-' {
		t.Fatalf("Bad format\n%s", res)
	}

	if !IsPEMEncrypted(res) || err != nil {
		t.Fatal("Result is not encrypted: ", err)
	}
}

func TestPEMToPrivateKey(t *testing.T) {
	key, _ := GeneratePrivateKey(512)
	key2, err := PEMToPrivateKey(PrivateKeyToPEM(key))
	if !reflect.DeepEqual(key, key2) || err != nil {
		t.Fatal(err)
	}
}

func TestEncryptedPEMToPrivateKey(t *testing.T) {
	key, _ := GeneratePrivateKey(512)
	res, _ := PrivateKeyToEncryptedPEM(key, "password")

	goodKey, err := EncryptedPEMToPrivateKey(res, "password")

	if !reflect.DeepEqual(key, goodKey) || err != nil {
		t.Fatal(err)
	}

	badKey, err := EncryptedPEMToPrivateKey(res, "badpass")

	if badKey != nil || err != x509.IncorrectPasswordError {
		t.Fatal(err)
	}
}

func ExampleEncryptedPEMToPrivateKey() {

	// Generate a new private key for example
	key, err := GeneratePrivateKey(512)

	if err != nil {
		panic(err)
	}

	// Get the encrypted PEM data

	p, err := PrivateKeyToEncryptedPEM(key, "myPassword")

	if err != nil {
		panic(err)
	}

	// Reverse the process

	newKey, err := EncryptedPEMToPrivateKey(p, "badPassword")

	if err == x509.IncorrectPasswordError {
		fmt.Println("Bad password")
	}

	newKey, err = EncryptedPEMToPrivateKey(p, "myPassword")

	if newKey != nil && err == nil {
		fmt.Println("OK")
	}

	// Output:
	// Bad password
	// OK
}
