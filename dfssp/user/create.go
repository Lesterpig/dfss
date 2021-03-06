// Package user handles user creation.
package user

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"regexp"
	"time"

	"dfss/auth"
	"dfss/dfssp/api"
	"dfss/dfssp/entities"
	"dfss/dfssp/templates"
	"dfss/mgdb"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2/bson"
)

var (
	mailRegex            = regexp.MustCompile(`.+@.+\..+`)
	maxRegistrationDelay = 24 * time.Hour
)

// Check if the registration request has usable fields
func checkRegisterRequest(in *api.RegisterRequest) *api.ErrorCode {
	if len(in.Email) == 0 {
		return &api.ErrorCode{Code: api.ErrorCode_INVARG, Message: "Invalid email length"}
	}

	if b := mailRegex.MatchString(in.Email); !b {
		return &api.ErrorCode{Code: api.ErrorCode_INVARG, Message: "Invalid mail"}
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

	// If there is already an entry with the same mail (case-insensitive), do nothing.
	var res []entities.User
	err = manager.Get("users").FindAll(bson.M{
		"$or": []bson.M{
			bson.M{"expiration": bson.M{"$gt": time.Now()}},                                  // authentified
			bson.M{"registration": bson.M{"$gt": time.Now().Add(-1 * maxRegistrationDelay)}}, // authentifying
		},
		"email": bson.M{"$regex": bson.RegEx{Pattern: "^" + in.Email + "$", Options: "i"}},
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
