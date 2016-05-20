// Package cmd handles flags and commands management.
package cmd

import (
	"dfss"
	dapi "dfss/dfssd/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RootCmd is the main command for the dfssc application
var RootCmd = &cobra.Command{
	Use:   "dfssc",
	Short: "DFSS command-line client",
	Long: `Command-line client v` + dfss.Version + ` for the
Distributed Fair Signing System project

A tool to sign multiparty contract using a secure cryptographic protocol`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		dapi.Configure(viper.GetString("demo") != "", viper.GetString("demo"), "client")
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		dapi.DClose()
	},
}

// All of the flags will be gathered by viper, this is why
// we do not store their values
func init() {
	// Bind flags to the dfssc command
	RootCmd.PersistentFlags().BoolP("verbose", "v", false, "print verbose messages")
	RootCmd.PersistentFlags().String("ca", "ca.pem", "path to the root certificate")
	RootCmd.PersistentFlags().String("cert", "cert.pem", "path to the user's certificate")
	RootCmd.PersistentFlags().String("key", "key.pem", "path to the user's private key")
	RootCmd.PersistentFlags().StringP("demo", "d", "", "demonstrator address and port, empty will disable it")
	RootCmd.PersistentFlags().String("host", "localhost:9000", "host of the dfss platform")
	RootCmd.PersistentFlags().IntP("port", "p", 9005, "port to use for P2P communication between clients")

	signCmd.Flags().Duration("slowdown", 0, "delay between each promises round (test only)")
	signCmd.Flags().Int("stopbefore", 0, "stop signature just before the promises round n, -1 to stop right before signature round (test only)")

	// Store flag values into viper
	_ = viper.BindPFlag("verbose", RootCmd.PersistentFlags().Lookup("verbose"))
	_ = viper.BindPFlag("file_ca", RootCmd.PersistentFlags().Lookup("ca"))
	_ = viper.BindPFlag("file_cert", RootCmd.PersistentFlags().Lookup("cert"))
	_ = viper.BindPFlag("file_key", RootCmd.PersistentFlags().Lookup("key"))
	_ = viper.BindPFlag("demo", RootCmd.PersistentFlags().Lookup("demo"))
	_ = viper.BindPFlag("local_port", RootCmd.PersistentFlags().Lookup("port"))
	_ = viper.BindPFlag("platform_addrport", RootCmd.PersistentFlags().Lookup("host"))

	// Bind subcommands to root
	RootCmd.AddCommand(dfss.VersionCmd, registerCmd, authCmd, newCmd, showCmd, fetchCmd, importCmd, exportCmd, signCmd)

}
