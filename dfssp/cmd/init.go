package cmd

import (
	"fmt"
	"os"

	"dfss/dfssc/common"
	"dfss/dfssp/authority"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "create and save the platform's private key and root certificate",
	Run: func(cmd *cobra.Command, args []string) {

		_ = viper.BindPFlag("cn", cmd.Flags().Lookup("cn"))
		_ = viper.BindPFlag("validity", cmd.Flags().Lookup("validity"))
		_ = viper.BindPFlag("country", cmd.Flags().Lookup("country"))
		_ = viper.BindPFlag("organization", cmd.Flags().Lookup("org"))
		_ = viper.BindPFlag("unit", cmd.Flags().Lookup("unit"))
		_ = viper.BindPFlag("key_size", cmd.Flags().Lookup("key"))

		_, err := authority.Initialize(common.SubViper("key_size", "validity", "country", "organization", "unit", "cn", "path"), nil, nil)
		if err != nil {
			fmt.Fprintln(os.Stderr, "An error occured during the initialization operation:", err)
			os.Exit(1)
		}
	},
}
