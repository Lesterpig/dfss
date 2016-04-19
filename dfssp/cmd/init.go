package cmd

import (
	"fmt"
	"os"

	"dfss/dfssc/common"
	dapi "dfss/dfssd/api"
	"dfss/dfssp/authority"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [cn, country, key, org, path, unit, root-validity]",
	Short: "create and save the platform's private key and root certificate",
	Run: func(cmd *cobra.Command, args []string) {
		err := authority.Initialize(common.SubViper("key_size", "root_validity", "country", "organization", "unit", "cn", "path"), nil, nil)
		if err != nil {
			fmt.Fprintln(os.Stderr, "An error occured during the initialization operation:", err)
			os.Exit(1)
		}
		dapi.DLog("Private key and root certificate generated")
	},
}
