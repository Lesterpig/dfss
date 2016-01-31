package contract

import (
	"crypto/sha512"
	"log"
	"time"

	"dfss/dfssp/api"
	"dfss/dfssp/entities"
	"dfss/mgdb"
	"gopkg.in/mgo.v2/bson"
)

// PostRoute is the GRPC route designed to create a contract.
func PostRoute(m *mgdb.MongoManager, in *api.PostContractRequest) *api.ErrorCode {

	inputError := checkInput(in)
	if inputError != nil {
		return inputError
	}

	signers, missingSigners, err := fetchSigners(m, in.Signer)
	if err != nil {
		log.Println(err)
		return &api.ErrorCode{Code: api.ErrorCode_INTERR, Message: "Database error"}
	}

	err = addContract(m, in, signers, missingSigners)
	if err != nil {
		log.Println(err)
		return &api.ErrorCode{Code: api.ErrorCode_INTERR}
	}

	if len(missingSigners) > 0 {
		return &api.ErrorCode{Code: api.ErrorCode_WARNING, Message: "Some users are not ready yet"}
	}
	return &api.ErrorCode{Code: api.ErrorCode_SUCCESS}

}

func checkInput(in *api.PostContractRequest) *api.ErrorCode {

	if len(in.Signer) == 0 {
		return &api.ErrorCode{Code: api.ErrorCode_INVARG, Message: "Expecting at least one signer"}
	}

	if len(in.Filename) == 0 {
		return &api.ErrorCode{Code: api.ErrorCode_INVARG, Message: "Expecting a valid filename"}
	}

	if len(in.Hash) != sha512.Size {
		return &api.ErrorCode{Code: api.ErrorCode_INVARG, Message: "Expecting a valid sha512 hash"}
	}

	return nil

}

func fetchSigners(m *mgdb.MongoManager, signers []string) ([]entities.User, []string, error) {
	var users []entities.User
	err := m.Get("users").FindAll(bson.M{
		"expiration": bson.M{"$gt": time.Now()},
		"email":      bson.M{"$in": signers},
	}, &users)
	if err != nil {
		return nil, nil, err
	}

	// Locate missing users
	var missing []string
	for _, s := range signers {
		found := false
		for _, u := range users {
			if s == u.Email {
				found = true
				break
			}
		}
		if !found {
			missing = append(missing, s)
		}
	}

	return users, missing, nil
}

func addContract(m *mgdb.MongoManager, in *api.PostContractRequest, signers []entities.User, missingSigners []string) error {
	contract := entities.NewContract()
	for _, s := range signers {
		contract.AddSigner(&s.ID, s.Email, s.CertHash)
	}
	for _, s := range missingSigners {
		contract.AddSigner(nil, s, "")
	}

	contract.Comment = in.Comment
	contract.Ready = len(missingSigners) == 0
	contract.File.Name = in.Filename
	contract.File.Hash = in.Hash
	contract.File.Hosted = false

	_, err := m.Get("contracts").Insert(contract)
	return err
}
