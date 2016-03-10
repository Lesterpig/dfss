package contract_test

import (
	"testing"

	"dfss/dfssp/api"
	"dfss/dfssp/contract"
	"dfss/dfssp/entities"
	"github.com/bmizerany/assert"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2/bson"
)

func TestReadySignBadContract(t *testing.T) {
	dropDataset()
	createDataset()

	client := clientTest(t)
	result, err := client.ReadySign(context.Background(), &api.ReadySignRequest{
		ContractUuid: bson.NewObjectId().Hex(),
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, api.ErrorCode_INVARG, result.ErrorCode.Code)
}

func TestReadySignBadUUID(t *testing.T) {
	dropDataset()
	createDataset()

	client := clientTest(t)
	result, err := client.ReadySign(context.Background(), &api.ReadySignRequest{
		ContractUuid: "VERY_VERY_BAD",
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, api.ErrorCode_INVARG, result.ErrorCode.Code)
}

func TestFindAndUpdatePendingSigner(t *testing.T) {
	signers := []entities.Signer{
		entities.Signer{Email: "a"},
		entities.Signer{Email: "b"},
		entities.Signer{Email: "c"},
	}

	signersReady := make([]bool, 3)
	assert.Equal(t, false, contract.FindAndUpdatePendingSigner("a", &signersReady, &signers))
	assert.Equal(t, false, contract.FindAndUpdatePendingSigner("c", &signersReady, &signers))
	assert.Equal(t, false, contract.FindAndUpdatePendingSigner("c", &signersReady, &signers))
	assert.Equal(t, false, contract.FindAndUpdatePendingSigner("a", &signersReady, &signers))
	assert.Equal(t, true, contract.FindAndUpdatePendingSigner("b", &signersReady, &signers))
	assert.Equal(t, true, contract.FindAndUpdatePendingSigner("b", &signersReady, &signers))
	assert.Equal(t, true, contract.FindAndUpdatePendingSigner("a", &signersReady, &signers))
}
