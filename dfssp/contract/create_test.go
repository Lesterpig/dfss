package contract_test

import (
	"crypto/sha512"
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"

	"dfss/auth"
	"dfss/dfssp/api"
	"dfss/dfssp/entities"
	"dfss/net"
	"github.com/bmizerany/assert"
	"golang.org/x/net/context"
)

var user1, user2, user3 *entities.User
var defaultHash = sha512.Sum512([]byte{0})
var defaultHashStr = string(defaultHash[:sha512.Size])

func createDataset() {

	user1 = entities.NewUser() // Regular user
	user2 = entities.NewUser() // Regular user
	user3 = entities.NewUser() // Non-auth user

	user1.Email = "user1@example.com"
	user1.Expiration = time.Now().AddDate(1, 0, 0)
	user1.Certificate = "Certificate1"
	user1.CertHash = "Hash1"

	user2.Email = "user2@example.com"
	user2.Expiration = time.Now().AddDate(1, 0, 0)
	user2.Certificate = "Certificate2"
	user2.CertHash = "Hash2"

	user3.Email = "user3@example.com"
	user3.Expiration = time.Now().AddDate(0, 0, -1)
	user3.Certificate = "Certificate3"
	user3.CertHash = "Hash3"

	_, _ = manager.Get("users").Insert(user1)
	_, _ = manager.Get("users").Insert(user2)
	_, _ = manager.Get("users").Insert(user3)

}

func dropDataset() {

	_ = manager.Get("users").Drop()
	_ = manager.Get("contracts").Drop()

}

func clientTest(t *testing.T) api.PlatformClient {
	// TODO if anyone needs this function in another test suite, please put it in a separate file
	// to avoid code duplication
	caData, _ := ioutil.ReadFile(filepath.Join("..", "testdata", "dfssp_rootCA.pem"))
	certData, _ := ioutil.ReadFile(filepath.Join("..", "..", "dfssc", "testdata", "cert.pem"))
	keyData, _ := ioutil.ReadFile(filepath.Join("..", "..", "dfssc", "testdata", "key.pem"))
	ca, _ := auth.PEMToCertificate(caData)
	cert, _ := auth.PEMToCertificate(certData)
	key, _ := auth.EncryptedPEMToPrivateKey(keyData, "password")

	conn, err := net.Connect("localhost:9090", cert, key, ca)
	if err != nil {
		t.Fatal("Unable to connect:", err)
	}

	return api.NewPlatformClient(conn)
}

func TestAddContractBadAuth(t *testing.T) {
	caData, _ := ioutil.ReadFile(filepath.Join("..", "testdata", "dfssp_rootCA.pem"))
	ca, _ := auth.PEMToCertificate(caData)
	conn, err := net.Connect("localhost:9090", nil, nil, ca)
	if err != nil {
		t.Fatal("Unable to connect:", err)
	}
	client := api.NewPlatformClient(conn)
	errorCode, err := client.PostContract(context.Background(), &api.PostContractRequest{})
	assert.Equal(t, nil, err)
	assert.Equal(t, api.ErrorCode_BADAUTH, errorCode.Code)
}

