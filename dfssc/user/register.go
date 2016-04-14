package user

import (
	"errors"
	"regexp"
	"time"

	"dfss/dfssc/common"
	"dfss/dfssc/security"
	pb "dfss/dfssp/api"
	"dfss/net"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// RegisterManager handles the registration of a user
type RegisterManager struct {
	fileCA       string
	fileCert     string
	fileKey      string
	addrPort     string
	passphrase   string
	country      string
	organization string
	unit         string
	mail         string
	bits         int
}

// NewRegisterManager return a new Register Manager to register a user
func NewRegisterManager(fileCA, fileCert, fileKey, addrPort, passphrase, country, organization, unit, mail string, bits int) (*RegisterManager, error) {
	m := &RegisterManager{
		fileCA:       fileCA,
		fileCert:     fileCert,
		fileKey:      fileKey,
		addrPort:     addrPort,
		passphrase:   passphrase,
		country:      country,
		organization: organization,
		unit:         unit,
		mail:         mail,
		bits:         bits,
	}

	if err := m.checkValidParams(); err != nil {
		return nil, err
	}

	if err := m.checkFilePresence(); err != nil {
		return nil, err
	}

	return m, nil
}

// Check the validity of the provided email, passphrase and bits
func (m *RegisterManager) checkValidParams() error {
	re, _ := regexp.Compile(`.+@.+\..+`)
	if b := re.MatchString(m.mail); !b {
		return errors.New("Provided mail is not valid")
	}

	if m.bits != 2048 && m.bits != 4096 {
		return errors.New("Length of the key should be 2048 or 4096 bits")
	}

	return nil
}

// Check there is no private key nor client certificate
// Check the CA is present and valid
// Check there is not a duplicate file
func (m *RegisterManager) checkFilePresence() error {
	if b := common.FileExists(m.fileKey); b {
		return errors.New("A private key is already present at path " + m.fileKey)
	}

	if b := common.FileExists(m.fileCert); b {
		return errors.New("A certificate is already present at path " + m.fileCert)
	}

	if m.fileKey == m.fileCert {
		return errors.New("Cannot store certificate and key in the same file")
	}

	if b := common.FileExists(m.fileCA); !b {
		return errors.New("You need the certificate of the platform at path " + m.fileCA)
	}

	data, err := security.GetCertificate(m.fileCA)
	if err != nil {
		return err
	}

	if time.Now().After(data.NotAfter) {
		return errors.New("Root certificate has expired")
	}

	return nil
}

// GetCertificate handles the creation of a certificate, delete private key upon failure
func (m *RegisterManager) GetCertificate() error {
	request, err := m.buildCertificateRequest()
	if err != nil {
		common.DeleteQuietly(m.fileKey)
		return err
	}

	code, err := m.sendRequest(request)
	if err != nil {
		common.DeleteQuietly(m.fileKey)
		return err
	}

	err = common.EvaluateErrorCodeResponse(code)
	if err != nil {
		common.DeleteQuietly(m.fileKey)
		return err
	}

	return nil
}

// Builds a certificate request
func (m *RegisterManager) buildCertificateRequest() (string, error) {
	key, err := security.GenerateKeys(m.bits, m.passphrase, m.fileKey)
	if err != nil {
		return "", err
	}

	request, err := security.GenerateCertificateRequest(m.country, m.organization, m.unit, m.mail, key)
	if err != nil {
		return "", err
	}

	return request, nil
}

// Send the request and returns the response
func (m *RegisterManager) sendRequest(certRequest string) (*pb.ErrorCode, error) {
	client, err := connect(m.fileCA, m.addrPort)
	if err != nil {
		return nil, err
	}

	// gRPC request
	request := &pb.RegisterRequest{
		Email:   m.mail,
		Request: certRequest,
	}

	// Stop the context if it takes too long for the platform to answer
	ctx, cancel := context.WithTimeout(context.Background(), net.DefaultTimeout)
	defer cancel()
	response, err := client.Register(ctx, request)
	if err != nil {
		return nil, errors.New(grpc.ErrorDesc(err))
	}

	return response, nil
}
