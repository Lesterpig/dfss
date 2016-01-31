package contract_test // Using another package to avoid import cycles

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"dfss/dfssp/entities"
	"dfss/dfssp/server"
	"dfss/mgdb"
	"dfss/net"
	"github.com/bmizerany/assert"
	"gopkg.in/mgo.v2/bson"
)

var err error
var collection *mgdb.MongoCollection
var manager *mgdb.MongoManager
var dbURI string

var repository *entities.ContractRepository

func TestMain(m *testing.M) {

	dbURI = os.Getenv("DFSS_MONGO_URI")
	if dbURI == "" {
		dbURI = "mongodb://localhost/dfss-test"
	}

	manager, err = mgdb.NewManager(dbURI)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	collection = manager.Get("demo")
	repository = entities.NewContractRepository(collection)

	// Start platform server
	keyPath := filepath.Join(os.Getenv("GOPATH"), "src", "dfss", "dfssp", "testdata")

	srv := server.GetServer(keyPath, dbURI, 365, true)
	go func() { _ = net.Listen("localhost:9090", srv) }()

	// Run
	code := m.Run()

	dropDataset()
	manager.Close()
	os.Exit(code)
}

func TestAddSigner(t *testing.T) {
	c := entities.NewContract()

	id := bson.NewObjectId()

	c.AddSigner(nil, "mail1", "hash1")
	c.AddSigner(&id, "mail2", "hash2")

	signers := c.Signers

	if len(signers) != 2 {
		t.Fatal("Signers are not inserted correctly")
	}

	assert.Equal(t, signers[0].Email, "mail1")
	assert.Equal(t, signers[0].Hash, "hash1")
	assert.Equal(t, signers[0].UserID.Hex(), "000000000000000000000000")
	assert.Equal(t, signers[1].Email, "mail2")
	assert.Equal(t, signers[1].Hash, "hash2")
	assert.Equal(t, signers[1].UserID.Hex(), id.Hex())
}

func assertContractEqual(t *testing.T, contract, fetched entities.Contract) {
	assert.Equal(t, contract.File, fetched.File)
	assert.Equal(t, contract.Date.Unix(), fetched.Date.Unix())
	assert.Equal(t, contract.Comment, fetched.Comment)
	assert.Equal(t, contract.Ready, fetched.Ready)
	assert.Equal(t, contract.Signers, fetched.Signers)
}

// Insert a contract with 2 users and check the fields are correctly persisted
func TestInsertContract(t *testing.T) {
	c := entities.NewContract()
	c.AddSigner(nil, "mail1", "hash1")
	c.AddSigner(nil, "mail1", "hash1")
	c.File.Name = "file"
	c.File.Hash = "hashFile"
	c.File.Hosted = false
	c.Comment = "comment"
	c.Ready = true

	_, err := repository.Collection.Insert(c)
	if err != nil {
		t.Fatal("Contract not inserted successfully in the database:", err)
	}

	fetched := entities.Contract{}
	selector := entities.Contract{
		ID: c.ID,
	}
	err = repository.Collection.FindByID(selector, &fetched)

	if err != nil {
		t.Fatal("Contract could not be successfully retrieved:", err)
	}

	assertContractEqual(t, *c, fetched)
}
