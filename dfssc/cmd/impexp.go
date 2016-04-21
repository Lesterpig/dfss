package cmd

import (
	"fmt"
	"os"

	"dfss/dfssc/common"
	"dfss/dfssc/user"
	"github.com/spf13/cobra"
)

// export the certificate and private key of the user
var exportCmd = &cobra.Command{
	Use:   "export <c>",
	Short: "export certificate and private key of the user to file c",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			_ = cmd.Usage()
			os.Exit(1)
			return
		}

		confFile := args[0]
		fmt.Println("Export user configuration")
		var keyPassphrase, confPassphrase string

		config, err := user.NewConfig(common.SubViper("file_key", "file_cert"))
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
	},
}

// import the configuration
var importCmd = &cobra.Command{
	Use:   "import <c>",
	Short: "import private key and certificate of the user from file c",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			_ = cmd.Usage()
			os.Exit(1)
			return
		}

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
	},
}

// Read two passphrases for the configuration
func readPassphrases(keyPassphrase, confPassphrase *string, second bool) error {
	fmt.Println("Enter the passphrase of the configuration")
	err := readPassword(confPassphrase, second)
	if err != nil {
		return err
	}

	fmt.Println("Enter the passphrase of your current key (if any)")
	return readPassword(keyPassphrase, false)
}
