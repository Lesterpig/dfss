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

var collection *mgdb.MongoCollection
var manager *mgdb.MongoManager
var dbURI string

var repository *entities.ContractRepository

func TestMain(m *testing.M) {

	dbURI = os.Getenv("DFSS_MONGO_URI")
	if dbURI == "" {
		dbURI = "mongodb://localhost/dfss-test"
	}

	var err error
	manager, err = mgdb.NewManager(dbURI)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	collection = manager.Get("contracts")
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

	c.AddSigner(nil, "mail1", []byte{0xaa})
	c.AddSigner(&id, "mail2", []byte{})

	signers := c.Signers

	if len(signers) != 2 {
		t.Fatal("Signers are not inserted correctly")
	}

	assert.Equal(t, signers[0].Email, "mail1")
	assert.Equal(t, signers[0].Hash, []byte{0xaa})
	assert.Equal(t, signers[0].UserID.Hex(), "000000000000000000000000")
	assert.Equal(t, signers[1].Email, "mail2")
	assert.Equal(t, signers[1].Hash, []byte{})
	assert.Equal(t, signers[1].UserID.Hex(), id.Hex())
}

func TestGetHashChain(t *testing.T) {
	c := entities.NewContract()
	c.AddSigner(nil, "mail1", []byte{0xaa})
	c.AddSigner(nil, "mail2", []byte{0xbb, 0xcc})
	c.AddSigner(nil, "mail3", []byte{})

	chain := c.GetHashChain()
	assert.Equal(t, 3, len(chain))
	assert.Equal(t, []byte{0xaa}, chain[0])
	assert.Equal(t, []byte{0xbb, 0xcc}, chain[1])
	assert.Equal(t, []byte{}, chain[2])
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
	dropDataset()
	c := entities.NewContract()
	c.AddSigner(nil, "mail1", []byte{0xaa})
	c.AddSigner(nil, "mail1", []byte{0xaa})
	c.File.Name = "file"
	c.File.Hash = []byte{0xff}
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

// Insert some contracts with missing user and test waiting contracts for this user
func TestGetWaitingForUser(t *testing.T) {
	knownID := bson.NewObjectId()
	dropDataset()

	c1 := entities.NewContract()
	c1.AddSigner(nil, "mail1", []byte{})
	c1.Ready = false

	c2 := entities.NewContract()
	c2.AddSigner(nil, "mail1", []byte{})
	c2.AddSigner(&knownID, "mail2", []byte{0x12})
	c2.Ready = false

	c3 := entities.NewContract()
	c3.AddSigner(nil, "mail2", []byte{})
	c3.AddSigner(&knownID, "mail1", []byte{0xaa})
	c3.Ready = false

	_, _ = repository.Collection.Insert(c1)
	_, _ = repository.Collection.Insert(c2)
	_, _ = repository.Collection.Insert(c3)

	contracts, err := repository.GetWaitingForUser("mail1")
	assert.Equal(t, nil, err)
	assert.Equal(t, 2, len(contracts))
}

func TestCheckAuthorization(t *testing.T) {
	dropDataset()
	createDataset()
	id := addTestContract()

	res, err := repository.GetWithSigner(user1.CertHash, id)
	assert.Equal(t, nil, err)
	assert.T(t, res != nil)
	res, err = repository.GetWithSigner(user1.CertHash, bson.NewObjectId())
	assert.Equal(t, nil, err)
	assert.T(t, res == nil)
	res, err = repository.GetWithSigner(user2.CertHash, id)
	assert.Equal(t, nil, err)
	assert.T(t, res == nil)
	res, err = repository.GetWithSigner(user2.CertHash, bson.NewObjectId())
	assert.Equal(t, nil, err)
	assert.T(t, res == nil)

	contract := entities.Contract{}
	_ = repository.Collection.FindByID(entities.Contract{ID: id}, &contract)
	contract.Ready = false
	_, _ = repository.Collection.UpdateByID(contract)

	// Still valid if contract is not ready
	res, _ = repository.GetWithSigner(user1.CertHash, id)
	assert.T(t, res != nil)
	assert.T(t, !res.Ready)
}
