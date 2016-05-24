// Package config handles basic configuration store for the GUI only.
// DFSS configuration is stored in the HOME/.dfss/config.json file
package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"dfss/net"
	"github.com/spf13/viper"
)

// Config is the structure that will be persisted in the configuration file
type Config struct {
	Email    string        `json:"email"`
	Platform string        `json:"platform"`
	Timeout  time.Duration `json:"timeout"`
}

// Load loads the configuration file into memory.
// If the file does not exist, the configuration will holds default values.
func Load() {
	// Load config file
	path := GetHomeDir()
	viper.AddConfigPath(GetHomeDir())
	viper.SetConfigName(viper.GetString("filename_config"))
	viper.ReadInConfig()

	// Alias for platform
	viper.RegisterAlias("platform_addrport", "platform")

	// Setup file paths
	viper.Set("home_dir", path)
	viper.Set("file_ca", filepath.Join(path, viper.GetString("filename_ca")))
	viper.Set("file_cert", filepath.Join(path, viper.GetString("filename_cert")))
	viper.Set("file_key", filepath.Join(path, viper.GetString("filename_key")))
	viper.Set("file_config", filepath.Join(path, viper.GetString("filename_config"))+".json")

	// Fill virtual-only fields
	viper.Set("registered", isFileValid(viper.GetString("file_key")))
	viper.Set("authenticated", isFileValid(viper.GetString("file_cert")))
	viper.Set("local_port", 9005)

	// Configure timeout
	if t := viper.GetDuration("timeout"); t > 0 {
		net.DefaultTimeout = t
	}

	return
}

// Save stores the current configuration object from memory.
func Save() {
	c := Config{viper.GetString("email"), viper.GetString("platform"), viper.GetDuration("timeout")}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return
	}

	_ = ioutil.WriteFile(viper.GetString("file_config"), data, 0600)
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

func isFileValid(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}
