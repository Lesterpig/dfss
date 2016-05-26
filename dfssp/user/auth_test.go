package user_test

import (
	"testing"
	"time"

	"dfss/dfssp/api"
	"dfss/dfssp/entities"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

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

func TestAuthTwice(t *testing.T) {
	email := "email"
	token := "token"
	user := entities.NewUser()
	user.Email = email
	user.RegToken = token
	user.Csr = string(csr)
	user.Certificate = "foo"
	user.CertHash = []byte{0xaa}

	_, err = repository.Collection.Insert(*user)
	if err != nil {
		t.Fatal(err)
	}

	// User is already registered
	client := clientTest(t, ValidServ)
	request := &api.AuthRequest{Email: email, Token: token}
	msg, err := client.Auth(context.Background(), request)
	assert.Equal(t, msg, (*api.RegisteredUser)(nil))
	if err == nil {
		t.Fatal("The user should have been evaluated as already registered")
	}
}

func TestWrongAuthRequestContext(t *testing.T) {
	mail := "right@right.right"
	token := "right"

	user := entities.NewUser()
	user.Email = mail
	user.RegToken = token
	user.Registration = time.Now().UTC().Add(time.Hour * -48)

	_, err := repository.Collection.Insert(*user)
	if err != nil {
		t.Fatal(err)
	}

	client := clientTest(t, ValidServ)

	request := &api.AuthRequest{Email: mail, Token: "foo"}

	// Token timeout
	msg, err := client.Auth(context.Background(), request)
	assert.Equal(t, (*api.RegisteredUser)(nil), msg)
	if err == nil {
		t.Fatal("The request should have been evaluated as invalid")
	}

	// Token mismatch
	user.Registration = time.Now().UTC()

	_, err = repository.Collection.UpdateByID(*user)
	if err != nil {
		t.Fatal(err)
	}

	msg, err = client.Auth(context.Background(), request)
	assert.Equal(t, (*api.RegisteredUser)(nil), msg)
	if err == nil {
		t.Fatal("The request should have been evaluated as invalid")
	}

	res := entities.User{}
	err = repository.Collection.FindByID(*user, &res)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, res.Certificate, "")
	assert.Equal(t, res.CertHash, []byte{})

	// Invalid certificate request (none here)
	request.Token = token
	msg, err = client.Auth(context.Background(), request)
	assert.Equal(t, (*api.RegisteredUser)(nil), msg)
	if err == nil {
		t.Fatal("The request should have been evaluated as invalid")
	}

	err = repository.Collection.FindByID(*user, &res)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, res.Certificate, "")
	assert.Equal(t, res.CertHash, []byte{})
}
