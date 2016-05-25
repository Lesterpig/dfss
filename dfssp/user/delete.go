package user

import (
	api "dfss/dfssp/api"
	"dfss/mgdb"
	"gopkg.in/mgo.v2/bson"
)

// Unregister delete a user based on the provided certificate hash
func Unregister(manager *mgdb.MongoManager, userCertificateHash []byte) *api.ErrorCode {
	count, err := manager.Get("users").DeleteAll(bson.M{
		"certHash": bson.M{"$eq": userCertificateHash},
	})
	if err != nil || count == 0 {
		return &api.ErrorCode{Code: api.ErrorCode_INVARG, Message: "No user matching provided certificate"}
	}

	return &api.ErrorCode{Code: api.ErrorCode_SUCCESS}
}
