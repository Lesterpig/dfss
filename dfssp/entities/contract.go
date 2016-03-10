package entities

import (
	"dfss/mgdb"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// File : Represents a file structure
type File struct {
	Name   string `key:"name" bson:"name"`     // Name of the File
	Hash   []byte `key:"hash" bson:"hash"`     // Hash of the File
	Hosted bool   `key:"hosted" bson:"hosted"` // True if hosted on the platform, else false
}

// Signer : Informations about the signer of a contract
type Signer struct {
	UserID bson.ObjectId `key:"userId" bson:"userId"`
	Email  string        `key:"email" bson:"email"`
	Hash   []byte        `key:"hash" bson:"hash"`
}

// Contract : Informations about a contract to be signed
type Contract struct {
	ID      bson.ObjectId `key:"_id" bson:"_id"`
	Date    time.Time     `key:"date" bson:"date"`
	Comment string        `key:"comment" bson:"comment"`
	Ready   bool          `key:"ready" bson:"ready"`
	File    *File         `key:"file" bson:"file"`
	Signers []Signer      `key:"signers" bson:"signers"`
}

// NewContract : Creates a new contract
func NewContract() *Contract {
	file := File{}
	var signers []Signer
	return &Contract{
		ID:      bson.NewObjectId(),
		Date:    time.Now(),
		File:    &file,
		Signers: signers,
	}
}

// AddSigner : Add a signer to the contract
func (c *Contract) AddSigner(id *bson.ObjectId, email string, hash []byte) {
	signer := &Signer{}
	signer.Email = email

	if id != nil {
		signer.UserID = *id
	} else {
		signer.UserID = bson.ObjectIdHex("000000000000000000000000")
	}

	signer.Hash = hash
	c.Signers = append(c.Signers, *signer)
}

// GetHashChain returns the ordered slice of signers hashes.
// It's used to check the dfss file if needed.
func (c *Contract) GetHashChain() [][]byte {
	chain := make([][]byte, len(c.Signers))
	for i, s := range c.Signers {
		chain[i] = s.Hash
	}
	return chain
}

// ContractRepository to contains every complex methods related to contract
type ContractRepository struct {
	Collection *mgdb.MongoCollection
}

// NewContractRepository : Creates a new Contract Repository
func NewContractRepository(collection *mgdb.MongoCollection) *ContractRepository {
	return &ContractRepository{
		collection,
	}
}

// GetWaitingForUser returns contracts waiting a specific unauthenticated user to start
func (r *ContractRepository) GetWaitingForUser(email string) ([]Contract, error) {
	var res []Contract
	err := r.Collection.FindAll(bson.M{
		"ready": false,
		"signers": bson.M{
			"$elemMatch": bson.M{
				"email": email,
				"hash":  []byte{},
			}},
	}, &res)
	return res, err
}

// CheckAuthorization checks that a client is allowed to sign a specific contract
func (r *ContractRepository) CheckAuthorization(signerHash []byte, contractID bson.ObjectId) bool {

	count, _ := r.Collection.Collection.Find(bson.M{
		"_id":   contractID,
		"ready": true,
		"signers": bson.M{
			"$elemMatch": bson.M{"hash": signerHash},
		},
	}).Count()

	return count == 1
}
