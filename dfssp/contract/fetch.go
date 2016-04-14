package contract

import (
	"dfss/dfssp/api"
	"dfss/dfssp/entities"
	"dfss/mgdb"
	"gopkg.in/mgo.v2/bson"
)

// Fetch returns the protobuf message when asking a specific contract containing a specific user.
func Fetch(db *mgdb.MongoManager, contractUUID string, clientHash []byte) *api.Contract {
	if !bson.IsObjectIdHex(contractUUID) {
		return &api.Contract{
			ErrorCode: &api.ErrorCode{Code: api.ErrorCode_INVARG},
		}
	}

	repository := entities.NewContractRepository(db.Get("contracts"))
	contract, _ := repository.GetWithSigner(clientHash, bson.ObjectIdHex(contractUUID))
	if contract == nil {
		return &api.Contract{
			ErrorCode: &api.ErrorCode{Code: api.ErrorCode_BADAUTH},
		}
	}

	data, err := GetJSON(contract, nil)
	if err != nil {
		return &api.Contract{
			ErrorCode: &api.ErrorCode{Code: api.ErrorCode_INTERR},
		}
	}

	return &api.Contract{
		ErrorCode: &api.ErrorCode{Code: api.ErrorCode_SUCCESS},
		Json:      data,
	}
}
