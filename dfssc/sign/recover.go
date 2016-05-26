package sign

import (
	"errors"
	"io/ioutil"
	"time"

	"dfss/dfssc/common"
	"dfss/dfssc/security"
	tAPI "dfss/dfsst/api"
	"dfss/net"

	"golang.org/x/net/context"
)

// Recover : performs a recover attempt using the provided recover file, and the user's passphrase to use secured connection
func Recover(filename, passphrase string) error {
	json, err := readRecoveryFile(filename)
	if err != nil {
		return err
	}

	auth := security.NewAuthContainer(passphrase)
	_, _, _, err = auth.LoadFiles()
	if err != nil {
		return err
	}

	ttp, err := connectToTTP(json, auth)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	response, err := ttp.Recover(ctx, &tAPI.RecoverRequest{SignatureUUID: json.SignatureUUID})
	if err != nil {
		return err
	}

	return treatTTPResponse(response, auth, json.SignatureUUID)
}

// readRecoveryFile : reads the recovery file from disk
func readRecoveryFile(filename string) (*common.RecoverDataJSON, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	json, err := common.UnmarshalRecoverDataFile(data)
	if err != nil {
		return nil, err
	}

	return json, nil
}

// connectToTTP : connects to the ttp, using the provided recover data and authentication information
func connectToTTP(json *common.RecoverDataJSON, auth *security.AuthContainer) (tAPI.TTPClient, error) {
	// TODO check that the connection spots missing TTP and returns an error quickly enough
	conn, err := net.Connect(json.TTPAddrport, auth.Cert, auth.Key, auth.CA, json.TTPHash)
	if err != nil {
		return nil, err
	}
	return tAPI.NewTTPClient(conn), nil
}

// treatTTPResponse : handles the writing of the recoverd contract on the disk if it was recovered.
func treatTTPResponse(response *tAPI.TTPResponse, auth *security.AuthContainer, signatureUUID string) error {
	if len(response.Contract) == 0 {
		return errors.New("Contract was not successfully signed. Impossible to recover.")
	}

	return ioutil.WriteFile(auth.Cert.Subject.CommonName+"-"+signatureUUID+".proof", response.Contract, 0600)
}
