package user

import (
	"dfss/dfssc/common"
	"fmt"
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func mockb(fca, fcert, addrPort string) *viper.Viper {
	return common.MockViper("file_ca", fca, "file_cert", fcert, "platform_addrport", addrPort)
}

func TestAuthenticationValidation(t *testing.T) {
	_, err := NewAuthManager("dummy", "token", mockb("fca", fcert, addrPort))
	assert.True(t, err != nil, "Email is invalid")

	f, _ := os.Create(fcert)
	_ = f.Close()
	_, err = NewAuthManager("mail@mail.mail", "token", mockb("fca", fcert, addrPort))
	assert.True(t, err != nil, "Cert file already there")

	_ = os.Remove(fcert)
	_, err = NewAuthManager("mail@mail.mail", "token", mockb("fca", fcert, addrPort))
	assert.True(t, err != nil, "CA file not there")

	_, err = NewAuthManager("mail@mail.mail", "token", mockb(fca, fcert, addrPort))
	assert.Equal(t, err, nil)
}

func ExampleAuthenticate() {
	manager, err := NewAuthManager("mail@mail.mail", "token", mockb(fca, fcert, addrPort))
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
