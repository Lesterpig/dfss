package cmd

import (
	"fmt"
	"os"

	"dfss/dfssc/sign"
	"github.com/spf13/cobra"
)

var recoverCmd = &cobra.Command{
	Use:   "recover <f>",
	Short: "try to recover signed contract from recover data file f",
	Run:   recover,
}

func recover(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		_ = cmd.Usage()
		return
	}

	var passphrase string
	_ = readPassword(&passphrase, false)
	filename := args[0]

	err := sign.Recover(filename, passphrase)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
		return
	}

	fmt.Println("Successfully recovered signed contract.")
	fmt.Println("Check .proof file.")
}
