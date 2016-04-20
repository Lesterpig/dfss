// Package user handles all the user-related logic
package user

import (
	"dfss/dfssc/common"
	"dfss/dfssc/security"
	pb "dfss/dfssp/api"
	"dfss/net"

	"github.com/spf13/viper"
)

// Register a user using the provided parameters
func Register(passphrase, country, organization, unit, mail string, bits int) error {
	manager, err := NewRegisterManager(passphrase, country, organization, unit, mail, bits, common.SubViper("file_key", "file_cert", "file_ca"))
	if err != nil {
		return err
	}
	return manager.GetCertificate()
}

// Authenticate a user using the provided parameters
func Authenticate(mail, token string) error {
	manager, err := NewAuthManager(mail, token, common.SubViper("file_ca", "file_cert"))
	if err != nil {
		return err
	}
	return manager.Authenticate()
}

func connect() (pb.PlatformClient, error) {
	ca, err := security.GetCertificate(viper.GetString("file_ca"))
	if err != nil {
		return nil, err
	}

	conn, err := net.Connect(viper.GetString("platform_addrport"), nil, nil, ca, nil)
	if err != nil {
		return nil, err
	}

	return pb.NewPlatformClient(conn), nil
}
