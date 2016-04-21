package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	osuser "os/user"
	"strconv"
	"strings"

	"dfss/dfssc/user"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "register a new client",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Registering a new user")
		// Initialize variables
		var country, mail, organization, unit, passphrase string
		var bits int

		name := "Jon Doe"
		u, err := osuser.Current()
		if err == nil {
			name = u.Name
		}

		// Get all the necessary parameters
		readStringParam("Mail", "", &mail)
		readStringParam("Country", "FR", &country)
		readStringParam("Organization", name, &organization)
		readStringParam("Organizational unit", name, &unit)
		readIntParam("Length of the key (2048 or 4096)", "2048", &bits)
		err = readPassword(&passphrase, true)
		if err != nil {
			fmt.Println("An error occurred:", err.Error())
			os.Exit(1)
			return
		}

		recapUser(mail, country, organization, unit)
		err = user.Register(passphrase, country, organization, unit, mail, bits)
		if err != nil {
			fmt.Fprintln(os.Stderr, "An error occurred:", err.Error())
			os.Exit(2)
		}
	},
}

// We need to use ONLY ONE reader: buffio buffers some data (= consumes from stdin)
var reader *bufio.Reader

// Get a string parameter from standard input
func readStringParam(message, def string, ptr *string) {
	fmt.Print(message)
	if len(def) > 0 {
		fmt.Printf(" [%s]", def)
	}
	fmt.Print(": ")

	if reader == nil {
		reader = bufio.NewReader(os.Stdin)
	}
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

	if !terminal.IsTerminal(0) {
		fmt.Println("+------------------------- WARNING --------------------------+")
		fmt.Println("| This is not a UNIX terminal, your password will be visible |")
		fmt.Println("+------------------------- WARNING --------------------------+")
		readStringParam("Enter your passphrase", "", ptr)
		return nil
	}

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
