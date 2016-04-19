package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"dfss/dfssc/sign"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "create a new contract",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Creating a new contract")

		passphrase, filepath, comment, signers := getContractInfo()
		err := sign.SendNewContract(passphrase, filepath, comment, signers)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

// getContractInfo asks user for contract informations
func getContractInfo() (passphrase string, path string, comment string, signers []string) {

	signers = make([]string, 1)

	var signersBuf string
	_ = readPassword(&passphrase, false)
	readStringParam("Contract path", "", &path)
	readStringParam("Comment", "(no comment)", &comment)
	readStringParam("Signer 1", "mail@example.com", &signers[0])

	i := 2
	for {
		readStringParam(fmt.Sprintf("Signer %d (return to end)", i), "", &signersBuf)
		if len(signersBuf) == 0 {
			break
		}
		signers = append(signers, signersBuf)
		i++
	}
	return
}
