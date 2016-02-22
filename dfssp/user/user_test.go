package user_test

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"dfss/auth"
	"dfss/dfssp/api"
	"dfss/dfssp/entities"
	"dfss/dfssp/server"
	"dfss/mgdb"
	"dfss/net"
	"github.com/bmizerany/assert"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2/bson"
	"time"
)

var (
	mail          string
	csr           []byte
	rootCA        *x509.Certificate
	rootKey, pkey *rsa.PrivateKey
)

func init() {
	mail = "foo@foo.foo"
	pkey, _ = auth.GeneratePrivateKey(512)

	path := filepath.Join(os.Getenv("GOPATH"), "src", "dfss", "dfssp", "testdata", "dfssp_rootCA.pem")
	CAData, _ := ioutil.ReadFile(path)

	rootCA, _ = auth.PEMToCertificate(CAData)

	path = filepath.Join(os.Getenv("GOPATH"), "src", "dfss", "dfssp", "testdata", "dfssp_pkey.pem")
	KeyData, _ := ioutil.ReadFile(path)

	rootKey, _ = auth.PEMToPrivateKey(KeyData)

	csr, _ = auth.GetCertificateRequest("country", "organization", "unit", mail, pkey)
}

var err error
var collection *mgdb.MongoCollection
var manager *mgdb.MongoManager
var dbURI string

var repository *entities.UserRepository

func TestMain(m *testing.M) {

	dbURI = os.Getenv("DFSS_MONGO_URI")
	if dbURI == "" {
		dbURI = "mongodb://localhost/dfss-test"
	}

	manager, err = mgdb.NewManager(dbURI)
	collection = manager.Get("users")

	repository = entities.NewUserRepository(collection)

	keyPath := filepath.Join(os.Getenv("GOPATH"), "src", "dfss", "dfssp", "testdata")

	// Valid server
	srv := server.GetServer(keyPath, dbURI, 365, true)
	go func() { _ = net.Listen(ValidServ, srv) }()

	// Server using invalid certificate duration
	srv2 := server.GetServer(keyPath, dbURI, -1, true)
	go func() { _ = net.Listen(InvalidServ, srv2) }()

	// Run
	err = collection.Drop()
	code := m.Run()
	err = collection.Drop()

	if err != nil {
		fmt.Println("An error occurred while droping the collection")
	}
	manager.Close()

	os.Exit(code)
}

func TestMongoFetchInexistantUser(t *testing.T) {
	user, erro := repository.FetchByMailAndHash("dummyMail", []byte{0x01})
	if user != nil || erro != nil {
		t.Fatal("User should not have been found and error should  be nil")
	}
}

func TestMongoInsertUser(t *testing.T) {
	user := entities.NewUser()
	user.Email = "dfss1@mpcs.tk"
	user.CertHash = []byte{0x01, 0x02}
	user.ConnInfo.IP = "127.0.0.1"
	user.ConnInfo.Port = 1111
	user.Csr = "csr1"
	user.RegToken = "regToken 1"

	_, err = repository.Collection.Insert(user)
	if err != nil {
		t.Fatal("An error occurred while inserting the user")
	}
}

func equalUsers(t *testing.T, user1, user2 *entities.User) {
	if user1.ID != user2.ID {
		t.Fatal("ID doesn't match : received ", user1.ID, " and ", user2.ID)
	}

	if string(user1.CertHash) != string(user2.CertHash) {
		t.Fatal("CertHash doesn't match : received ", user1.CertHash, " and ", user2.CertHash)
	}

	if user1.Email != user2.Email {
		t.Fatal("Email doesn't match : received ", user1.Email, " and ", user2.Email)
	}

	if user1.Registration.Unix() != user2.Registration.Unix() {
		t.Fatal("Registration doesn't match : received ", user1.Registration, " and ", user2.Registration)
	}

	if user1.RegToken != user2.RegToken {
		t.Fatal("RegToken doesn't match : received ", user1.RegToken, " and ", user2.RegToken)
	}

	if user1.Csr != user2.Csr {
		t.Fatal("Csr doesn't match : received ", user1.Csr, " and ", user2.Csr)
	}

	if user1.ConnInfo.IP != user2.ConnInfo.IP {
		t.Fatal("ConnInfo.IP doesn't match : received ", user1.ConnInfo.IP, " and ", user2.ConnInfo.IP)
	}

	if user1.ConnInfo.Port != user2.ConnInfo.Port {
		t.Fatal("ConnInfo.Port doesn't match : received ", user1.ConnInfo.Port, " and ", user2.ConnInfo.Port)
	}

	if user1.Certificate != user2.Certificate {
		t.Fatal("Certificate doesn't match : received ", user1.Certificate, " and ", user2.Certificate)
	}
}

