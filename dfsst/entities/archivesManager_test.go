package entities

import (
	"crypto/sha512"
	"fmt"
	"os"
	"testing"

	cAPI "dfss/dfssc/api"
	"dfss/mgdb"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

var (
	db         string
	collection *mgdb.MongoCollection
	dbManager  *mgdb.MongoManager

	sequence             []uint32
	signers              [][]byte
	contractDocumentHash []byte
	signatureUUID        string
	signatureUUIDBson    bson.ObjectId

	seal []byte

	signersEntities []Signer

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
	signatureUUID = signatureUUIDBson.Hex()

	seal = []byte{}

	signersEntities = make([]Signer, 0)
	for _, s := range signers {
		signerEntity := NewSigner(s)
		signersEntities = append(signersEntities, *signerEntity)
	}
}

func TestMain(m *testing.M) {
	dbManager, err = mgdb.NewManager(db)
	if err != nil {
		fmt.Fprintln(os.Stderr, "An error occured during the connection to MongoDB:", err)
		os.Exit(2)
	}

	collection = dbManager.Get("signatures")
	err = collection.Drop()
	code := m.Run()
	err = collection.Drop()

	if err != nil {
		fmt.Println("An error occurred while droping the collection")
	}

	dbManager.Close()

	os.Exit(code)
}

func TestInitializeArchives(t *testing.T) {
	promise := &cAPI.Promise{
		Context: &cAPI.Context{
			ContractDocumentHash: contractDocumentHash,
			Sequence:             sequence,
			Signers:              signers,
			SignatureUUID:        signatureUUID,
			Seal:                 seal,
		},
	}
	manager := &ArchivesManager{
		DB: dbManager,
	}
	arch := NewSignatureArchives(signatureUUIDBson, sequence, signersEntities, contractDocumentHash, seal)

	err = manager.InitializeArchives(promise, signatureUUIDBson, &signersEntities)
	assert.Nil(t, err)
	arch.Signers = manager.Archives.Signers
	assert.Equal(t, manager.Archives, arch)

	ok, err := collection.DeleteByID(*manager.Archives)
	assert.Equal(t, err, nil)
	assert.Equal(t, ok, true)
}

func TestContainsSignature(t *testing.T) {
	archives := NewSignatureArchives(signatureUUIDBson, sequence, signersEntities, contractDocumentHash, seal)
	manager := &ArchivesManager{
		DB:       dbManager,
		Archives: archives,
	}

	b, arch := manager.ContainsSignature(signatureUUIDBson)
	assert.Equal(t, b, false)
	assert.Equal(t, arch, &SignatureArchives{})

	ok, err := collection.Insert(archives)
	assert.Equal(t, ok, true)
	assert.Equal(t, err, nil)

	b, arch = manager.ContainsSignature(signatureUUIDBson)
	assert.Equal(t, b, true)
	assert.Equal(t, arch, archives)

	ok, err = collection.DeleteByID(*archives)
	assert.Equal(t, err, nil)
	assert.Equal(t, ok, true)
}

func TestHasReceivedAbortToken(t *testing.T) {
	archives := NewSignatureArchives(signatureUUIDBson, sequence, signersEntities, contractDocumentHash, seal)
	manager := &ArchivesManager{
		DB:       dbManager,
		Archives: archives,
	}
	signerIndex := uint32(1)
	abortedSigner0 := NewAbortedSigner(uint32(0), uint32(1))
	abortedSigner1 := NewAbortedSigner(signerIndex, uint32(1))

	assert.Equal(t, len(archives.AbortedSigners), 0)

	aborted := manager.HasReceivedAbortToken(signerIndex)
	assert.Equal(t, aborted, false)

	archives.AbortedSigners = append(archives.AbortedSigners, *abortedSigner0)

	aborted = manager.HasReceivedAbortToken(signerIndex)
	assert.Equal(t, aborted, false)

	archives.AbortedSigners = append(archives.AbortedSigners, *abortedSigner1)

	aborted = manager.HasReceivedAbortToken(signerIndex)
	assert.Equal(t, aborted, true)
}

func TestWasContractSigned(t *testing.T) {
	archives := NewSignatureArchives(signatureUUIDBson, sequence, signersEntities, contractDocumentHash, seal)
	manager := &ArchivesManager{
		DB:       dbManager,
		Archives: archives,
	}

	signed, contract := manager.WasContractSigned()
	assert.Equal(t, signed, false)
	assert.Equal(t, len(contract), 0)

	c := []byte{0}
	archives.SignedContract = c
	signed, contract = manager.WasContractSigned()
	assert.Equal(t, signed, true)
	assert.Equal(t, contract, c)
}

func TestHasSignerPromised(t *testing.T) {
	archives := NewSignatureArchives(signatureUUIDBson, sequence, signersEntities, contractDocumentHash, seal)
	manager := &ArchivesManager{
		DB:       dbManager,
		Archives: archives,
	}
	ok := manager.HasSignerPromised(1)
	assert.Equal(t, len(archives.ReceivedPromises), 0)
	assert.Equal(t, ok, false)

	promise0 := &Promise{
		RecipientKeyIndex: 1,
		SenderKeyIndex:    0,
	}
	archives.ReceivedPromises = append(archives.ReceivedPromises, *promise0)
	assert.Equal(t, len(archives.ReceivedPromises), 1)

	ok = manager.HasSignerPromised(1)
	assert.Equal(t, ok, false)

	promise1 := &Promise{
		RecipientKeyIndex: 1,
		SenderKeyIndex:    1,
	}
	archives.ReceivedPromises = append(archives.ReceivedPromises, *promise1)
	assert.Equal(t, len(archives.ReceivedPromises), 2)

	ok = manager.HasSignerPromised(1)
	assert.Equal(t, ok, true)

	promise2 := &Promise{
		RecipientKeyIndex: 0,
		SenderKeyIndex:    1,
	}
	archives.ReceivedPromises = append(archives.ReceivedPromises, *promise2)
	assert.Equal(t, len(archives.ReceivedPromises), 3)

	ok = manager.HasSignerPromised(1)
	assert.Equal(t, ok, true)
}

// TO MODIFY WHEN SOURCE FUNCTION WILL BE UPDATED
func TestAddToAbort(t *testing.T) {
	// TODO
	// Test the abortedIndex field, when promises will be implemented
	archives := NewSignatureArchives(signatureUUIDBson, sequence, signersEntities, contractDocumentHash, seal)
	manager := &ArchivesManager{
		DB:       dbManager,
		Archives: archives,
	}
	promise := &cAPI.Promise{
		Context: &cAPI.Context{
			Signers:  [][]byte{},
			Sequence: sequence,
		},
	}

	assert.Equal(t, len(archives.AbortedSigners), 0)

	sIndex, err := GetIndexOfSigner(promise, signers[1])
	assert.Equal(t, err.Error(), "Signer's hash couldn't be matched")

	promise.Context.Signers = signers
	sIndex, err = GetIndexOfSigner(promise, signers[1])
	assert.Equal(t, err, nil)
	assert.Equal(t, sIndex, uint32(1))

	manager.AddToAbort(sIndex)
	assert.Equal(t, len(archives.AbortedSigners), 1)
	assert.Equal(t, archives.AbortedSigners[0].SignerIndex, uint32(1))

	manager.AddToAbort(sIndex)
	assert.Equal(t, len(archives.AbortedSigners), 1)
	assert.Equal(t, archives.AbortedSigners[0].SignerIndex, uint32(1))
}

func TestAddToDishonest(t *testing.T) {
	archives := NewSignatureArchives(signatureUUIDBson, sequence, signersEntities, contractDocumentHash, seal)
	manager := &ArchivesManager{
		DB:       dbManager,
		Archives: archives,
	}
	promise := &cAPI.Promise{
		Context: &cAPI.Context{
			Signers:  [][]byte{},
			Sequence: sequence,
		},
	}

	assert.Equal(t, len(archives.DishonestSigners), 0)

	sIndex, err := GetIndexOfSigner(promise, signers[1])
	assert.Equal(t, err.Error(), "Signer's hash couldn't be matched")

	promise.Context.Signers = signers
	sIndex, err = GetIndexOfSigner(promise, signers[1])
	assert.Equal(t, err, nil)
	assert.Equal(t, sIndex, uint32(1))

	manager.AddToDishonest(sIndex)
	assert.Equal(t, len(archives.DishonestSigners), 1)
	assert.Equal(t, archives.DishonestSigners[0], uint32(1))

	manager.AddToDishonest(sIndex)
	assert.Equal(t, len(archives.DishonestSigners), 1)
	assert.Equal(t, archives.DishonestSigners[0], uint32(1))
}

func TestAddPromise(t *testing.T) {
	archives := NewSignatureArchives(signatureUUIDBson, sequence, signersEntities, contractDocumentHash, seal)
	manager := &ArchivesManager{
		DB:       dbManager,
		Archives: archives,
	}
	assert.Equal(t, len(archives.ReceivedPromises), 0)

	promise0 := &Promise{
		RecipientKeyIndex: 1,
		SenderKeyIndex:    0,
		SequenceIndex:     0,
	}
	promise1 := &Promise{
		RecipientKeyIndex: 0,
		SenderKeyIndex:    1,
		SequenceIndex:     1,
	}

	manager.AddPromise(promise0)
	assert.Equal(t, len(archives.ReceivedPromises), 1)
	assert.Equal(t, archives.ReceivedPromises[0].RecipientKeyIndex, uint32(1))
	assert.Equal(t, archives.ReceivedPromises[0].SenderKeyIndex, uint32(0))
	assert.Equal(t, archives.ReceivedPromises[0].SequenceIndex, uint32(0))

	manager.AddPromise(promise1)
	assert.Equal(t, len(archives.ReceivedPromises), 2)
	assert.Equal(t, archives.ReceivedPromises[1].RecipientKeyIndex, uint32(0))
	assert.Equal(t, archives.ReceivedPromises[1].SenderKeyIndex, uint32(1))
	assert.Equal(t, archives.ReceivedPromises[1].SequenceIndex, uint32(1))

	manager.AddPromise(promise0)
	assert.Equal(t, len(archives.ReceivedPromises), 2)
}
