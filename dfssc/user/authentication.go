package user

import (
	"io/ioutil"
	"regexp"
	"time"

	"dfss/dfssc/common"
	"dfss/dfssc/security"
	pb "dfss/dfssp/api"
	"errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// AuthManager handles the authentication of a user
type AuthManager struct {
	fileCA   string
	fileCert string
	addrPort string
	mail     string
	token    string
}

// NewAuthManager creates a new authentication manager with the given parameters
func NewAuthManager(fileCA, fileCert, addrPort, mail, token string) (*AuthManager, error) {
	m := &AuthManager{
		fileCA:   fileCA,
		fileCert: fileCert,
		addrPort: addrPort,
		mail:     mail,
		token:    token,
	}

	if err := m.checkValidParams(); err != nil {
		return nil, err
	}

	if err := m.checkFilePresence(); err != nil {
		return nil, err
	}

	return m, nil
}

func (m *AuthManager) checkValidParams() error {
	re, _ := regexp.Compile(`.+@.+\..+`)
	if b := re.MatchString(m.mail); !b {
		return errors.New("Provided mail is not valid")
	}

	return nil
}

func (m *AuthManager) checkFilePresence() error {
	if b := common.FileExists(m.fileCert); b {
		return errors.New("A certificate is already present at path " + m.fileCert)
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

// Authenticate performs the authentication request
// (ie connection to the platform grpc server, sending of the request, handling the response)
func (m *AuthManager) Authenticate() error {
	response, err := m.sendRequest()
	if err != nil {
		return err
	}

	return m.evaluateResponse(response)
}

// Creates the associated authentication request and sends it to the platform grpc server
func (m *AuthManager) sendRequest() (*pb.RegisteredUser, error) {
	client, err := connect(m.fileCA, m.addrPort)
	if err != nil {
		return nil, err
	}

	// gRPC request
	request := &pb.AuthRequest{
		Email: m.mail,
		Token: m.token,
	}

	// Stop the context if it takes too long for the platform to answer
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	response, err := client.Auth(ctx, request)
	if err != nil {
		return nil, errors.New(grpc.ErrorDesc(err))
	}

	return response, nil
}

// Handle the platform grpc server's reponse to the authentication request
func (m *AuthManager) evaluateResponse(response *pb.RegisteredUser) error {
	cert := []byte(response.ClientCert)

	return ioutil.WriteFile(m.fileCert, cert, 0600)
}