func TestAddContract(t *testing.T) {
	dropDataset()
	createDataset()

	client := clientTest(t)
	errorCode, err := client.PostContract(context.Background(), &api.PostContractRequest{
		Hash:     defaultHashStr,
		Filename: "ContractFilename",
		Signer:   []string{user1.Email, user2.Email},
		Comment:  "ContractComment",
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, api.ErrorCode_SUCCESS, errorCode.Code)

	// Check database content
	var contracts []entities.Contract
	err = manager.Get("contracts").FindAll(nil, &contracts)
	if err != nil {
		t.Fatal("Unexpected db error:", err)
	}

	assert.Equal(t, 1, len(contracts))
	assert.Equal(t, defaultHashStr, contracts[0].File.Hash)
	assert.Equal(t, "ContractFilename", contracts[0].File.Name)
	assert.Equal(t, "ContractComment", contracts[0].Comment)
	assert.T(t, contracts[0].Ready)

	assert.Equal(t, 2, len(contracts[0].Signers))
	assert.Equal(t, user1.ID, contracts[0].Signers[0].UserID)
	assert.Equal(t, user1.CertHash, contracts[0].Signers[0].Hash)
	assert.Equal(t, user1.Email, contracts[0].Signers[0].Email)
	assert.Equal(t, user2.ID, contracts[0].Signers[1].UserID)
	assert.Equal(t, user2.CertHash, contracts[0].Signers[1].Hash)
	assert.Equal(t, user2.Email, contracts[0].Signers[1].Email)
}

func TestAddContractMissingUser(t *testing.T) {
	dropDataset()
	createDataset()

	client := clientTest(t)
	errorCode, err := client.PostContract(context.Background(), &api.PostContractRequest{
		Hash:     defaultHashStr,
		Filename: "ContractFilename",
		Signer:   []string{user1.Email, user3.Email},
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, api.ErrorCode_WARNING, errorCode.Code)

	// Check database content
	var contracts []entities.Contract
	err = manager.Get("contracts").FindAll(nil, &contracts)
	if err != nil {
		t.Fatal("Unexpected db error:", err)
	}

	assert.Equal(t, 1, len(contracts))
	assert.Equal(t, defaultHashStr, contracts[0].File.Hash)
	assert.Equal(t, "ContractFilename", contracts[0].File.Name)
	assert.Equal(t, "", contracts[0].Comment)
	assert.T(t, !contracts[0].Ready)

	assert.Equal(t, 2, len(contracts[0].Signers))
	assert.Equal(t, user1.ID, contracts[0].Signers[0].UserID)
	assert.Equal(t, user1.CertHash, contracts[0].Signers[0].Hash)
	assert.Equal(t, user1.Email, contracts[0].Signers[0].Email)
	assert.Equal(t, "000000000000000000000000", contracts[0].Signers[1].UserID.Hex())
	assert.Equal(t, "", contracts[0].Signers[1].Hash)
	assert.Equal(t, user3.Email, contracts[0].Signers[1].Email)
}

func TestAddContractNoUser(t *testing.T) {
	dropDataset()
	createDataset()

	client := clientTest(t)
	errorCode, err := client.PostContract(context.Background(), &api.PostContractRequest{
		Hash:     defaultHashStr,
		Filename: "ContractFilename",
		Signer:   []string{},
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, api.ErrorCode_INVARG, errorCode.Code)

	// Check database content
	var contracts []entities.Contract
	err = manager.Get("contracts").FindAll(nil, &contracts)
	if err != nil {
		t.Fatal("Unexpected db error:", err)
	}

	assert.Equal(t, 0, len(contracts))
}

func TestAddContractDuplicatedUser(t *testing.T) {
	dropDataset()
	createDataset()

	client := clientTest(t)
	errorCode, err := client.PostContract(context.Background(), &api.PostContractRequest{
		Hash:     defaultHashStr,
		Filename: "ContractFilename",
		Signer:   []string{user1.Email, user1.Email, user2.Email},
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, api.ErrorCode_SUCCESS, errorCode.Code)

	// Check database content
	var contracts []entities.Contract
	err = manager.Get("contracts").FindAll(nil, &contracts)
	if err != nil {
		t.Fatal("Unexpected db error:", err)
	}

	assert.Equal(t, 1, len(contracts))
	assert.Equal(t, 2, len(contracts[0].Signers))
}

func TestAddContractNoFilename(t *testing.T) {
	dropDataset()
	createDataset()

	client := clientTest(t)
	errorCode, err := client.PostContract(context.Background(), &api.PostContractRequest{
		Hash:   defaultHashStr,
		Signer: []string{user1.Email},
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, api.ErrorCode_INVARG, errorCode.Code)

	// Check database content
	var contracts []entities.Contract
	err = manager.Get("contracts").FindAll(nil, &contracts)
	if err != nil {
		t.Fatal("Unexpected db error:", err)
	}

	assert.Equal(t, 0, len(contracts))
}

func TestAddContractBadHash(t *testing.T) {
	dropDataset()
	createDataset()

	client := clientTest(t)
	errorCode, err := client.PostContract(context.Background(), &api.PostContractRequest{
		Hash:     "aVeryBadHash",
		Filename: "ContractFilename",
		Signer:   []string{user1.Email},
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, api.ErrorCode_INVARG, errorCode.Code)

	// Check database content
	var contracts []entities.Contract
	err = manager.Get("contracts").FindAll(nil, &contracts)
	if err != nil {
		t.Fatal("Unexpected db error:", err)
	}

	assert.Equal(t, 0, len(contracts))
}
