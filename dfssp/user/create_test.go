package user_test

import (
	"dfss/dfssp/api"
	"dfss/net"
	"github.com/bmizerany/assert"
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
	conn, err := net.Connect(hostPort, nil, nil, rootCA)
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

func TestWrongAuthRequest(t *testing.T) {
	// Get a client to the invalid server (cert duration is -1)
	client := clientTest(t, InvalidServ)

	// Invalid mail length
	inv := &api.AuthRequest{}
	msg, err := client.Auth(context.Background(), inv)
	if msg != nil || err == nil {
		t.Fatal("The request should have been evaluated as invalid")
	}

	// Invalid token length
	inv.Email = "foo"
	msg, err = client.Auth(context.Background(), inv)
	if msg != nil || err == nil {
		t.Fatal("The request should have been evaluated as invalid")
	}

	// Invalid certificate validity duration
	inv.Token = "foo"
	msg, err = client.Auth(context.Background(), inv)
	if msg != nil || err == nil {
		t.Fatal("The request should have been evaluated as invalid")
	}
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
