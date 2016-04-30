// Package cmd handles flags and commands management.
package cmd

import (
	"dfss"
	dapi "dfss/dfssd/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RootCmd is the main command for the dfssp application
var RootCmd = &cobra.Command{
	Use:   "dfssp",
	Short: "Platform for the DFSS protocol",
	Long: "Platform v " + dfss.Version + ` for the
Distributed Fair Signing System project

Users and Contracts manager`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		dapi.Configure(viper.GetString("demo") != "", viper.GetString("demo"), "platform")
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		dapi.DClose()
	},
}

func init() {
	// Add flags to dfssp
	RootCmd.PersistentFlags().BoolP("verbose", "v", false, "print verbose messages")
	RootCmd.PersistentFlags().StringP("demo", "d", "", "demonstrator address and port, let empty for no debug")
	RootCmd.PersistentFlags().String("path", ".", "path to get the platform's private key and root certificate")

	initCmd.Flags().String("cn", "dfssp", "common name for the root certificate")
	initCmd.Flags().IntP("validity", "r", 365, "validity duration for the root certificate (days)")
	initCmd.Flags().String("country", "FR", "country for the root certificate")
	initCmd.Flags().String("org", "DFSS", "organization for the root certificate")
	initCmd.Flags().String("unit", "INSA Rennes", "organizational unit for the root certificate")
	initCmd.Flags().IntP("key", "k", 2048, "encoding size for the private key of the platform")

	ttpCmd.Flags().String("cn", "ttp", "common name for the ttp certificate")
	ttpCmd.Flags().IntP("validity", "c", 365, "validity duration for the ttp certificate (days)")
	ttpCmd.Flags().String("country", "FR", "country for the ttp certificate")
	ttpCmd.Flags().String("org", "DFSS", "organization for the ttp certificate")
	ttpCmd.Flags().String("unit", "INSA Rennes", "organizational unit for the ttp certificate")
	ttpCmd.Flags().IntP("key", "k", 2048, "encoding size for the private key of the ttp")

	startCmd.Flags().IntP("validity", "c", 365, "validity duration for the child certificates (days)")
	startCmd.Flags().StringP("address", "a", "0.0.0.0", "address to bind for listening")
	startCmd.Flags().StringP("port", "p", "9000", "port to bind for listening")
	startCmd.Flags().String("db", "mongodb://localhost/dfss", "server url in standard MongoDB format for accessing database")

	// Bind viper to flags
	_ = viper.BindPFlag("verbose", RootCmd.PersistentFlags().Lookup("verbose"))
	_ = viper.BindPFlag("demo", RootCmd.PersistentFlags().Lookup("demo"))
	_ = viper.BindPFlag("path", RootCmd.PersistentFlags().Lookup("path"))

	viper.SetDefault("pkey_filename", "dfssp_pkey.pem")
	viper.SetDefault("ca_filename", "dfssp_rootCA.pem")

	// Register subcommands here
	RootCmd.AddCommand(dfss.VersionCmd, ttpCmd, initCmd, startCmd)
}
