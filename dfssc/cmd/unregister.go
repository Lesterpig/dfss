package cmd

import (
	"fmt"
	"os"

	"dfss/dfssc/security"
	"dfss/dfssc/user"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var unregisterCmd = &cobra.Command{
	Use:   "unregister",
	Short: "delete current client information on platform",
	Run: func(cmd *cobra.Command, args []string) {
		// Read info from provided certificate
		cert, err := security.GetCertificate(viper.GetString("file_cert"))
		if err != nil {
			fmt.Fprintln(os.Stderr, "An error occurred:", err.Error())
			os.Exit(2)
		}

		// Confirmation
		var passphrase, ready string
		_ = readPassword(&passphrase, false)
		readStringParam("Do you REALLY want to delete "+cert.Subject.CommonName+"? Type 'yes' to confirm", "", &ready)
		if ready != "yes" {
			fmt.Fprintln(os.Stderr, "Unregistering aborted!")
			os.Exit(1)
		}

		err = user.Unregister(passphrase)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Cannot unregister:", err.Error())
			os.Exit(2)
		}
	},
}
