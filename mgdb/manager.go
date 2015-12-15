package mgdb

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"os"
)

// errorConnection represents an error to be thrown upon connection
type errorConnection struct {
	s string
}

// Return the string contained in the error
func (e *errorConnection) Error() string {
	return e.s
}

// Creates a new error with the given message
func newErrorConnection(message string) error {
	return &errorConnection{message}
}

// MongoManager is aimed at handling the Mongo connection through the mgo driver
type MongoManager struct {

	// Session is the mgo.Session struct
	Session *mgo.Session

	// Database is the mgo.Database struct
	Database *mgo.Database

	Collections map[string]*MongoCollection
}

// NewManager a new Manager, the environment variable MONGOHQ_URL needs to be set
// up with mongo uri, else it throws an error
func NewManager(database string) (*MongoManager, error) {
	uri := os.Getenv("MGDB_URL")
	if uri == "" {
		err := newErrorConnection("No uri provided, please set the MONGOHG_URL to connect to mongo")
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

// Close closes the current connection
// Be careful, you won't be able to query the Collections anymore
func (m *MongoManager) Close() {
	m.Session.Close()
}

// Get returns a MongoCollection over a specified Collection
// The Collections are cached when they are called at least once
func (m *MongoManager) Get(Collection string) *MongoCollection {
	coll, ok := m.Collections[Collection]
	if !ok {
		coll = newCollection(m.Database.C(Collection))
		m.Collections[Collection] = coll
	}
	return coll
}

// MongoCollection is a wrapped around an mgo Collection to query to database
type MongoCollection struct {
	// Collection is the mgo.Collection struct
	Collection *mgo.Collection
	factory    *MetadataFactory
}

// newCollection returns a new MongoCollection
func newCollection(coll *mgo.Collection) *MongoCollection {
	return &MongoCollection{
		coll,
		NewMetadataFactory(),
	}
}

// Insert persists an Entity into the selected Collection
// The _id field must be present in the mapping (see example provided)
func (manager *MongoCollection) Insert(entity interface{}) (bool, error) {
	err := manager.Collection.Insert(entity)
	return err == nil, err
}

// UpdateByID updates the entity with the new value provided.
// The _id of an Entity cannot be changed this way
func (manager *MongoCollection) UpdateByID(entity interface{}) (bool, error) {
	m := manager.factory.ToMap(entity)
	err := manager.Collection.Update(map[string]interface{}{"_id": m["_id"]}, entity)
	return err == nil, err
}

// UpdateAll updates the entities matching the selector with the query
// The format of the parameters is expected to follow the one
// provided in mgo's documentation
// Return the number of updated entities
func (manager *MongoCollection) UpdateAll(selector interface{}, update interface{}) (int, error) {
	info, err := manager.Collection.UpdateAll(selector, update)
	return info.Updated, err
}

// FindByID fill the entity from the document with matching id
func (manager *MongoCollection) FindByID(id interface{}, result interface{}) error {
	m := manager.factory.ToMap(id)
	err := manager.Collection.Find(map[string]interface{}{"_id": m["_id"]}).One(result)
	return err
}

// FindAll finds all entities matching the selector and put them into the result slice
// The format of the selector is expected to follow the one
// provided in mgo's documentation
func (manager *MongoCollection) FindAll(query interface{}, result interface{}) error {
	return manager.Collection.Find(query).All(result)
}

// DeleteByID deletes the entity matching the id
// Return true if the delection was successful
func (manager *MongoCollection) DeleteByID(id interface{}) (bool, error) {
	m := manager.factory.ToMap(id)
	err := manager.Collection.Remove(bson.M{"_id": m["_id"]})
	return err == nil, err
}

// DeleteAll deletes all the entities matching the selector
// The format of the selector is expected to follow the one
// provided in mgo's documentation
// Return the number of deleted entities
func (manager *MongoCollection) DeleteAll(query interface{}) (int, error) {
	info, err := manager.Collection.RemoveAll(query)
	return info.Removed, err
}

// Count returns the number of entities currently in the Collection
func (manager *MongoCollection) Count() int {
	count, _ := manager.Collection.Count()
	return count
}

// Drop drops the current Collection
// This action is irreversible !
func (manager *MongoCollection) Drop() error {
	return manager.Collection.DropCollection()
}
