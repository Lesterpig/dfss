package user

import (
	"dfss/auth"
	"dfss/dfssc/common"
	"dfss/dfssc/security"
	"encoding/json"
	"fmt"

	"github.com/spf13/viper"
)

// Config represents the config file to be marshalled in json
type Config struct {
	KeyFile  string `json:"key"`
	KeyData  []byte `json:"keyData"`
	CertFile string `json:"cert"`
	CertData []byte `json:"certData"`
}

// NewConfig creates a new config object from the key and certificate provided
// The validity of those is checked later
func NewConfig(viper *viper.Viper) (*Config, error) {
	keyFile := viper.GetString("file_key")
	if !common.FileExists(keyFile) {
		return nil, fmt.Errorf("No such file: %s", keyFile)
	}

	certFile := viper.GetString("file_cert")
	if !common.FileExists(certFile) {
		return nil, fmt.Errorf("No such file: %s", certFile)
	}

	key, err := common.ReadFile(keyFile)

	if err != nil {
		return nil, err
	}

	cert, err := common.ReadFile(certFile)

	if err != nil {
		return nil, err
	}

	return &Config{
		KeyData:  key,
		KeyFile:  keyFile,
		CertData: cert,
		CertFile: certFile,
	}, nil

}

// SaveConfigToFile marshals checks the  validity of the certificate and private key,
// marshals the struct in JSON, encrypt the string using AES-256 with the provided passphrase,
// and finally save it to a file
func (c *Config) SaveConfigToFile(fileName, passphrase, keyPassphrase string) error {
	if common.FileExists(fileName) {
		return fmt.Errorf("Cannot overwrite file: %s", fileName)
	}

	if len(passphrase) < 4 {
		return fmt.Errorf("Passphrase should be at least 4 characters long")
	}

	err := c.checkData(keyPassphrase)
	if err != nil {
		return err
	}

	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	encodedData, err := security.EncryptStringAES(passphrase, data)
	if err != nil {
		return err
	}

	err = common.SaveToDisk(encodedData, fileName)
	return err
}

// DecodeConfiguration : decrypt and unmarshal the given configuration file
// to create a Config object. It also checks the validity of the certificate and private key
func DecodeConfiguration(fileName, keyPassphrase, confPassphrase string) (*Config, error) {
	if !common.FileExists(fileName) {
		return nil, fmt.Errorf("No such file: %s", fileName)
	}

	if len(confPassphrase) < 4 {
		return nil, fmt.Errorf("Passphrase should be at least 4 characters long")
	}

	encodedData, err := common.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	decodedData, err := security.DecryptAES(confPassphrase, encodedData)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(decodedData, &config)
	if err != nil {
		return nil, err
	}

	err = config.checkData(keyPassphrase)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// Check that the certificate is valid, and that the private key is valid too
// using the passphrase
func (c *Config) checkData(keyPassphrase string) error {
	_, err := auth.PEMToCertificate(c.CertData)
	if err != nil {
		return err
	}

	_, err = auth.EncryptedPEMToPrivateKey(c.KeyData, keyPassphrase)
	return err
}

// SaveUserInformations save the certificate and private key to the files specified in the Config struct
func (c *Config) SaveUserInformations() error {
	if common.FileExists(c.KeyFile) {
		return fmt.Errorf("Cannot overwrite file: %s", c.KeyFile)
	}

	if common.FileExists(c.CertFile) {
		return fmt.Errorf("Cannot overwrite file: %s", c.CertFile)
	}

	err := common.SaveToDisk(c.KeyData, c.KeyFile)
	if err != nil {
		common.DeleteQuietly(c.KeyFile)
		return err
	}

	err = common.SaveToDisk(c.CertData, c.CertFile)
	if err != nil {
		common.DeleteQuietly(c.KeyFile)
		common.DeleteQuietly(c.CertFile)
		return err
	}

	return nil
}
