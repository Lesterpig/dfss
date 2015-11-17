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
		t.Fail()
	}

	if !testing.Short() {
		GeneratePrivateKey(4096)
	}
}

func TestPrivateKeyToPEM(t *testing.T) {
	key, _ := GeneratePrivateKey(2048)
	res := PrivateKeyToPEM(key)

	if IsPEMEncrypted(res) {
		t.Fail()
	}
}

func TestPrivateKeyToEncryptedPEM(t *testing.T) {
	key, _ := GeneratePrivateKey(2048)
	res, err := PrivateKeyToEncryptedPEM(key, "password")

	if !IsPEMEncrypted(res) || err != nil {
		t.Fail()
	}
}

func TestPEMToPrivateKey(t *testing.T) {
	key, _ := GeneratePrivateKey(2048)
	key2, err := PEMToPrivateKey(PrivateKeyToPEM(key))
	if !reflect.DeepEqual(key, key2) || err != nil {
		t.Fail()
	}
}

func TestEncryptedPEMToPrivateKey(t *testing.T) {
	key, _ := GeneratePrivateKey(2048)
	res, _ := PrivateKeyToEncryptedPEM(key, "password")

	goodKey, err := EncryptedPEMToPrivateKey(res, "password")

	if !reflect.DeepEqual(key, goodKey) || err != nil {
		t.Fail()
	}

	badKey, err := EncryptedPEMToPrivateKey(res, "badpass")

	if badKey != nil || err != x509.IncorrectPasswordError {
		t.Fail()
	}
}

func ExampleEncryptedPEMToPrivateKey() {

	// Generate a new private key for example
	key, err := GeneratePrivateKey(2048)

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
