package user

import (
	"dfss/dfssc/common"
	"dfss/dfssc/security"
	"dfss/dfssp/api"
	"dfss/net"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
)

// Unregister a user from the platform
func Unregister(passphrase string) error {
	auth := security.NewAuthContainer(passphrase)
	ca, cert, key, err := auth.LoadFiles()
	if err != nil {
		return err
	}

	conn, err := net.Connect(viper.GetString("platform_addrport"), cert, key, ca, nil)
	if err != nil {
		return err
	}

	client := api.NewPlatformClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), net.DefaultTimeout)
	defer cancel()
	response, err := client.Unregister(ctx, &api.Empty{})

	if err != nil {
		return err
	}

	return common.EvaluateErrorCodeResponse(response)
}
