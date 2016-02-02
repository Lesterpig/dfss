package user

import (
	"fmt"
	"os"
	"testing"

	"dfss/dfssp/entities"
	"dfss/mgdb"
	"gopkg.in/mgo.v2/bson"
)

var err error
var collection *mgdb.MongoCollection
var manager *mgdb.MongoManager
var dbURI string

var repository *entities.UserRepository

func TestMain(m *testing.M) {

	dbURI = os.Getenv("DFSS_MONGO_URI")
	if dbURI == "" {
		dbURI = "mongodb://localhost/dfss"
	}

	manager, err = mgdb.NewManager(dbURI)
	collection = manager.Get("demo")

	repository = entities.NewUserRepository(collection)

	// Run
	code := m.Run()

	// Teardown
	// The collection is created automatically on
	// first connection, that's why we do not recreate it manually
	err = collection.Drop()

	if err != nil {
		fmt.Println("An error occurred while droping the collection")
	}
	manager.Close()

	os.Exit(code)
}

func TestMongoFetchInexistantUser(t *testing.T) {
	user, erro := repository.FetchByMailAndHash("dummyMail", "dummyHash")
	if user != nil || erro != nil {
		t.Fatal("User should not have been found and error should  be nil")
	}
}

func TestMongoInsertUser(t *testing.T) {
	user := entities.NewUser()
	user.Email = "dfss1@mpcs.tk"
	user.CertHash = "dummy_hash"
	user.ConnInfo.IP = "127.0.0.1"
	user.ConnInfo.Port = 1111
	user.Csr = "csr1"
	user.RegToken = "regToken 1"

	_, err = repository.Collection.Insert(user)
	if err != nil {
		t.Fatal("An error occurred while inserting the user")
	}

	fmt.Println("Successfully inserted a user")
}

func equalUsers(t *testing.T, user1, user2 *entities.User) {
	if user1.ID != user2.ID {
		t.Fatal("ID doesn't match : received ", user1.ID, " and ", user2.ID)
	}

	if user1.CertHash != user2.CertHash {
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
	user.CertHash = "dummy_hash"
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
