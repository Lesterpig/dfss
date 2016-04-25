package sign

import (
	"io/ioutil"

	"dfss/dfssc/common"
	"dfss/dfssc/security"
	"dfss/dfssp/api"
	"dfss/net"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
)

// FetchContract tries to download contract metadata from specified uuid, and stores the resulting json at path
func FetchContract(passphrase, uuid, path string) error {
	auth := security.NewAuthContainer(passphrase)
	ca, cert, key, err := auth.LoadFiles()
	if err != nil {
		return err
	}

	conn, err := net.Connect(viper.GetString("platform_addrport"), cert, key, ca, nil)
	if err != nil {
		return err
	}

	request := &api.GetContractRequest{
		Uuid: uuid,
	}
	client := api.NewPlatformClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), net.DefaultTimeout)
	defer cancel()
	response, err := client.GetContract(ctx, request)
	if err != nil {
		return err
	}

	err = common.EvaluateErrorCodeResponse(response.ErrorCode)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, response.Json, 0600)
}
