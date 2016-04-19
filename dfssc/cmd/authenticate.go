package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"dfss/dfssc/user"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "authenticate a new client",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Authenticating user")
		var mail, token string

		readStringParam("Mail", "", &mail)
		readStringParam("Token", "", &token)

		err := user.Authenticate(mail, token)
		if err != nil {
			fmt.Println("An error occurred : ", err.Error())
			os.Exit(3)
		}
	},
}
