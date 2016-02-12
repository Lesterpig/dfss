package main

import (
	"fmt"

	"dfss/dfssc/sign"
)

func newContract() {
	fmt.Println("Creating a new contract")

	passphrase, filepath, comment, signers := getContractInfo()
	err := sign.NewCreateManager(fca, fcert, fkey, addrPort, passphrase, filepath, comment, signers)
	if err != nil {
		fmt.Println(err)
	}
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
