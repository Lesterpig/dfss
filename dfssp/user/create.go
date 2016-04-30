// Package user handles user creation.
package user

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"time"

	"dfss/auth"
	"dfss/dfssp/api"
	"dfss/dfssp/authority"
	"dfss/dfssp/contract"
	"dfss/dfssp/entities"
	"dfss/dfssp/templates"
	"dfss/mgdb"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2/bson"
)

// Check if the registration request has usable fields
func checkRegisterRequest(in *api.RegisterRequest) *api.ErrorCode {
	if len(in.Email) == 0 {
		return &api.ErrorCode{Code: api.ErrorCode_INVARG, Message: "Invalid email length"}
	}

	if len(in.Request) == 0 {
		return &api.ErrorCode{Code: api.ErrorCode_INVARG, Message: "Invalid request length"}
	}

	_, err := auth.PEMToCertificateRequest([]byte(in.Request))

	if err != nil {
		return &api.ErrorCode{Code: api.ErrorCode_INVARG, Message: err.Error()}
	}

	return nil
}

// Send the verification email in response to the specified registration request
//
// This method should only be called AFTER checking the RegisterRequest for validity
func sendVerificationMail(in *api.RegisterRequest, token string) {
	conn := templates.MailConn()
	if conn == nil {
		log.Println("Couldn't connect to the dfssp mail server")
		return
	}
	defer func() { _ = conn.Close() }()

	rcpts := []string{in.Email}

	mail := templates.VerificationMail{Token: token}
	content, err := templates.Get("verificationMail", mail)
	if err != nil {
		log.Println(err)
		return
	}

	err = conn.Send(
		rcpts,
		"[DFSS] Registration email validation",
		content,
		nil,
		nil,
		nil,
	)
	if err != nil {
		log.Println(err)
		return
	}
}

// Register checks if the registration request is valid, and if so,
// creates the user entry in the database
//
// If there is already an entry in the database with the same email,
// evaluates the request as invalid
//
// The user's ConnectionInfo field is NOT handled here
// This data should be gathered upon beginning the signing sequence
func Register(manager *mgdb.MongoManager, in *api.RegisterRequest) (*api.ErrorCode, error) {
	// Check the request validity
	errCode := checkRegisterRequest(in)
	if errCode != nil {
		return errCode, nil
	}

	// Generating the random token
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		return &api.ErrorCode{Code: api.ErrorCode_INTERR, Message: "Error during the generation of the token"}, nil
	}
	token := fmt.Sprintf("%x", b)

	// If there is already an entry with the same mail, do nothing
	var res []entities.User
	err = manager.Get("users").FindAll(bson.M{
		"email": bson.M{"$eq": in.Email},
	}, &res)
	if len(res) != 0 {
		return &api.ErrorCode{Code: api.ErrorCode_INVARG, Message: "An entry already exists with the same mail"}, nil
	}

	// Creating the new user
	user := entities.NewUser()
	user.Email = in.Email
	user.RegToken = token
	user.Csr = in.Request

	// Adding the new user in the database
	ok, err := manager.Get("users").Insert(*user)
	if !ok {
		return &api.ErrorCode{Code: api.ErrorCode_INTERR, Message: "Error during the insertion of the new user"}, err
	}

	// Sending the email
	sendVerificationMail(in, token)

	return &api.ErrorCode{Code: api.ErrorCode_SUCCESS, Message: "Registration successful ; email sent"}, nil
}

// Check if the authentication request has usable fields
func checkAuthRequest(in *api.AuthRequest) error {
	if len(in.Email) == 0 {
		return errors.New("Invalid email length")
	}

	if len(in.Token) == 0 {
		return errors.New("Invalid token length")
	}

	if viper.GetInt("validity") < 1 {
		return errors.New("Invalid validity duration")
	}

	return nil
}

