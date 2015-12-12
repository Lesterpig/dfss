package mgdb

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"os"
)

type errorConnection struct {
	s string
}

func (e *errorConnection) Error() string {
	return e.s
}

func NewErrorConnection(s string) error {
	return &errorConnection{s}
}

// The Manager handling mongoDB connection
type MongoManager struct {
	session     *mgo.Session
	database    *mgo.Database
	collections map[string]*MongoCollection
}

// Create a new Manager, the environment variable MONGOHQ_URL needs to be set
// up with mongo uri, else it throws an error
func NewManager(database string) (*MongoManager, error) {
	uri := os.Getenv("MONGOHQ_URL")
	if uri == "" {
		err := NewErrorConnection("No uri provided, please set the MONGOHG_URL to connect to mongo")
		return nil, err
	}

	sess, err := mgo.Dial(uri)

	if err != nil {
		return nil, err
	}

	db := sess.DB(database)
	return &MongoManager{
		sess,
		db,
		make(map[string]*MongoCollection),
	}, nil
}

// Get a MongoCollection over a specified collection
// It is then possible to query via the MongoCollection
func (m *MongoManager) Get(collection string) *MongoCollection {
	coll, ok := m.collections[collection]
	if !ok {
		coll = NewCollection(m.database.C(collection))
		m.collections[collection] = coll
	}
	return coll
}

// The Manager handling mongoDB querying
type MongoCollection struct {
	collection *mgo.Collection
	factory    *MetadataFactory
}

// Open a connection to the database
func NewCollection(coll *mgo.Collection) *MongoCollection {
	return &MongoCollection{
		coll,
		NewMetadataFactory(),
	}
}

// Inserts an Entity into the selected collection
// The _id field must be present in the mapping (see example provided)
func (manager *MongoCollection) Insert(entity interface{}) (bool, error) {
	err := manager.collection.Insert(entity)
	return err == nil, err
}

// Update the entity with the new value provided.
// The _id of an Entity cannot be changed this way
func (manager *MongoCollection) UpdateById(entity interface{}) (bool, error) {
	m := manager.factory.ToMap(entity)
	err := manager.collection.Update(map[string]interface{}{"_id": m["_id"]}, entity)
	return err == nil, err
}

// Update the entities matching the selector with the query
// The format of the parameters is expected to follow the one
// provided in mgo's documentation
// Return the number of updated entities
func (manager *MongoCollection) UpdateAll(selector interface{}, update interface{}) (int, error) {
	info, err := manager.collection.UpdateAll(selector, update)
	return info.Updated, err
}

// Given an id, result will contain the associated entity
func (manager *MongoCollection) FindById(id interface{}, result interface{}) error {
	m := manager.factory.ToMap(id)
	err := manager.collection.Find(map[string]interface{}{"_id": m["_id"]}).One(result)
	return err
}

// Finds all entities matching the selector and put them into the result slice
// The format of the selector is expected to follow the one
// provided in mgo's documentation
func (manager *MongoCollection) FindAll(query interface{}, results []interface{}) error {
	return manager.collection.Find(query).All(results)
}

// Delete the entity matching the id
// Return true if the delection was successful
func (manager *MongoCollection) DeleteById(id interface{}) (bool, error) {
	m := manager.factory.ToMap(id)
	err := manager.collection.Remove(bson.M{"_id": m["_id"]})
	return err == nil, err
}

// Delete all the entities matching the selector
// The format of the selector is expected to follow the one
// provided in mgo's documentation
// Return the number of deleted entities
func (manager *MongoCollection) DeleteAll(query interface{}) (int, error) {
	info, err := manager.collection.RemoveAll(query)
	return info.Removed, err
}

// Switch collection in the database
func (manager *MongoCollection) SwitchCollection(coll string) {
	manager.collection = manager.database.C(coll)
}

// Count the number of entities currently in the database
func (manager *MongoCollection) Count() int {
	count, _ := manager.collection.Count()
	return count
}

// Close the connection to the session
func (manager *MongoCollection) Close() {
	manager.session.Close()
}
