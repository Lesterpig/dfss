package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	osuser "os/user"
	"strconv"
	"strings"
	"time"

	"dfss/dfssc/user"
	"golang.org/x/crypto/ssh/terminal"
)

func registerUser() {
	fmt.Println("Registering a new user")
	// Initialize variables
	var country, mail, organization, unit, passphrase string
	var bits int

	u, err := osuser.Current()
	if err != nil {
		fmt.Println("An error occurred : ", err.Error())
		return
	}

	// Get all the necessary parameters
	readStringParam("Mail", "", &mail)
	readStringParam("Country", time.Now().Location().String(), &country)
	readStringParam("Organization", u.Name, &organization)
	readStringParam("Organizational unit", u.Name, &unit)
	readIntParam("Length of the key (2048 or 4096)", "2048", &bits)
	err = readPassword(&passphrase, true)
	if err != nil {
		fmt.Println("An error occurred:", err.Error())
		return
	}

	recapUser(mail, country, organization, unit)
	err = user.Register(fca, fcert, fkey, addrPort, passphrase, country, organization, unit, mail, bits)
	if err != nil {
		fmt.Println("An error occurred:", err.Error())
	}
}

// Get a string parameter from standard input
func readStringParam(message, def string, ptr *string) {
	fmt.Print(message)
	if len(def) > 0 {
		fmt.Printf(" [%s]", def)
	}
	fmt.Print(": ")

	reader := bufio.NewReader(os.Stdin)
	value, _ := reader.ReadString('\n')

	// Trim newline symbols
	value = strings.TrimRight(value, "\n")
	value = strings.TrimRight(value, "\r")

	*ptr = value
	if value == "" {
		*ptr = def
	}

}

func readIntParam(message, def string, ptr *int) {
	var str string
	readStringParam(message, def, &str)
	value, err := strconv.Atoi(str)
	if err != nil {
		*ptr = 0
	} else {
		*ptr = value
	}
}

// Get the password from standard input
func readPassword(ptr *string, needConfirm bool) error {
	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		return err
	}

	fmt.Print("Enter your passphrase: ")
	passphrase, err := terminal.ReadPassword(0)
	fmt.Println()
	if err != nil {
		return err
	}

	if needConfirm {
		fmt.Print("Confirm your passphrase: ")
		confirm, err := terminal.ReadPassword(0)
		fmt.Println()
		if err != nil {
			return err
		}

		if fmt.Sprintf("%s", passphrase) != fmt.Sprintf("%s", confirm) {
			return errors.New("Password do not match")
		}
	}

	*ptr = fmt.Sprintf("%s", passphrase)
	_ = terminal.Restore(0, oldState)

	return nil
}

func recapUser(mail, country, organization, unit string) {
	fmt.Println("Summary of the new user:")
	fmt.Println("  Common Name:", mail)
	fmt.Println("  Country:", country)
	fmt.Println("  Organization:", organization)
	fmt.Println("  Organizational unit:", unit)
}