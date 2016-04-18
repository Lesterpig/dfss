package user

import (
	"dfss/dfssc/common"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestAuthenticationValidation(t *testing.T) {
	_, err := NewAuthManager("fca", fcert, addrPort, "dummy", "token")
	assert.True(t, err != nil, "Email is invalid")

	f, _ := os.Create(fcert)
	_ = f.Close()
	_, err = NewAuthManager("fca", fcert, addrPort, "mail@mail.mail", "token")
	assert.True(t, err != nil, "Cert file already there")

	_ = os.Remove(fcert)
	_, err = NewAuthManager("fca", fcert, addrPort, "mail@mail.mail", "token")
	assert.True(t, err != nil, "CA file not there")

	_, err = NewAuthManager(fca, fcert, addrPort, "mail@mail.mail", "token")
	assert.Equal(t, err, nil)
}

func ExampleAuthenticate() {
	manager, err := NewAuthManager(fca, fcert, addrPort, "mail@mail.mail", "token")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Manager successfully created")

	err = manager.Authenticate()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Manager said authentication went fine")

	if b := common.FileExists(fcert); !b {
		fmt.Println("The cert file was not saved to disk")
	} else {
		fmt.Println("The cert file was saved to disk")
		_ = os.Remove(fcert)
	}

	// Output:
	// Manager successfully created
	// Manager said authentication went fine
	// The cert file was saved to disk
}
