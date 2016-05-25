package user_test

import (
	"dfss/dfssp/api"
	"testing"
	"time"

	"dfss/dfssp/entities"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestSimpleRegister(t *testing.T) {
	client := clientTest(t, ValidServ)

	request := &api.RegisterRequest{
		Email:   "simple@simple.simple",
		Request: string(csr),
	}
	errCode, err := client.Register(context.Background(), request)
	assert.Nil(t, err)
	assert.Equal(t, errCode.Code, api.ErrorCode_SUCCESS)
}

func TestWrongRegisterRequest(t *testing.T) {
	client := clientTest(t, ValidServ)

	request := &api.RegisterRequest{}
	errCode, err := client.Register(context.Background(), request)
	assert.Nil(t, err)
	assert.Equal(t, errCode.Code, api.ErrorCode_INVARG)

	// Wrong email, good request
	request.Email = "foo"
	request.Request = string(csr)
	errCode, err = client.Register(context.Background(), request)
	assert.Nil(t, err)
	assert.Equal(t, errCode.Code, api.ErrorCode_INVARG)

	// Good email, wrong request
	request.Email = "foo@foo.foo"
	request.Request = "foo"
	errCode, err = client.Register(context.Background(), request)
	assert.Nil(t, err)
	assert.Equal(t, errCode.Code, api.ErrorCode_INVARG)
}

// An entry already exists with the same email and a very close registration date
// -> Unable to register, even with a different case
func TestRegisterTwice(t *testing.T) {
	user := entities.NewUser()
	user.Email = "twice@twice.twice"

	_, err = repository.Collection.Insert(*user)
	assert.Nil(t, err)

	client := clientTest(t, ValidServ)

	request := &api.RegisterRequest{Email: "twice@twice.twIce", Request: string(csr)}
	errCode, err := client.Register(context.Background(), request)
	assert.Equal(t, err, nil)
	assert.Equal(t, errCode.Code, api.ErrorCode_INVARG)
}

// An entry already exists with the same email, BUT expiration is behind us
// -> Able to register
func TestRegisterRenew(t *testing.T) {
	user := entities.NewUser()
	user.Email = "renew@renew.renew"
	user.Registration = time.Now().AddDate(0, 0, -2)
	user.Expiration = time.Now().Add(-1 * time.Hour)

	_, err = repository.Collection.Insert(*user)
	assert.Nil(t, err)

	client := clientTest(t, ValidServ)

	request := &api.RegisterRequest{Email: "renew@renew.renew", Request: string(csr)}
	errCode, err := client.Register(context.Background(), request)
	assert.Equal(t, err, nil)
	assert.Equal(t, errCode.Code, api.ErrorCode_SUCCESS)
}
