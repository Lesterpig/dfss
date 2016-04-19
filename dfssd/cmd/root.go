package cmd

import (
	"dfss"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RootCmd is the main command for the dfssd application
var RootCmd = &cobra.Command{
	Use:   "dfssd",
	Short: "Demonstrator for the DFSS project",
	Long: "Demonstrator v" + dfss.Version + ` for the
Distributed Fair Signing System project

Debug tool to trace remote transmissions`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	// Add flag to the command
	RootCmd.PersistentFlags().IntP("port", "p", 3000, "port to use for listening transmissions")

	// Bind the flag to viper
	_ = viper.BindPFlag("port", RootCmd.PersistentFlags().Lookup("port"))

	// Register subcommands
	RootCmd.AddCommand(versionCmd, noguiCmd)
}
