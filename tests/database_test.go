package tests

import (
	"dfss/dfssp/entities"
	"dfss/mgdb"
	"gopkg.in/mgo.v2/bson"
)

var dbURI string
var dbManager *mgdb.MongoManager

// EraseDatabase drops the test database.
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
