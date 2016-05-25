// Package contract manages contracts and signatures (creation and execution).
package contract

import (
	"crypto/sha512"
	"log"
	"strings"
	"time"

	"dfss/dfssp/api"
	"dfss/dfssp/entities"
	"dfss/dfssp/templates"
	"dfss/mgdb"
	"gopkg.in/mgo.v2/bson"
)

// Builder contains internal information to create a new contract.
type Builder struct {
	m              *mgdb.MongoManager
	in             *api.PostContractRequest
	signers        []entities.User
	missingSigners []string
	Contract       *entities.Contract
}

// NewContractBuilder creates a new builder from current context.
// Call Execute() on the builder to get a result from it.
func NewContractBuilder(m *mgdb.MongoManager, in *api.PostContractRequest) *Builder {
	return &Builder{
		m:  m,
		in: in,
	}
}

// Execute triggers the creation of a new contract.
func (c *Builder) Execute() *api.ErrorCode {

	inputError := c.checkInput()
	if inputError != nil {
		return inputError
	}

	err := c.fetchSigners()
	if err != nil {
		log.Println(err)
		return &api.ErrorCode{Code: api.ErrorCode_INTERR, Message: "Database error"}
	}

	err = c.addContract()
	if err != nil {
		log.Println(err)
		return &api.ErrorCode{Code: api.ErrorCode_INTERR}
	}

	if len(c.missingSigners) > 0 {
		c.sendPendingContractMail()
		return &api.ErrorCode{Code: api.ErrorCode_WARNING, Message: "Some users are not ready yet"}
	}
	c.SendNewContractMail()
	return &api.ErrorCode{Code: api.ErrorCode_SUCCESS}

}

// checkInput checks that a PostContractRequest is well-formed
func (c *Builder) checkInput() *api.ErrorCode {
	if len(c.in.Signer) == 0 {
		return &api.ErrorCode{Code: api.ErrorCode_INVARG, Message: "Expecting at least one signer"}
	}

	if len(c.in.Filename) == 0 {
		return &api.ErrorCode{Code: api.ErrorCode_INVARG, Message: "Expecting a valid filename"}
	}

	if len(c.in.Hash) != sha512.Size {
		return &api.ErrorCode{Code: api.ErrorCode_INVARG, Message: "Expecting a valid sha512 hash"}
	}

	return nil
}

// fetchSigners fetches authenticated users for this contract from the DB
func (c *Builder) fetchSigners() error {
	var users []entities.User

	// Convert emails to case-tolerant emails
	var conditions []bson.RegEx
	for _, s := range c.in.Signer {
		conditions = append(conditions, bson.RegEx{Pattern: "^" + s + "$", Options: "i"})
	}

	// Fetch users where email is part of the signers slice in request
	// and authentication is valid
	err := c.m.Get("users").FindAll(bson.M{
		"expiration": bson.M{"$gt": time.Now()},
		"email":      bson.M{"$in": conditions},
	}, &users)
	if err != nil {
		return err
	}

	// Locate missing users
	for _, s := range c.in.Signer {
		found := false
		lowerEmail := strings.ToLower(s)
		for _, u := range users {
			if lowerEmail == strings.ToLower(u.Email) {
				found = true
				break
			}
		}
		if !found {
			c.missingSigners = append(c.missingSigners, s) // a list of not valid mail adress
		}
	}

	c.signers = users
	return nil
}

// addContract inserts the contract into the DB
func (c *Builder) addContract() error {
	contract := entities.NewContract()
	for _, s := range c.signers {
		contract.AddSigner(&s.ID, s.Email, s.CertHash)
	}
	for _, s := range c.missingSigners {
		contract.AddSigner(nil, s, nil)
	}

	contract.Comment = c.in.Comment
	contract.Ready = len(c.missingSigners) == 0
	contract.File.Name = c.in.Filename
	contract.File.Hash = c.in.Hash
	contract.File.Hosted = false

	_, err := c.m.Get("contracts").Insert(contract)
	c.Contract = contract

	return err
}

// SendNewContractMail sends a mail to each known signer in a contract containing the DFSS file
func (c *Builder) SendNewContractMail() {
	conn := templates.MailConn()
	if conn == nil {
		return
	}
	defer func() { _ = conn.Close() }()

	rcpts := make([]string, len(c.Contract.Signers))
	for i, s := range c.Contract.Signers {
		rcpts[i] = s.Email
	}

	content, err := templates.Get("contract", c.Contract)
	if err != nil {
		log.Println(err)
		return
	}

	file, err := GetJSON(c.Contract)
	if err != nil {
		log.Println(err)
		return
	}

	fileSmallHash := c.Contract.ID.Hex()

	_ = conn.Send(
		rcpts,
		"[DFSS] You are invited to sign "+c.Contract.File.Name,
		content,
		[]string{"application/json"},
		[]string{fileSmallHash + ".json"},
		[][]byte{file},
	)
}

// sendPendingContractMail sends a mail to non-authenticated signers to invite them
func (c *Builder) sendPendingContractMail() {
	conn := templates.MailConn()
	if conn == nil {
		return
	}
	defer func() { _ = conn.Close() }()

	content, err := templates.Get("invitation", c.Contract)
	if err != nil {
		log.Println(err)
		return
	}

	_ = conn.Send(c.missingSigners, "[DFSS] You are invited to sign "+c.Contract.File.Name, content, nil, nil, nil)
}
