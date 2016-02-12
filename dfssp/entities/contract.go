package entities

import (
	"dfss/mgdb"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// File : Represents a file structure
type File struct {
	Name   string `key:"name" bson:"name"`     // Name of the File
	Hash   string `key:"hash" bson:"hash"`     // Hash of the File
	Hosted bool   `key:"hosted" bson:"hosted"` // True if hosted on the platform, else false
}

// Signer : Informations about the signer of a contract
type Signer struct {
	UserID bson.ObjectId `key:"userId" bson:"userId"`
	Email  string        `key:"email" bson:"email"`
	Hash   string        `key:"hash" bson:"hash"`
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
func (c *Contract) AddSigner(id *bson.ObjectId, email, hash string) {
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
				"email" : email,
				"hash" : "",
			}},
	}, &res)
	return res, err
}
