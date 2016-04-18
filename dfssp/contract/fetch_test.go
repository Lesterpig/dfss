package contract_test

import (
	"encoding/json"
	"testing"

	"dfss/dfssp/api"
	"dfss/dfssp/contract"
	"dfss/dfssp/entities"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2/bson"
)

var contract1 *entities.Contract

func insertTestContract(insertUser1 bool) {
	contract1 = entities.NewContract()
	if insertUser1 {
		contract1.AddSigner(&user1.ID, user1.Email, user1.CertHash)
	}
	contract1.AddSigner(&user2.ID, user2.Email, user2.CertHash)
	_, _ = manager.Get("contracts").Insert(contract1)
}

func TestGetContract(t *testing.T) {
	dropDataset()
	createDataset()
	insertTestContract(true)

	client := clientTest(t)
	c, err := client.GetContract(context.Background(), &api.GetContractRequest{
		Uuid: contract1.ID.Hex(),
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, api.ErrorCode_SUCCESS, c.ErrorCode.Code)
	assert.True(t, len(c.Json) > 0)

	// Trying to unmarshal contract
	var entity contract.JSON
	err = json.Unmarshal(c.Json, &entity)
	assert.Equal(t, nil, err)
	assert.Equal(t, 2, len(entity.Signers))
}

func TestGetContractWrongSigner(t *testing.T) {
	dropDataset()
	createDataset()
	insertTestContract(false)

	client := clientTest(t)
	c, err := client.GetContract(context.Background(), &api.GetContractRequest{
		Uuid: contract1.ID.Hex(),
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, api.ErrorCode_BADAUTH, c.ErrorCode.Code)
	assert.Equal(t, 0, len(c.Json))
}

func TestGetContractWrongContract(t *testing.T) {
	dropDataset()
	createDataset()
	insertTestContract(true)

	client := clientTest(t)
	c, err := client.GetContract(context.Background(), &api.GetContractRequest{
		Uuid: bson.NewObjectId().Hex(),
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, api.ErrorCode_BADAUTH, c.ErrorCode.Code)
	assert.Equal(t, 0, len(c.Json))
}

func TestGetContractWrongContractUUID(t *testing.T) {
	dropDataset()
	createDataset()
	insertTestContract(true)

	client := clientTest(t)
	c, err := client.GetContract(context.Background(), &api.GetContractRequest{
		Uuid: "wrongUUID",
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, api.ErrorCode_INVARG, c.ErrorCode.Code)
	assert.Equal(t, 0, len(c.Json))
}