// Check if the authentication request was made in time
func checkTokenTimeout(user *entities.User) error {
	now := time.Now().UTC()
	bad := now.After(user.Registration.Add(time.Hour * 24))
	if bad {
		return errors.New("Registration request is over 24 hours old")
	}

	return nil
}

// Gerenate the user's certificate and certificate hash according to the specified parameters
//
// This function should only be called AFTER checking the AuthRequest for validity
func generateUserCert(csr string, parent *x509.Certificate, key *rsa.PrivateKey) ([]byte, []byte, error) {
	x509csr, err := auth.PEMToCertificateRequest([]byte(csr))
	if err != nil {
		return nil, nil, err
	}

	cert, err := auth.GetCertificate(viper.GetInt("validity"), auth.GenerateUID(), x509csr, parent, key)
	if err != nil {
		return nil, nil, err
	}

	c, _ := auth.PEMToCertificate(cert)
	certHash := auth.GetCertificateHash(c)

	return cert, certHash, nil
}

// Auth checks if the authentication request is valid, and if so,
// generate the certificate and certificate hash for the user, and
// updates the user's entry in the database
//
// If there is already an entry in the database with the same email,
// and that this entry already has a certificate and certificate hash,
// evaluates the request as invalid
//
// The user's ConnectionInfo field is NOT handled here
// This data should be gathered upon beginning the signing sequence
func Auth(pid *authority.PlatformID, manager *mgdb.MongoManager, in *api.AuthRequest) (*api.RegisteredUser, error) {
	// Check the request validity
	err := checkAuthRequest(in)
	if err != nil {
		return nil, err
	}

	// Find the user in the database
	var user entities.User

	err = manager.Get("users").Collection.Find(bson.M{
		"email": bson.M{"$eq": in.Email},
	}).One(&user)
	if err != nil {
		return nil, err
	}

	// If the user already has a certificate and certificate hash in the database, does nothing
	if user.Certificate != "" || len(user.CertHash) != 0 {
		return nil, errors.New("User is already registered")
	}

	// Check if the delta between now and the moment the user was created (ie the moment he sent the register request) is in bound of 24h
	err = checkTokenTimeout(&user)
	if err != nil {
		return nil, err
	}

	// Check if the token is correct
	if in.Token != user.RegToken {
		return nil, errors.New("Token mismatch")
	}

	// Generate the certificates and hash
	cert, certHash, err := generateUserCert(user.Csr, pid.RootCA, pid.Pkey)
	if err != nil {
		return nil, err
	}

	user.Certificate = string(cert)
	user.CertHash = certHash
	user.Expiration = time.Now().AddDate(0, 0, viper.GetInt("validity"))

	// Updating the database
	ok, err := manager.Get("users").UpdateByID(user)
	if !ok {
		return nil, err
	}

	// Update missed contracts in background
	go launchMissedContracts(manager, &user)

	// Returning the RegisteredUser message
	return &api.RegisteredUser{ClientCert: user.Certificate}, nil
}

func launchMissedContracts(manager *mgdb.MongoManager, user *entities.User) {

	repository := entities.NewContractRepository(manager.Get("contracts"))
	contracts, err := repository.GetWaitingForUser(user.Email)
	if err != nil {
		log.Println("Cannot get missed contracts for user", user.Email+":", err)
	}

	for _, c := range contracts {

		c.Ready = true
		for i := range c.Signers {
			if c.Signers[i].Email == user.Email {
				c.Signers[i].Hash = user.CertHash
				c.Signers[i].UserID = user.ID
			}
			if len(c.Signers[i].Hash) == 0 {
				c.Ready = false
			}
		}

		// Update contract in database
		_, err = repository.Collection.UpdateByID(c)
		if err != nil {
			log.Println("Cannot update missed contract", c.ID, "for user", user.Email+":", err)
		}

		if c.Ready {
			// Send required mails
			builder := contract.NewContractBuilder(manager, nil)
			builder.Contract = &c
			builder.SendNewContractMail()
		}
	}

}
