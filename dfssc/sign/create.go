package sign

import (
	"crypto/sha512"
	"io/ioutil"
	"path/filepath"
	"time"

	"dfss/dfssc/common"
	"dfss/dfssc/security"
	"dfss/dfssp/api"
	"dfss/net"
	"golang.org/x/net/context"
)

// CreateManager handles the creation of a new contract.
type CreateManager struct {
	auth     *security.AuthContainer
	filepath string
	comment  string
	signers  []string
	hash     []byte
	filename string
}

// NewCreateManager tries to create a contract on the platform and returns an error or nil
func NewCreateManager(fileCA, fileCert, fileKey, addrPort, passphrase, filepath, comment string, signers []string) error {
	m := &CreateManager{
		auth:     security.NewAuthContainer(fileCA, fileCert, fileKey, addrPort, passphrase),
		filepath: filepath,
		comment:  comment,
		signers:  signers,
	}

	err := m.computeFile()
	if err != nil {
		return err
	}

	result, err := m.sendRequest()
	if err != nil {
		return err
	}

	return common.EvaluateErrorCodeResponse(result)
}

// computeFile computes hash and filename providing the contract filepath
func (m *CreateManager) computeFile() error {
	data, err := ioutil.ReadFile(m.filepath)
	if err != nil {
		return err
	}

	hash := sha512.Sum512(data)
	m.hash = hash[:]
	m.filename = filepath.Base(m.filepath)

	return nil
}

// sendRequest sends a new contract request for the platform and send it
func (m *CreateManager) sendRequest() (*api.ErrorCode, error) {

	ca, cert, key, err := m.auth.LoadFiles()
	if err != nil {
		return nil, err
	}

	conn, err := net.Connect(m.auth.AddrPort, cert, key, ca)
	if err != nil {
		return nil, err
	}

	request := &api.PostContractRequest{
		Hash:     m.hash,
		Filename: m.filename,
		Signer:   m.signers,
		Comment:  m.comment,
	}

	client := api.NewPlatformClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	response, err := client.PostContract(ctx, request)
	if err != nil {
		return nil, err
	}

	return response, nil
}
