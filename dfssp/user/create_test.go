package user_test

import (
	"dfss/dfssp/api"
	"dfss/net"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"testing"
)

const (
	// ValidServ is a host/port adress to a platform server with bad setup
	ValidServ = "localhost:9090"
	// InvalidServ is a host/port adress to a platform server with bad setup
	InvalidServ = "localhost:9091"
)

func clientTest(t *testing.T, hostPort string) api.PlatformClient {
	conn, err := net.Connect(hostPort, nil, nil, rootCA, nil)
	if err != nil {
		t.Fatal("Unable to connect: ", err)
	}

	return api.NewPlatformClient(conn)
}

func TestWrongRegisterRequest(t *testing.T) {
	client := clientTest(t, ValidServ)

	request := &api.RegisterRequest{}
	errCode, err := client.Register(context.Background(), request)
	assert.Equal(t, nil, err)
	assert.Equal(t, errCode.Code, api.ErrorCode_INVARG)

	request.Email = "foo"
	errCode, err = client.Register(context.Background(), request)
	assert.Equal(t, nil, err)
	assert.Equal(t, errCode.Code, api.ErrorCode_INVARG)

	request.Request = "foo"
	errCode, err = client.Register(context.Background(), request)
	assert.Equal(t, nil, err)
	assert.Equal(t, errCode.Code, api.ErrorCode_INVARG)
}

func TestAuthUserNotFound(t *testing.T) {
	mail := "wrong@wrong.wrong"
	token := "wrong"
	client := clientTest(t, ValidServ)

	request := &api.AuthRequest{Email: mail, Token: token}
	msg, err := client.Auth(context.Background(), request)
	assert.Equal(t, (*api.RegisteredUser)(nil), msg)
	if err == nil {
		t.Fatal("The request user should not have been found in the database")
	}
}
