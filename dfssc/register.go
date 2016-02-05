package main

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/ssh/terminal"
)

// Get a string parameter from standard input
func readStringParam(message, def string, ptr *string) {
	fmt.Printf("%s [%s]: ", message, def)
	fmt.Scanf("%s", ptr)
	if *ptr == "" {
		*ptr = def
	}
}

// Get an integer parameter from standard input
func readIntParam(message string, def int, ptr *int) {
	fmt.Printf("%s [%d]: ", message, def)
	fmt.Scanf("%d", ptr)
	if *ptr == 0 {
		*ptr = def
	}
}

// Get the password from standard input
func readPassword(ptr *string) error {
	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		return err
	}

	fmt.Println("Enter your passphrase :")
	passphrase, err := terminal.ReadPassword(0)
	if err != nil {
		return err
	}

	fmt.Println("Confirm your passphrase :")
	confirm, err := terminal.ReadPassword(0)
	if err != nil {
		return err
	}

	if fmt.Sprintf("%s", passphrase) != fmt.Sprintf("%s", confirm) {
		return errors.New("Password do not match")
	}

	*ptr = fmt.Sprintf("%s", passphrase)

	_ = terminal.Restore(0, oldState)

	return nil
}

func recapUser(fca, fcert, fkey, addrPort, mail, country, organization, unit string, bits int) {
	// Recap informations
	fmt.Println(fmt.Sprintf("Summary of the new user : Mail : %s; Country : %s; Organization : %s; Organizational unit : %s; bits : %d", mail, country, organization, unit, bits))
	fmt.Println(fmt.Sprintf("Address of the platform is %s", addrPort))
	fmt.Println(fmt.Sprintf("File storing the CA : %s", fca))
	fmt.Println(fmt.Sprintf("File to store the certificate : %s", fcert))
	fmt.Println(fmt.Sprintf("File to store the private key : %s", fkey))
}
