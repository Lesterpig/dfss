package main

import (
	"fmt"
	"os"

	"dfss/dfssc/sign"
)

func signContract(args []string) {
	filename := args[0]
	fmt.Println("You are going to sign the following contract:")
	contract := getContract(filename)
	if contract == nil {
		os.Exit(1)
	}

	var passphrase string
	_ = readPassword(&passphrase, false)

	// Preparation
	manager, err := sign.NewSignatureManager(fca, fcert, fkey, addrPort, passphrase, localPort, contract)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	err = manager.ConnectToPeers()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(3)
	}

	// Confirmation
	var ready string
	readStringParam("Do you REALLY want to sign "+contract.File.Name+"? Type 'yes' to confirm", "", &ready)
	if ready != "yes" {
		os.Exit(4)
	}

	// Ignition
	fmt.Println("Waiting for other signers to be ready...")
	signatureUUID, err := manager.SendReadySign()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(5)
	}

	// TODO Warning, integration tests are checking Stdout
	fmt.Println("Everybody is ready, starting the signature", signatureUUID)

	// Signature
}
