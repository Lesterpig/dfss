package entities

import (
	"dfss/mgdb"

	"time"

	"gopkg.in/mgo.v2/bson"
)

// User : User stored in mongo
type User struct {
	ID           bson.ObjectId `key:"_id" bson:"_id"`                   // Internal id of a User
	Email        string        `key:"email" bson:"email"`               // Email of a User
	Registration time.Time     `key:"registration" bson:"registration"` // Time of registration of the User
	Expiration   time.Time     `key:"expiration" bson:"expiration"`     // Certificate expiration of the User
	RegToken     string        `key:"regToken" bson:"regToken"`         // Token used for registering a User
	Csr          string        `key:"csr" bson:"csr"`                   // Certificate request at PEM format
	Certificate  string        `key:"certificate" bson:"certificate"`   // Certificate of the User
	CertHash     []byte        `key:"certHash" bson:"certHash"`         // Hash of the certificate
}

// NewUser : Create a new User
func NewUser() *User {
	return &User{
		ID:           bson.NewObjectId(),
		Registration: time.Now().UTC(),
	}
}

// UserRepository : Holds all the complex methods regarding a user
type UserRepository struct {
	Collection *mgdb.MongoCollection
}

// NewUserRepository : Creates a new user repository from the given connection
func NewUserRepository(collection *mgdb.MongoCollection) *UserRepository {
	return &UserRepository{
		collection,
	}
}

// FetchByMailAndHash : Fetches a User from its email and certificate hash
func (repository *UserRepository) FetchByMailAndHash(email string, hash []byte) (*User, error) {
	var users []User
	err := repository.Collection.FindAll(bson.M{"email": email, "certHash": hash}, &users)
	if err != nil || len(users) == 0 {
		return nil, err
	}

	users[0].Registration = users[0].Registration.UTC()
	return &users[0], err
}
