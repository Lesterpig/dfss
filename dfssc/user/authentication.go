package user

import (
	"errors"
	"io/ioutil"
	"regexp"
	"time"

	"dfss/dfssc/common"
	"dfss/dfssc/security"
	pb "dfss/dfssp/api"
	"dfss/net"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// AuthManager handles the authentication of a user
type AuthManager struct {
	viper *viper.Viper
	mail  string
	token string
}

// NewAuthManager creates a new authentication manager with the given parameters
func NewAuthManager(mail, token string, viper *viper.Viper) (*AuthManager, error) {
	m := &AuthManager{
		viper: viper,
		mail:  mail,
		token: token,
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
	fileCert := m.viper.GetString("file_cert")
	if b := common.FileExists(fileCert); b {
		return errors.New("A certificate is already present at path " + fileCert)
	}

	fileCA := m.viper.GetString("file_ca")
	if b := common.FileExists(fileCA); !b {
		return errors.New("You need the certificate of the platform at path " + fileCA)
	}

	data, err := security.GetCertificate(fileCA)
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
	client, err := connect()
	if err != nil {
		return nil, err
	}

	// gRPC request
	request := &pb.AuthRequest{
		Email: m.mail,
		Token: m.token,
	}

	// Stop the context if it takes too long for the platform to answer
	ctx, cancel := context.WithTimeout(context.Background(), net.DefaultTimeout)
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

	return ioutil.WriteFile(m.viper.GetString("file_cert"), cert, 0600)
}
