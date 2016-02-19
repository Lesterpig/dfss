package tests

import (
	"dfss/dfssp/entities"
	"dfss/mgdb"
	"gopkg.in/mgo.v2/bson"
)

var dbURI string
var dbManager *mgdb.MongoManager

func eraseDatabase() {
	_ = dbManager.Database.DropDatabase()
}

func getRegistrationToken(mail string) string {
	var user entities.User
	_ = dbManager.Get("users").Collection.Find(bson.M{
		"email": mail,
	}).One(&user)

	if len(user.RegToken) == 0 {
		return "badToken"
	}
	return user.RegToken
}

func getContract(file string, skip int) *entities.Contract {
	var contract entities.Contract
	_ = dbManager.Get("contracts").Collection.Find(bson.M{
		"file.name": file,
	}).Sort("_id").Skip(skip).One(&contract)
	return &contract
}
