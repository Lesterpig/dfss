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
	Use:   "ttp [cn, country, key, org, path, unit, cert-validity]",
	Short: "create and save the TTP's private key and certificate",
	Run: func(cmd *cobra.Command, args []string) {
		path := viper.GetString("path")

		pid, err := authority.Start(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Bad root CA or key; please use the `init` command before the `ttp` one.\n", err)
		}
		ttpPath := filepath.Join(path, "ttp")
		v := common.SubViper("key_size", "root_validity", "country", "organization", "unit", "cn")
		v.Set("path", ttpPath)
		err = authority.Initialize(v, pid.RootCA, pid.Pkey)
		if err != nil {
			fmt.Fprintln(os.Stderr, "An error occured during TTP credentials generation:", err)
		}
		dapi.DLog("Private key and certificate generated for TTP")
	},
}
