// Package config handles basic configuration store for the GUI only.
// DFSS configuration is stored in the HOME/.dfss/config.json file
package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

// CAFile is the filename for the root certificate
const CAFile = "ca.pem"

// CertFile is the filename for the user certificate
const CertFile = "cert.pem"

// KeyFile is the filename for the private user key
const KeyFile = "key.pem"

// ConfigFile is the filename for the DFSS configuration file
const ConfigFile = "config.json"

// Config is the structure that will be persisted in the configuration file
type Config struct {
	Email    string
	Platform string

	// Virtual-only fields
	Registered    bool `json:"-"`
	Authenticated bool `json:"-"`
}

// Load loads the configuration file into memory.
// If the file does not exist, the configuration will holds default values.
func Load() (conf Config) {
	data, err := ioutil.ReadFile(getConfigFilename())
	if err != nil {
		return
	}

	_ = json.Unmarshal(data, &conf)

	// Fill virtual-only fields
	path := GetHomeDir()
	conf.Registered = isFileValid(filepath.Join(path, KeyFile))
	conf.Authenticated = isFileValid(filepath.Join(path, CertFile))
	return
}

// Save stores the current configuration object from memory.
func Save(c Config) {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return
	}

	_ = ioutil.WriteFile(getConfigFilename(), data, 0600)
}

// GetHomeDir is a helper to get the .dfss store directory
func GetHomeDir() string {
	u, err := user.Current()
	if err != nil {
		return ""
	}

	dfssPath := filepath.Join(u.HomeDir, ".dfss")
	if err := os.MkdirAll(dfssPath, os.ModeDir|0700); err != nil {
		return ""
	}

	return dfssPath + string(filepath.Separator)
}

func getConfigFilename() string {
	return filepath.Join(GetHomeDir(), ConfigFile)
}

func isFileValid(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}
