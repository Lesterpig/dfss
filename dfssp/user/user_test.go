package user_test

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"dfss/auth"
	"dfss/dfssp/api"
	"dfss/dfssp/entities"
	"dfss/dfssp/server"
	u "dfss/dfssp/user"
	"dfss/mgdb"
	"dfss/net"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2/bson"
)

var (
	mail          string
	csr           []byte
	rootCA        *x509.Certificate
	rootKey, pkey *rsa.PrivateKey
)

const (
	// ValidServ is a host/port adress to a platform server with good setup
	ValidServ = "localhost:9090"
)

func clientTest(t *testing.T, hostPort string) api.PlatformClient {
	conn, err := net.Connect(hostPort, nil, nil, rootCA, nil)
	if err != nil {
		t.Fatal("Unable to connect: ", err)
	}

	return api.NewPlatformClient(conn)
}

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
	viper.Set("ca_filename", "dfssp_rootCA.pem")
	viper.Set("pkey_filename", "dfssp_pkey.pem")

	dbURI = os.Getenv("DFSS_MONGO_URI")
	if dbURI == "" {
		dbURI = "mongodb://localhost/dfss-test"
	}

	manager, err = mgdb.NewManager(dbURI)
	collection = manager.Get("users")

	repository = entities.NewUserRepository(collection)

	keyPath := filepath.Join(os.Getenv("GOPATH"), "src", "dfss", "dfssp", "testdata")

	// Valid server
	viper.Set("path", keyPath)
	viper.Set("dbURI", dbURI)
	viper.Set("validity", 365)
	viper.Set("verbose", true)
	srv := server.GetServer()
	go func() { log.Fatal(net.Listen(ValidServ, srv)) }()

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
	user.Csr = "csr1"
	user.RegToken = "regToken 1"

	_, err = repository.Collection.Insert(user)
	if err != nil {
		t.Fatal("An error occurred while inserting the user")
	}
}

func TestUnregisterUser(t *testing.T) {
	var hash = []byte{0xde, 0xad, 0xbe, 0xef}
	user := entities.NewUser()
	user.Email = "dfss2@mpcs.tk"
	user.CertHash = hash
	user.Csr = "csr2"
	user.RegToken = "regToken 2"

	_, err = repository.Collection.Insert(user)
	if err != nil {
		t.Fatal("An error occurred while inserting the user")
	}

	response := u.Unregister(manager, hash)
	if response.Code != api.ErrorCode_SUCCESS {
		t.Fatal("An error occured while deleting the user:" + response.Message)
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

	if user1.Certificate != user2.Certificate {
		t.Fatal("Certificate doesn't match : received ", user1.Certificate, " and ", user2.Certificate)
	}
}

func TestMongoFetchUser(t *testing.T) {
	user := entities.NewUser()
	user.Email = "dfss2@mpcs.tk"
	user.CertHash = nil
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

	conn, err := net.Connect("localhost:9090", nil, nil, rootCA, nil)
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

	fmt.Println("Certificate successfully received")

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
	// Certificate successfully received
	// Database successfully updated with cert and certHash
}
