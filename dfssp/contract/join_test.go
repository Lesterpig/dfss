package contract_test

import (
	"testing"

	"dfss/dfssp/api"
	"dfss/dfssp/entities"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2/bson"
)

func addTestContract() bson.ObjectId {
	contract := entities.NewContract()
	contract.AddSigner(&user1.ID, user1.Email, user1.CertHash)
	contract.Ready = true
	_, _ = manager.Get("contracts").Insert(contract)
	return contract.ID
}

func TestJoinSignature(t *testing.T) {
	dropDataset()
	createDataset()
	contractID := addTestContract()

	client := clientTest(t)
	stream, err := client.JoinSignature(context.Background(), &api.JoinSignatureRequest{
		ContractUuid: contractID.Hex(),
		Port:         5050,
	})
	assert.Equal(t, nil, err)

	user, err := stream.Recv()
	assert.Equal(t, nil, err)
	assert.Equal(t, "", user.ErrorCode.Message)
	assert.Equal(t, api.ErrorCode_SUCCESS, user.ErrorCode.Code)
	assert.Equal(t, contractID.Hex(), user.ContractUuid)
	assert.Equal(t, "test@test.com", user.User.Email)
	assert.Equal(t, uint32(5050), user.User.Port)
}

func TestJoinSignatureBadContract(t *testing.T) {
	dropDataset()
	createDataset()

	client := clientTest(t)
	stream, err := client.JoinSignature(context.Background(), &api.JoinSignatureRequest{
		ContractUuid: bson.NewObjectId().Hex(),
		Port:         5050,
	})
	assert.Equal(t, nil, err)

	user, err := stream.Recv()
	assert.Equal(t, nil, err)
	assert.Equal(t, "unauthorized signature", user.ErrorCode.Message)
	assert.Equal(t, api.ErrorCode_INVARG, user.ErrorCode.Code)
}

func TestJoinSignatureBadUUID(t *testing.T) {
	dropDataset()
	createDataset()

	client := clientTest(t)
	stream, err := client.JoinSignature(context.Background(), &api.JoinSignatureRequest{
		ContractUuid: "VERY_BAD",
		Port:         5050,
	})
	assert.Equal(t, nil, err)

	user, err := stream.Recv()
	assert.Equal(t, nil, err)
	assert.Equal(t, "invalid contract uuid", user.ErrorCode.Message)
	assert.Equal(t, api.ErrorCode_INVARG, user.ErrorCode.Code)
}
