package contract

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

// NewSigner : Creates a signer associated to a contract
func NewSigner() *Signer {
	return &Signer{
		UserID: bson.NewObjectId(),
	}
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
func (c *Contract) AddSigner(email, hash string) {
	signer := NewSigner()
	signer.Email = email
	signer.Hash = hash
	c.Signers = append(c.Signers, *signer)
}

// Repository to contains every complex methods related to contract
type Repository struct {
	Collection *mgdb.MongoCollection
}

// NewRepository : Creates a new Contract Repository
func NewRepository(collection *mgdb.MongoCollection) *Repository {
	return &Repository{
		collection,
	}
}
