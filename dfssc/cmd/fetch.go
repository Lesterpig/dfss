package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"dfss/dfssc/sign"
	"github.com/spf13/cobra"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "get a contract hosted on the platform",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Fetching a saved contract")

		var passphrase, uuid, directory string
		_ = readPassword(&passphrase, false)
		readStringParam("Contract UUID", "", &uuid)
		readStringParam("Save directory", ".", &directory)

		path := filepath.Join(directory, uuid+".json")
		err := sign.FetchContract(passphrase, uuid, path)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}
