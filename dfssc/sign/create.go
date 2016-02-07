package sign

import (
	"crypto/sha512"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"dfss/dfssc/security"
	"dfss/dfssp/api"
	"dfss/net"
	"golang.org/x/net/context"
)

// CreateManager handles the creation of a new contract
type CreateManager struct {
	fileCA     string
	fileCert   string
	fileKey    string
	addrPort   string
	passphrase string
	filepath   string
	comment    string
	signers    []string
	hash       string
	filename   string
}

// NewCreateManager tries to create a contract on the platform and returns an error or nil
func NewCreateManager(fileCA, fileCert, fileKey, addrPort, passphrase, filepath, comment string, signers []string) error {
	m := &CreateManager{
		fileCA:     fileCA,
		fileCert:   fileCert,
		fileKey:    fileKey,
		addrPort:   addrPort,
		passphrase: passphrase,
		filepath:   filepath,
		comment:    comment,
		signers:    signers,
	}

	err := m.computeFile()
	if err != nil {
		return err
	}

	result, err := m.sendRequest()
	if err != nil {
		return err
	}

	return m.evaluateResponse(result)
}

// computeFile computes hash and filename providing the contract filepath
func (m *CreateManager) computeFile() error {
	data, err := ioutil.ReadFile(m.filepath)
	if err != nil {
		return err
	}

	hash := sha512.Sum512(data)
	m.hash = fmt.Sprintf("%x", hash)
	m.filename = filepath.Base(m.filepath)

	return nil
}

// sendRequest crafts a new contract request for the platform and send it
func (m *CreateManager) sendRequest() (*api.ErrorCode, error) {
	ca, err := security.GetCertificate(m.fileCA)
	if err != nil {
		return nil, err
	}

	cert, err := security.GetCertificate(m.fileCert)
	if err != nil {
		return nil, err
	}

	key, err := security.GetPrivateKey(m.fileKey, m.passphrase)
	if err != nil {
		return nil, err
	}

	conn, err := net.Connect(m.addrPort, cert, key, ca)
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

// evaluateResponse evaluates the errorCode received
func (m *CreateManager) evaluateResponse(code *api.ErrorCode) error {
	if code.Code == api.ErrorCode_SUCCESS {
		return nil
	} else if code.Code == api.ErrorCode_WARNING {
		return errors.New("Contract created with a warning message: " + code.Message)
	}
	return errors.New("Received error code " + (code.Code).String() + ", message is: " + code.Message)
}
