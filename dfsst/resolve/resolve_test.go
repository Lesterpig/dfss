package resolve

import (
	"fmt"
	"os"
	"testing"

	"crypto/sha512"
	"dfss/dfsst/entities"
	"dfss/mgdb"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

var (
	db        string
	dbManager *mgdb.MongoManager

	sequence             []uint32
	signers              [][]byte
	contractDocumentHash []byte
	signatureUUIDBson    bson.ObjectId

	signedHash []byte

	signersEntities []entities.Signer

	err error
)

func init() {
	db = os.Getenv("DFSS_MONGO_URI")
	if db == "" {
		db = "mongodb://localhost/dfss-test"
	}
	sequence = []uint32{0, 1, 2, 0, 1, 2, 0, 1, 2}

	for i := 0; i < 3; i++ {
		h := sha512.Sum512([]byte{byte(i)})
		signer := h[:]
		signers = append(signers, signer)
	}

	contractDocumentHash = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}
	signatureUUIDBson = bson.NewObjectId()

	signersEntities = make([]entities.Signer, 0)
	for _, s := range signers {
		signerEntity := entities.NewSigner(s)
		signersEntities = append(signersEntities, *signerEntity)
	}
}

func TestMain(m *testing.M) {
	dbManager, err = mgdb.NewManager(db)
	if err != nil {
		fmt.Fprintln(os.Stderr, "An error occured during the connection to MongoDB:", err)
		os.Exit(2)
	}

	collection := dbManager.Get("signatures")
	err = collection.Drop()
	code := m.Run()
	err = collection.Drop()

	if err != nil {
		fmt.Println("An error occurred while droping the collection")
	}

	dbManager.Close()

	os.Exit(code)
}

func TestArePromisesComplete(t *testing.T) {
	// TODO
	// This requires the implementation of the call to the ttp
}

func TestSolve(t *testing.T) {
	archives := entities.NewSignatureArchives(signatureUUIDBson, sequence, signersEntities, contractDocumentHash, signedHash)
	manager := &entities.ArchivesManager{
		DB:       dbManager,
		Archives: archives,
	}

	promise0 := &entities.Promise{
		RecipientKeyIndex: 1,
		SenderKeyIndex:    0,
	}
	promise1 := &entities.Promise{
		RecipientKeyIndex: 2,
		SenderKeyIndex:    1,
	}
	promise2 := &entities.Promise{
		RecipientKeyIndex: 1,
		SenderKeyIndex:    2,
	}

	ok, contract := Solve(manager)
	assert.Equal(t, ok, false)
	assert.Equal(t, len(contract), 0)

	manager.Archives.ReceivedPromises = append(manager.Archives.ReceivedPromises, *promise0)
	ok, contract = Solve(manager)
	assert.Equal(t, ok, false)
	assert.Equal(t, len(contract), 0)

	manager.Archives.ReceivedPromises = append(manager.Archives.ReceivedPromises, *promise1)
	ok, contract = Solve(manager)
	assert.Equal(t, ok, false)
	assert.Equal(t, len(contract), 0)

	manager.Archives.ReceivedPromises = append(manager.Archives.ReceivedPromises, *promise2)
	ok, contract = Solve(manager)
	assert.Equal(t, ok, true)
	if len(contract) == 0 {
		t.Fatal("Contract should have need generated")
	}
}

// TO MODIFY WHEN SOURCE FUNCTION WILL BE UPDATED
func TestGenerateSignedContract(t *testing.T) {
	// TODO
	assert.Equal(t, true, true)
}

func TestComputeDishonestSigners(t *testing.T) {
	var tmpPromises []*entities.Promise
	archives := entities.NewSignatureArchives(signatureUUIDBson, sequence, signersEntities, contractDocumentHash, signedHash)

	assert.Equal(t, len(archives.AbortedSigners), 0)

	dishonestSigners := ComputeDishonestSigners(archives, tmpPromises)
	assert.Equal(t, len(dishonestSigners), 0)

	abortedSigner := entities.NewAbortedSigner(uint32(1), uint32(1))
	archives.AbortedSigners = append(archives.AbortedSigners, *abortedSigner)
	assert.Equal(t, len(archives.AbortedSigners), 1)

	dishonestSigners = ComputeDishonestSigners(archives, tmpPromises)
	assert.Equal(t, len(dishonestSigners), 0)

	promise0 := &entities.Promise{
		RecipientKeyIndex: 1,
		SenderKeyIndex:    0,
		SequenceIndex:     0,
	}
	tmpPromises = append(tmpPromises, promise0)
	assert.Equal(t, len(tmpPromises), 1)

	dishonestSigners = ComputeDishonestSigners(archives, tmpPromises)
	assert.Equal(t, len(dishonestSigners), 0)

	promise1 := &entities.Promise{
		RecipientKeyIndex: 0,
		SenderKeyIndex:    1,
		SequenceIndex:     1,
	}
	tmpPromises = append(tmpPromises, promise1)
	assert.Equal(t, len(tmpPromises), 2)

	dishonestSigners = ComputeDishonestSigners(archives, tmpPromises)
	assert.Equal(t, len(dishonestSigners), 0)

	promise2 := &entities.Promise{
		RecipientKeyIndex: 0,
		SenderKeyIndex:    1,
		SequenceIndex:     4,
	}
	tmpPromises = append(tmpPromises, promise2)
	assert.Equal(t, len(tmpPromises), 3)

	dishonestSigners = ComputeDishonestSigners(archives, tmpPromises)
	assert.Equal(t, len(dishonestSigners), 1)
	assert.Equal(t, dishonestSigners[0], uint32(1))

	promise3 := &entities.Promise{
		RecipientKeyIndex: 0,
		SenderKeyIndex:    1,
		SequenceIndex:     2,
	}
	tmpPromises = append(tmpPromises, promise3)
	assert.Equal(t, len(tmpPromises), 4)

	dishonestSigners = ComputeDishonestSigners(archives, tmpPromises)
	assert.Equal(t, len(dishonestSigners), 1)
	assert.Equal(t, dishonestSigners[0], uint32(1))
}
