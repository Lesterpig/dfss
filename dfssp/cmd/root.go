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
	RootCmd.PersistentFlags().StringP("address", "a", "0.0.0.0", "address to bind for listening")
	RootCmd.PersistentFlags().StringP("port", "p", "9000", "port to bind for listening")
	RootCmd.PersistentFlags().String("path", ".", "path to get the platform's private key and root certificate")
	RootCmd.PersistentFlags().String("country", "France", "country for the root certificate")
	RootCmd.PersistentFlags().String("org", "DFSS", "Organization for the root certificate")
	RootCmd.PersistentFlags().String("unit", "INSA Rennes", "Organizational unit for the rot certificate")
	RootCmd.PersistentFlags().String("cn", "dfssp", "Common name for the root certificate")
	RootCmd.PersistentFlags().IntP("key", "k", 512, "Encoding size for the private key")
	RootCmd.PersistentFlags().IntP("root-validity", "r", 365, "Validity duration for the root certificate (days)")
	RootCmd.PersistentFlags().IntP("cert-validity", "c", 365, "Validity duration for the child certificates (days)")
	RootCmd.PersistentFlags().String("db", "mongodb://localhost/dfss", "server url in standard MongoDB format for accessing database")

	// Bind viper to flags
	_ = viper.BindPFlag("verbose", RootCmd.PersistentFlags().Lookup("verbose"))
	_ = viper.BindPFlag("demo", RootCmd.PersistentFlags().Lookup("demo"))
	_ = viper.BindPFlag("address", RootCmd.PersistentFlags().Lookup("address"))
	_ = viper.BindPFlag("port", RootCmd.PersistentFlags().Lookup("port"))
	_ = viper.BindPFlag("path", RootCmd.PersistentFlags().Lookup("path"))
	_ = viper.BindPFlag("country", RootCmd.PersistentFlags().Lookup("country"))
	_ = viper.BindPFlag("organization", RootCmd.PersistentFlags().Lookup("org"))
	_ = viper.BindPFlag("unit", RootCmd.PersistentFlags().Lookup("unit"))
	_ = viper.BindPFlag("cn", RootCmd.PersistentFlags().Lookup("cn"))
	_ = viper.BindPFlag("key_size", RootCmd.PersistentFlags().Lookup("key"))
	_ = viper.BindPFlag("root_validity", RootCmd.PersistentFlags().Lookup("root-validity"))
	_ = viper.BindPFlag("cert_validity", RootCmd.PersistentFlags().Lookup("cert-validity"))
	_ = viper.BindPFlag("dbURI", RootCmd.PersistentFlags().Lookup("db"))

	viper.SetDefault("pkey_filename", "dfssp_pkey.pem")
	viper.SetDefault("ca_filename", "dfssp_rootCA.pem")

	// Register subcommands here
	RootCmd.AddCommand(versionCmd, ttpCmd, initCmd, startCmd)

}
