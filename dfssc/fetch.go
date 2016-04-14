package main

import (
	"fmt"
	"os"
	"path/filepath"

	"dfss/dfssc/sign"
)

func fetchContract(_ []string) {
	fmt.Println("Fetching a saved contract")

	var passphrase, uuid, directory string
	_ = readPassword(&passphrase, false)
	readStringParam("Contract UUID", "", &uuid)
	readStringParam("Save directory", ".", &directory)

	path := filepath.Join(directory, uuid+".json")
	err := sign.FetchContract(fca, fcert, fkey, addrPort, passphrase, uuid, path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
