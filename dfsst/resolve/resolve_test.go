package resolve

import (
	"crypto/sha512"
	"fmt"
	"os"
	"testing"

	cAPI "dfss/dfssc/api"
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
	var promises []*entities.Promise
	promise := &cAPI.Promise{
		Context: &cAPI.Context{
			RecipientKeyHash: []byte{},
			Sequence:         sequence,
			Signers:          signers,
		},
	}

	complete := ArePromisesComplete(promises, promise, 2)
	assert.False(t, complete)

	promise.Context.RecipientKeyHash = signers[2]
	complete = ArePromisesComplete(promises, promise, 2)
	assert.False(t, complete)

	for i := 1; i > -1; i-- {
		p := &entities.Promise{
			RecipientKeyIndex: 2,
			SenderKeyIndex:    sequence[i],
			SequenceIndex:     uint32(i),
		}
		promises = append(promises, p)
	}

	complete = ArePromisesComplete(promises, promise, 2)
	assert.False(t, complete)

	selfPromise := &entities.Promise{
		RecipientKeyIndex: 2,
		SenderKeyIndex:    2,
		SequenceIndex:     2,
	}
	promises = append(promises, selfPromise)

	complete = ArePromisesComplete(promises, promise, 2)
	assert.True(t, complete)

	promises = []*entities.Promise{}
	for i := 7; i > 5; i-- {
		p := &entities.Promise{
			RecipientKeyIndex: 2,
			SenderKeyIndex:    sequence[i],
			SequenceIndex:     uint32(i),
		}
		promises = append(promises, p)
	}

	complete = ArePromisesComplete(promises, promise, 8)
	assert.False(t, complete)

	selfPromise = &entities.Promise{
		RecipientKeyIndex: 2,
		SenderKeyIndex:    2,
		SequenceIndex:     8,
	}
	promises = append(promises, selfPromise)
	complete = ArePromisesComplete(promises, promise, 8)
	assert.True(t, complete)
}

func TestGenerateExpectedPromises(t *testing.T) {
	var expected []*entities.Promise
	promise := &cAPI.Promise{
		Context: &cAPI.Context{
			RecipientKeyHash: []byte{},
			Sequence:         sequence,
			Signers:          signers,
		},
	}

	promises, err := generateExpectedPromises(promise, 1)
	assert.Equal(t, err.Error(), "Signer's hash couldn't be matched")
	assert.Nil(t, promises)

	promise.Context.RecipientKeyHash = signers[1]
	promises, err = generateExpectedPromises(promise, 1)
	assert.Nil(t, err)
	assert.Equal(t, len(promises), 2)
	assert.True(t, promises[0].Equal(&entities.Promise{
		RecipientKeyIndex: 1,
		SenderKeyIndex:    0,
		SequenceIndex:     0,
	}))
	assert.True(t, promises[1].Equal(&entities.Promise{
		RecipientKeyIndex: 1,
		SenderKeyIndex:    1,
		SequenceIndex:     1,
	}))

	promise.Context.RecipientKeyHash = signers[2]
	promises, err = generateExpectedPromises(promise, 1)
	assert.Equal(t, err.Error(), "Signer at step is not recipient")
	assert.Nil(t, promises)

	promises, err = generateExpectedPromises(promise, 2)
	assert.Nil(t, err)

	for i := 1; i > -1; i-- {
		p := &entities.Promise{
			RecipientKeyIndex: 2,
			SenderKeyIndex:    sequence[i],
			SequenceIndex:     uint32(i),
		}
		expected = append(expected, p)
	}
	selfPromise := &entities.Promise{
		RecipientKeyIndex: 2,
		SenderKeyIndex:    2,
		SequenceIndex:     2,
	}
	expected = append(expected, selfPromise)
	assert.Equal(t, len(promises), 3)
	assert.Equal(t, len(promises), len(expected))
	for i := 0; i < len(promises); i++ {
		assert.True(t, expected[i].Equal(promises[i]))
	}

	expected = []*entities.Promise{}
	promises, err = generateExpectedPromises(promise, 8)
	assert.Nil(t, err)

	for i := 7; i > 5; i-- {
		p := &entities.Promise{
			RecipientKeyIndex: 2,
			SenderKeyIndex:    sequence[i],
			SequenceIndex:     uint32(i),
		}
		expected = append(expected, p)
	}
	selfPromise = &entities.Promise{
		RecipientKeyIndex: 2,
		SenderKeyIndex:    2,
		SequenceIndex:     8,
	}
	expected = append(expected, selfPromise)
	assert.Equal(t, len(promises), 3)
	assert.Equal(t, len(promises), len(expected))
	for i := 0; i < len(promises); i++ {
		assert.True(t, expected[i].Equal(promises[i]))
	}
}

