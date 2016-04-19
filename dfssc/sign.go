package main

import (
	"fmt"
	"os"

	"dfss/dfssc/sign"
)

func signContract(args []string) {
	filename := args[0]
	fmt.Println("You are going to sign the following contract:")
	showContract(args)

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

	manager.OnSignerStatusUpdate = signFeedbackFn
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
	manager.OnProgressUpdate = signProgressFn
	err = manager.Sign()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(5)
	}

	// Persist evidencies, if any
	err = manager.PersistSignaturesToFile()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(5)
	}

	fmt.Println("Signature complete! See .proof file for evidences.")
}

func signFeedbackFn(mail string, status sign.SignerStatus, data string) {
	if status == sign.StatusConnecting {
		fmt.Println("- Trying to connect with", mail, "/", data)
	} else if status == sign.StatusConnected {
		fmt.Println("  Successfully connected!", "[", data, "]")
	}
}

func signProgressFn(current int, max int) {}