func TestMongoFetchUser(t *testing.T) {
	user := entities.NewUser()
	user.Email = "dfss2@mpcs.tk"
	user.CertHash = nil
	user.ConnInfo.IP = "127.0.0.2"
	user.ConnInfo.Port = 2222
	user.Csr = "csr2"
	user.RegToken = "regToken 2"

	_, err = repository.Collection.Insert(user)

	if err != nil {
		t.Fatal("An error occured while inserting a user: ", err)
	}

	fetched, erro := repository.FetchByMailAndHash(user.Email, user.CertHash)

	if erro != nil {
		t.Fatal("An error occurred while fetching the previously inserted user", err)
	}

	equalUsers(t, user, fetched)
}

func TestMongoFetchIncompleteUser(t *testing.T) {
	user := entities.User{
		ID: bson.NewObjectId(),
	}

	_, err = repository.Collection.Insert(user)

	if err != nil {
		t.Fatal("An error occured while inserting a user: ", err)
	}

	fetched, erro := repository.FetchByMailAndHash(user.Email, user.CertHash)

	if erro != nil {
		t.Fatal("An error occurred while fetching the previously inserted user", err)
	}

	equalUsers(t, &user, fetched)
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

func ExampleAuth() {
	mail := "example@example.example"
	token := "example"

	user := entities.NewUser()
	user.Email = mail
	user.RegToken = token
	user.Csr = string(csr)

	_, err = repository.Collection.Insert(*user)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("User successfully inserted")

	conn, err := net.Connect("localhost:9090", nil, nil, rootCA)
	if err != nil {
		fmt.Println("Unable to connect: ", err)
	}
	fmt.Println("Client successfully connected")

	client := api.NewPlatformClient(conn)

	request := &api.AuthRequest{Email: user.Email, Token: user.RegToken}
	msg, err := client.Auth(context.Background(), request)

	fmt.Println("AuthRequest successfully sent")

	if msg == (*api.RegisteredUser)(nil) {
		fmt.Println("The request should have been evaluated as valid")
	}
	if err != nil {
		fmt.Println(err)
	}

	if msg.ClientCert == "" {
		fmt.Println("The certificate should have been given as an answer")
	}

	fmt.Println("Certificate successfully recieved")

	res := entities.User{}
	err = repository.Collection.FindByID(*user, &res)
	if err != nil {
		fmt.Println(err)
	}

	if res.Certificate == "" || res.CertHash == nil {
		fmt.Println("The database should have been updated")
	}

	fmt.Println("Database successfully updated with cert and certHash")

	// Output:
	// User successfully inserted
	// Client successfully connected
	// AuthRequest successfully sent
	// Certificate successfully recieved
	// Database successfully updated with cert and certHash
}

func TestRegisterTwice(t *testing.T) {
	mail := "done@done.done"

	user := entities.NewUser()
	user.Email = mail

	_, err = repository.Collection.Insert(*user)
	if err != nil {
		fmt.Println(err)
	}

	client := clientTest(t, ValidServ)

	// An entry already exists with the same email
	request := &api.RegisterRequest{Email: mail, Request: string(csr)}
	errCode, err := client.Register(context.Background(), request)
	assert.Equal(t, err, nil)
	assert.Equal(t, errCode.Code, api.ErrorCode_INVARG)
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
