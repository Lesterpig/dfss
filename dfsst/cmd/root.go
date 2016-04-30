// Package cmd handles flags and commands management.
package cmd

import (
	"fmt"

	"dfss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RootCmd is the main command for the dfsst application
var RootCmd = &cobra.Command{
	Use:   "dfsst",
	Short: "DFSS TTP v" + dfss.Version,
	Long: `DFSS TTP v ` + dfss.Version + `

Trusted third party resolver for the
Distributed Fair Signing System project

Sign your contract using a secure cryptographic protocol`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

// All of the flags will be gathered by viper, this is why
// we do not store their values
func init() {
	// Bind flags to the dfsst command
	RootCmd.PersistentFlags().BoolP("verbose", "v", false, "print verbose messages")
	RootCmd.PersistentFlags().String("ca", "ca.pem", "path to the root certificate")
	RootCmd.PersistentFlags().String("cert", "cert.pem", "path to the ttp's certificate")
	RootCmd.PersistentFlags().String("key", "key.pem", "path to the ttp's private key")
	RootCmd.PersistentFlags().StringP("demo", "d", "", "demonstrator address and port, empty will disable it")

	startCmd.Flags().StringP("address", "a", "0.0.0.0", "address to bind for listening")
	startCmd.Flags().String("db", "mongodb://localhost/dfss", "server url in standard MongoDB format to access the database")
	startCmd.Flags().IntP("port", "p", 9020, "port to bind for listening")

	// Store flag values into viper
	_ = viper.BindPFlag("verbose", RootCmd.PersistentFlags().Lookup("verbose"))
	_ = viper.BindPFlag("file_ca", RootCmd.PersistentFlags().Lookup("ca"))
	_ = viper.BindPFlag("file_cert", RootCmd.PersistentFlags().Lookup("cert"))
	_ = viper.BindPFlag("file_key", RootCmd.PersistentFlags().Lookup("key"))
	_ = viper.BindPFlag("demo", RootCmd.PersistentFlags().Lookup("demo"))

	_ = viper.BindPFlag("port", startCmd.Flags().Lookup("port"))
	_ = viper.BindPFlag("address", startCmd.Flags().Lookup("address"))
	_ = viper.BindPFlag("dbURI", startCmd.Flags().Lookup("db"))

	if err := viper.BindEnv("password", "DFSS_TTP_PASSWORD"); err != nil {
		fmt.Println("Warning: The DFSS_TTP_PASSWORD environment variable is not set, assuming the private key is decrypted")
		viper.Set("password", "")
	}

	// Register Sub Commands
	RootCmd.AddCommand(dfss.VersionCmd, startCmd)

}