func TestGenerationRound(t *testing.T) {
	var expected []*entities.Promise

	roundPromises, err := generationRound(sequence, 2, -1)
	assert.Equal(t, err.Error(), "Index out of range")
	assert.Equal(t, len(roundPromises), 0)

	roundPromises, err = generationRound(sequence, 2, 2)
	for i := 1; i > -1; i-- {
		p := &entities.Promise{
			RecipientKeyIndex: 2,
			SenderKeyIndex:    sequence[i],
			SequenceIndex:     uint32(i),
		}
		expected = append(expected, p)
	}

	assert.Equal(t, err, nil)
	assert.Equal(t, len(roundPromises), 2)
	for i := 0; i < 2; i++ {
		assert.True(t, expected[i].Equal(roundPromises[i]))
	}

	expected = []*entities.Promise{}
	roundPromises, err = generationRound(sequence, 2, 5)
	for i := 4; i > 2; i-- {
		p := &entities.Promise{
			RecipientKeyIndex: 2,
			SenderKeyIndex:    sequence[i],
			SequenceIndex:     uint32(i),
		}
		expected = append(expected, p)
	}

	assert.Equal(t, err, nil)
	assert.Equal(t, len(roundPromises), 2)
	for i := 0; i < 2; i++ {
		assert.True(t, expected[i].Equal(roundPromises[i]))
	}

	expected = []*entities.Promise{}
	roundPromises, err = generationRound(sequence, 2, 8)
	for i := 7; i > 5; i-- {
		p := &entities.Promise{
			RecipientKeyIndex: 2,
			SenderKeyIndex:    sequence[i],
			SequenceIndex:     uint32(i),
		}
		expected = append(expected, p)
	}

	assert.Equal(t, err, nil)
	assert.Equal(t, len(roundPromises), 2)
	for i := 0; i < 2; i++ {
		assert.True(t, expected[i].Equal(roundPromises[i]))
	}

	roundPromises, err = generationRound(sequence, 2, 9)
	assert.Equal(t, err.Error(), "Index out of range")
	assert.Equal(t, len(roundPromises), 0)
}

func TestAddPromiseToExpected(t *testing.T) {
	promise0 := &entities.Promise{
		RecipientKeyIndex: 2,
		SenderKeyIndex:    0,
		SequenceIndex:     0,
	}

	promise1 := &entities.Promise{
		RecipientKeyIndex: 2,
		SenderKeyIndex:    1,
		SequenceIndex:     1,
	}

	promise2 := &entities.Promise{
		RecipientKeyIndex: 2,
		SenderKeyIndex:    1,
		SequenceIndex:     2,
	}

	expected := []*entities.Promise{}

	assert.Equal(t, len(expected), 0)

	expected = addPromiseToExpected(expected, promise0)
	assert.Equal(t, len(expected), 1)
	assert.True(t, promise0.Equal(expected[0]))

	expected = addPromiseToExpected(expected, promise0)
	assert.Equal(t, len(expected), 1)
	assert.True(t, promise0.Equal(expected[0]))

	expected = addPromiseToExpected(expected, promise1)
	assert.Equal(t, len(expected), 2)
	assert.True(t, promise0.Equal(expected[0]))
	assert.True(t, promise1.Equal(expected[1]))

	expected = addPromiseToExpected(expected, promise2)
	assert.Equal(t, len(expected), 2)
	assert.True(t, promise0.Equal(expected[0]))
	assert.True(t, promise2.Equal(expected[1]))
}

func TestContainsPreviousPromise(t *testing.T) {
	promise0 := &entities.Promise{
		RecipientKeyIndex: 2,
		SenderKeyIndex:    0,
		SequenceIndex:     0,
	}

	promise1 := &entities.Promise{
		RecipientKeyIndex: 2,
		SenderKeyIndex:    1,
		SequenceIndex:     1,
	}

	promise2 := &entities.Promise{
		RecipientKeyIndex: 2,
		SenderKeyIndex:    1,
		SequenceIndex:     2,
	}

	promises := []*entities.Promise{}

	assert.Equal(t, containsPreviousPromise(promises, promise1), -1)

	promises = append(promises, promise0)
	assert.Equal(t, containsPreviousPromise(promises, promise1), -1)

	promises = append(promises, promise1)
	assert.Equal(t, containsPreviousPromise(promises, promise1), -1)
	assert.Equal(t, containsPreviousPromise(promises, promise2), 1)

	promises = append(promises, promise2)
	assert.Equal(t, containsPreviousPromise(promises, promise2), 1)
}

func TestContainsPromise(t *testing.T) {
	promise0 := &entities.Promise{
		RecipientKeyIndex: 2,
		SenderKeyIndex:    0,
		SequenceIndex:     0,
	}

	promise1 := &entities.Promise{
		RecipientKeyIndex: 2,
		SenderKeyIndex:    1,
		SequenceIndex:     1,
	}

	promises := []*entities.Promise{promise0}

	assert.False(t, containsPromise(promises, promise1))
	assert.True(t, containsPromise(promises, promise0))
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
		t.Fatal("Contract should have beed generated")
	}
}

// TO MODIFY WHEN SOURCE FUNCTION WILL BE UPDATED
func TestGenerateSignedContract(t *testing.T) {
	// TODO
	id := bson.NewObjectId()
	var promises []entities.Promise
	for i := 1; i > -1; i-- {
		p := entities.Promise{
			RecipientKeyIndex: 2,
			SenderKeyIndex:    sequence[i],
			SequenceIndex:     uint32(i),
		}
		promises = append(promises, p)
	}
	selfPromise := entities.Promise{
		RecipientKeyIndex: 2,
		SenderKeyIndex:    2,
		SequenceIndex:     2,
	}
	promises = append(promises, selfPromise)

	archives := &entities.SignatureArchives{
		ID:               id,
		Signers:          signersEntities,
		ReceivedPromises: promises,
	}

	contract := GenerateSignedContract(archives)

	var pseudoContract string
	for _, p := range promises {
		signature := "SIGNATURE FROM SIGNER " + fmt.Sprintf("%x", signersEntities[p.SenderKeyIndex].Hash)
		signature += " ON SIGNATURE n° " + fmt.Sprint(id) + "\n"
		pseudoContract += signature
	}

	expected := []byte(pseudoContract)
	assert.Equal(t, contract, expected)
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
