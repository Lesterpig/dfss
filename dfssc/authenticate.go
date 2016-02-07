package main

import (
	"fmt"

	"dfss/dfssc/user"
)

func authUser() {
	fmt.Println("Authenticating user")
	var mail, token string

	readStringParam("Mail", "", &mail)
	readStringParam("Token", "", &token)

	err := user.Authenticate(fca, fcert, addrPort, mail, token)
	if err != nil {
		fmt.Println("An error occurred : ", err.Error())
	}
}
