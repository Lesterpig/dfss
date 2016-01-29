package contract

import (
	"dfss/mgdb"
	"fmt"
	"os"
	"testing"
)

var err error
var collection *mgdb.MongoCollection
var manager *mgdb.MongoManager
var dbURI string

var repository *Repository

func TestMain(m *testing.M) {

	dbURI = os.Getenv("DFSS_MONGO_URI")
	if dbURI == "" {
		dbURI = "mongodb://localhost/dfss"
	}

	manager, err = mgdb.NewManager(dbURI)
	collection = manager.Get("demo")

	repository = NewRepository(collection)

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

func TestAddSigner(t *testing.T) {
	contract := NewContract()
	contract.AddSigner("mail1", "hash1")
	contract.AddSigner("mail2", "hash2")

	signers := contract.Signers
	var fixedSigners [2]Signer
	copy(fixedSigners[:], signers[:2])

	if len(signers) != 2 {
		t.Fatal("Signers are not inserted correctly")
	}

	helperCompareEmailAndHash(t, fixedSigners)
}

func helperCompareFiles(t *testing.T, contract, fetched Contract) {
	if contract.File.Name != fetched.File.Name || contract.File.Hash != fetched.File.Hash || contract.File.Hosted != fetched.File.Hosted {
		t.Fatal("Contract files doesn't match")
	}
}
func helperCompareEmailAndHash(t *testing.T, signers [2]Signer) {
	if signers[0].Email != "mail1" && signers[0].Email != "mail2" || signers[1].Email != "mail1" && signers[1].Email != "mail2" {
		t.Fatal("Signers are not correctly added")
	}
	if signers[0].Hash != "hash1" && signers[0].Hash != "hash2" || signers[1].Hash != "hash1" && signers[1].Hash != "hash2" {
		t.Fatal("Signers are not correctly added")
	}
}

func helperCompareInformations(t *testing.T, contract, fetched Contract) {
	if contract.Date.Unix() != fetched.Date.Unix() || contract.Comment != fetched.Comment || contract.Ready != fetched.Ready || len(contract.Signers) != len(fetched.Signers) {
		t.Fatal("Contract informations doesn't match")
	}
}

// Insert a contract with 2 users and check the fields are correctly persisted
func TestInsertContract(t *testing.T) {
	contract := NewContract()
	contract.AddSigner("mail1", "hash1")
	contract.AddSigner("mail1", "hash1")
	contract.File.Name = "file"
	contract.File.Hash = "hashFile"
	contract.File.Hosted = false
	contract.Comment = "comment"
	contract.Ready = true

	_, err := repository.Collection.Insert(contract)
	if err != nil {
		t.Fatal("Contract not inserted successfully in the database")
	}

	fmt.Println("Successfully persisted contract")

	fetched := Contract{}
	selector := Contract{
		ID: contract.ID,
	}
	err = repository.Collection.FindByID(selector, &fetched)

	if err != nil {
		t.Fatal("Contract could not be successfully retrieved")
	}

	helperCompareFiles(t, *contract, fetched)
	helperCompareInformations(t, *contract, fetched)

	if contract.Signers[0].Hash != fetched.Signers[0].Hash || contract.Signers[0].Hash != fetched.Signers[1].Hash {
		t.Fatal("Signers hash doesn't match")
	}
	if contract.Signers[0].Email != fetched.Signers[0].Email || contract.Signers[0].Email != fetched.Signers[1].Email {
		t.Fatal("Signers hash doesn't match")
	}
	if contract.Signers[0].UserID != fetched.Signers[0].UserID || contract.Signers[1].UserID != fetched.Signers[1].UserID {
		t.Fatal("Signers id doesn't match")
	}
	fmt.Println("Successfully fetched contract")
}
