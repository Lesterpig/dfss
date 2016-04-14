package main

import (
	"dfss/dfssc/user"
	"fmt"
	"os"
)

// export the certificate and private key of the user
func exportConf(args []string) {
	confFile := args[0]
	fmt.Println("Export user configuration")
	var keyPassphrase, confPassphrase string

	config, err := user.NewConfig(fkey, fcert)
	if err != nil {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("Couldn't open the files: %s", err))
		os.Exit(1)
		return
	}

	err = readPassphrases(&keyPassphrase, &confPassphrase, true)
	if err != nil {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("An error occurred: %s", err))
		os.Exit(1)
		return
	}

	err = config.SaveConfigToFile(confFile, confPassphrase, keyPassphrase)
	if err != nil {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("Couldn't save the configuration on the disk: %s", err))
		os.Exit(1)
		return
	}

}

// Read two passphrases for the configuration
func readPassphrases(keyPassphrase, confPassphrase *string, second bool) error {
	fmt.Println("Enter the passphrase of the configuration")
	err := readPassword(confPassphrase, second)
	if err != nil {
		return err
	}

	fmt.Println("Enter the passphrase of your current key (if any)")
	err = readPassword(keyPassphrase, false)
	if err != nil {
		return err
	}

	return nil
}

// import the configuration
func importConf(args []string) {
	confFile := args[0]
	var keyPassphrase, confPassphrase string
	err := readPassphrases(&keyPassphrase, &confPassphrase, false)
	if err != nil {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("An error occurred: %s", err))
		os.Exit(1)
		return
	}

	config, err := user.DecodeConfiguration(confFile, keyPassphrase, confPassphrase)
	if err != nil {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("Couldnd't decrypt the configuration: %s", err))
		os.Exit(1)
		return
	}

	err = config.SaveUserInformations()
	if err != nil {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("Couldn't save the certificate and private key: %s", err))
		os.Exit(1)
		return
	}
}
