package sign

import (
	"io/ioutil"
	"time"

	"dfss/dfssc/common"
	"dfss/dfssc/security"
	"dfss/dfssp/api"
	"dfss/net"
	"golang.org/x/net/context"
)

// FetchContract tries to download contract metadata from specified uuid, and stores the resulting json at path
func FetchContract(fileCA, fileCert, fileKey, addrPort, passphrase, uuid, path string) error {
	auth := security.NewAuthContainer(fileCA, fileCert, fileKey, addrPort, passphrase)
	ca, cert, key, err := auth.LoadFiles()
	if err != nil {
		return err
	}

	conn, err := net.Connect(auth.AddrPort, cert, key, ca)
	if err != nil {
		return err
	}

	request := &api.GetContractRequest{
		Uuid: uuid,
	}
	client := api.NewPlatformClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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
