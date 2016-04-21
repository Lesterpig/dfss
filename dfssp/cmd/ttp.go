package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"dfss/dfssc/common"
	dapi "dfss/dfssd/api"
	"dfss/dfssp/authority"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ttpCmd = &cobra.Command{
	Use:   "ttp",
	Short: "create and save the TTP's private key and certificate",
	Run: func(cmd *cobra.Command, args []string) {

		_ = viper.BindPFlag("cn", cmd.Flags().Lookup("cn"))
		_ = viper.BindPFlag("validity", cmd.Flags().Lookup("validity"))
		_ = viper.BindPFlag("country", cmd.Flags().Lookup("country"))
		_ = viper.BindPFlag("organization", cmd.Flags().Lookup("org"))
		_ = viper.BindPFlag("unit", cmd.Flags().Lookup("unit"))
		_ = viper.BindPFlag("key_size", cmd.Flags().Lookup("key"))

		path := viper.GetString("path")

		pid, err := authority.Start(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Bad root CA or key; please use the `init` command before the `ttp` one.\n", err)
			os.Exit(1)
		}
		ttpPath := filepath.Join(path, "ttp")
		v := common.SubViper("key_size", "validity", "country", "organization", "unit", "cn")
		v.Set("path", ttpPath)
		err = authority.Initialize(v, pid.RootCA, pid.Pkey)
		if err != nil {
			fmt.Fprintln(os.Stderr, "An error occured during TTP credentials generation:", err)
			os.Exit(1)
		}
		dapi.DLog("Private key and certificate generated for TTP")
	},
}
