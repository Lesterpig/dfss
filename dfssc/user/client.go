// Package user handles all the user-related logic
package user

import (
	"dfss/dfssc/security"
	pb "dfss/dfssp/api"
	"dfss/net"
)

// Register a user using the provided parameters
func Register(fileCA, fileCert, fileKey, addrPort, passphrase, country, organization, unit, mail string, bits int) error {
	manager, err := NewRegisterManager(fileCA, fileCert, fileKey, addrPort, passphrase, country, organization, unit, mail, bits)
	if err != nil {
		return err
	}

	err = manager.GetCertificate()
	if err != nil {
		return err
	}

	return nil
}

// Authenticate a user using the provided parameters
func Authenticate(fileCA, fileCert, addrPort, mail, token string) error {
	manager, err := NewAuthManager(fileCA, fileCert, addrPort, mail, token)
	if err != nil {
		return err
	}

	err = manager.Authenticate()
	if err != nil {
		return err
	}

	return nil
}

func connect(fileCA, addrPort string) (pb.PlatformClient, error) {
	ca, err := security.GetCertificate(fileCA)
	if err != nil {
		return nil, err
	}

	conn, err := net.Connect(addrPort, nil, nil, ca)
	if err != nil {
		return nil, err
	}

	return pb.NewPlatformClient(conn), nil
}
