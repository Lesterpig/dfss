package security

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"dfss/auth"
	"dfss/dfssc/common"
)

// GetCertificate return the Certificate stored on the disk
func GetCertificate(filename string) (*x509.Certificate, error) {
	data, err := common.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	cert, err := auth.PEMToCertificate(data)
	if err != nil {
		return nil, err
	}

	return cert, nil
}

// GetPrivateKey return the private key stored on the disk
func GetPrivateKey(filename, passphrase string) (*rsa.PrivateKey, error) {

	data, err := common.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	key, err := auth.EncryptedPEMToPrivateKey(data, passphrase)
	if err != nil {
		return nil, err
	}

	return key, nil
}

// AES-256 requires a 32 bytes key, this function extend the key to this length
func extendKey(key string) string {
	key = strings.Repeat(key, 32/len(key)+1)
	return key[:32]
}

// EncryptStringAES enciphers the data using AES-256 algorithm
func EncryptStringAES(key string, data []byte) ([]byte, error) {
	key = extendKey(key)
	block, err := aes.NewCipher(bytes.NewBufferString(key).Bytes())
	if err != nil {
		return nil, err
	}
	b := base64.StdEncoding.EncodeToString(data)
	ciphertext := make([]byte, aes.BlockSize+len(b))

	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	return ciphertext, nil
}

// DecryptAES deciphers the data using AES-256 algorithm
func DecryptAES(key string, data []byte) ([]byte, error) {
	key = extendKey(key)
	block, err := aes.NewCipher(bytes.NewBufferString(key).Bytes())
	if err != nil {
		return nil, err
	}
	if len(data) < aes.BlockSize {
		return nil, fmt.Errorf("Ciphertext is not correctly encrypted")
	}
	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(data, data)
	decoded, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, err
	}
	return decoded, nil
}
