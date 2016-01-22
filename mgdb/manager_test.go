package mgdb

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"os"
	"testing"
)

type card struct {
	ID    bson.ObjectId `key:"_id" bson:"_id"`
	Value string        `key:"value" bson:"value"`
	Color string        `key:"color" bson:"color"`
}

type hand struct {
	ID      bson.ObjectId `key:"_id" bson:"_id"`
	CardOne card          `key:"card_one" bson:"card_one"`
	CardTwo card          `key:"card_two" bson:"card_two"`
}

const defaultDBUrl = "MGDB_URL"

var collection *MongoCollection
var manager *MongoManager
var err error

func TestMain(m *testing.M) {
	// Setup
	fmt.Println("Try to connect to : " + os.Getenv(defaultDBUrl))

	db := os.Getenv(defaultDBUrl)
	if db == "" {
		db = "demo"
	}

	manager, err = NewManager(db, defaultDBUrl)

	collection = manager.Get("demo")

	// Run
	code := m.Run()

	// Teardown
	// The collection is created automatically on
	// first connection, that's why we do not recreate it manually
	err = collection.Drop()

	if err != nil {
		fmt.Println("An error occurred while droping the collection")
	}
	manager.Close()

	os.Exit(code)
}

func TestMongoConnection(t *testing.T) {
	if err != nil {
		t.Fatal("Couldn't connect to the database :", err)
	}
}

func helperInsert(value, color string) (card, bool) {
	c := card{
		bson.NewObjectId(),
		value,
		color,
	}

	fmt.Println("Inserting card : ", c)
	ok, err := collection.Insert(c)
	if !ok {
		fmt.Println("A problem occurred during insert : ", err)
		return c, false
	}
	return c, true
}

func TestMongoInsert(t *testing.T) {
	_, ok := helperInsert("five", "Hearts")
	if !ok {
		t.Fatal("Couldn't insert")
	}
}

func TestMongoFindByID(t *testing.T) {
	c, ok := helperInsert("king", "Spades")
	if !ok {
		t.Fatal("Couldn't insert")
	}

	res := card{}
	err := collection.FindByID(c, &res)
	if err != nil {
		t.Fatal("Couldn't fetch the card : ", err)
	}
	fmt.Println("Fetched the card : ", res)
}

func TestMongoUpdateByID(t *testing.T) {
	// Create and insert new card
	c, ok := helperInsert("Ace", "Diamonds")
	if !ok {
		t.Fatal("Couldn't insert")
	}

	// Update and persist the card
	c.Value = "Jack"
	c.Color = ""
	ok, err := collection.UpdateByID(c)
	if !ok {
		t.Fatal("Couldn't update the card : ", err)
	}
	fmt.Println("Updated to : ", c)

	// Assert the changes have been persisted
	res := card{}
	err = collection.FindByID(c, &res)
	if err != nil {
		t.Fatal("Couldn't fetch the previously updated card")
	}
	fmt.Println("Fetched the card : ", res)

	if c.ID != res.ID || c.Color != res.Color || c.Value != res.Value {
		t.Fatal(fmt.Sprintf("Updated card with %v and fetched %v", c, res))
	}
}

func TestMongoUpdateByIDNestedTypes(t *testing.T) {
	// Create and insert a hand
	c1 := card{bson.NewObjectId(), "Ace", "Spades"}
	c2 := card{bson.NewObjectId(), "Ace", "Hearts"}
	h := hand{bson.NewObjectId(), c1, c2}

	ok, err := collection.Insert(h)
	if !ok {
		t.Fatal("Couldn't insert hand :", err)
	}
	fmt.Println("Hand is : ", h)

	// Update the hand and persist the changes
	h.CardOne.Value = "Three"
	h.CardOne.Color = "Clubs"
	h.CardTwo.Value = "Six"
	h.CardTwo.Color = "Diamonds"
	ok, err = collection.UpdateByID(h)
	if !ok {
		t.Fatal("An error occured while updating the hand :", err)
	}
	fmt.Println("Update hand to : ", h)

	// Find the hand and assert the changes were made
	res := hand{}
	err = collection.FindByID(h, &res)
	if err != nil {
		t.Fatal("Couldn't fetch the previously update hand")
	}
	fmt.Println("Fetched hand : ", res)

	if h.CardOne.Value != res.CardOne.Value || h.CardTwo.Value != res.CardTwo.Value || h.CardOne.Color != res.CardOne.Color || h.CardTwo.Color != res.CardTwo.Color {
		t.Fatal(fmt.Sprintf("Fetched card is %v; expected %v", res, h))
	}
	fmt.Println("Update was successful")
}

func TestMongoDeleteByID(t *testing.T) {
	c, ok := helperInsert("Three", "Hearts")
	if !ok {
		t.Fatal("Couldn't insert")
	}

	ok, err := collection.DeleteByID(c)
	if !ok {
		t.Fatal("Couldn't remove the card : ", err)
	}
	fmt.Println("Removed the card")
}

func ExampleMongoManager() {

	//Define an animal to be use in further tests

	type animal struct {
		Name string `key:"_id" bson:"_id"`
		Race string `key:"race" bson:"race"`
		Age  int    `key:"age" bson:"age"`
	}

	//Initializes a MongoManager for the 'demo' database
	manager, err := NewManager("demo", defaultDBUrl)
	if err != nil { /* Handle error */
	}

	// Connects to the collection named 'animals'
	// If inexistant, it is created
	animals := manager.Get("animals")

	// Creates then insert a new animal into the collection
	tails := animal{"Tails", "Fox", 15}
	ok, _ := animals.Insert(tails)
	fmt.Println(fmt.Sprintf("Transaction went ok : %v", ok))

	// Get the previously inserted animal
	ani := animal{Name: "Tails"}
	res := animal{}
	err = animals.FindByID(ani, &res)
	if err != nil { /* Handle error */
	}

	// res now contains the animal {"Tails", "Fox", 15}
	// It is also possible to provided a struct with several field filled
	// For example, the following code would have produced the same result
	err = animals.FindByID(tails, &res)
	if err != nil { /* Handle error */
	}

	// Update an entity and persist the changes in the database
	res.Age += 2
	ok, _ = animals.UpdateByID(res)
	fmt.Println(fmt.Sprintf("Transaction went ok : %v", ok))

	// The database now contains the document {"_id": "Tails", "race": "Fox", age: 17}

	ok, _ = animals.DeleteByID(res)
	fmt.Println(fmt.Sprintf("Transaction went ok : %v", ok))

	// Tails has been successfully deleted from the database

	// Insert a bunch of data for following examples

	ok, _ = animals.Insert(animal{"Sonic", "Hedgehog", 12})
	fmt.Println(fmt.Sprintf("Transaction went ok : %v", ok))
	ok, _ = animals.Insert(animal{"Eggman", "Robot", 15})
	fmt.Println(fmt.Sprintf("Transaction went ok : %v", ok))
	ok, _ = animals.Insert(animal{"Amy", "Hedgehog", 12})
	fmt.Println(fmt.Sprintf("Transaction went ok : %v", ok))
	ok, _ = animals.Insert(animal{"Tails", "Fox", 12})
	fmt.Println(fmt.Sprintf("Transaction went ok : %v", ok))
	ok, _ = animals.Insert(animal{"Metal Sonic", "Robot", 14})
	fmt.Println(fmt.Sprintf("Transaction went ok : %v", ok))
	ok, _ = animals.Insert(animal{"Knuckles", "Echidna", 13})
	fmt.Println(fmt.Sprintf("Transaction went ok : %v", ok))
	ok, _ = animals.Insert(animal{"EggRobo", "Robot", 15})
	fmt.Println(fmt.Sprintf("Transaction went ok : %v", ok))
	ok, _ = animals.Insert(animal{"Tikal", "Echidna", 14})
	fmt.Println(fmt.Sprintf("Transaction went ok : %v", ok))
	ok, _ = animals.Insert(animal{"Shadow", "Hedgehog", 13})
	fmt.Println(fmt.Sprintf("Transaction went ok : %v", ok))
	ok, _ = animals.Insert(animal{"Silver", "Hedgehog", 15})
	fmt.Println(fmt.Sprintf("Transaction went ok : %v", ok))

	// Get all documents in the collection
	var all []animal
	err = animals.FindAll(bson.M{}, &all)
	if err != nil { /* Handle error */
	}

	// Get all hedgehogs
	// The type bson.M is provided by mgo/bson, it is an alias for map[string]interface{}
	// To learn how to make a proper query, just refer to mongoDB documentation
	var hedgehogs []animal
	err = animals.FindAll(bson.M{"race": "Hedgehog"}, &hedgehogs)
	if err != nil { /* Handle error */
	}

	// Fetch Tails, Eggman and Silver
	var tailsEggmanSilver []animal
	names := make([]string, 3)
	names[0] = "Tails"
	names[1] = "Eggman"
	names[2] = "Silver"
	err = animals.FindAll(bson.M{"_id": bson.M{"$in": names}}, &tailsEggmanSilver)
	if err != nil { /* Handle error */
	}

	// Update all animals with age > 12 and decrement by one
	// The first argument is used to select some documents, and the second argument contains the modification to apply
	count, _ := animals.UpdateAll(bson.M{"age": bson.M{"$gt": 12}}, bson.M{"$inc": bson.M{"age": -1}})
	fmt.Println(fmt.Sprintf("%d animals were uodated", count))

	// UpdateAll animals with race = 'Robot' and change it to 'Machine'
	count, _ = animals.UpdateAll(bson.M{"race": "Robot"}, bson.M{"$set": bson.M{"race": "Machine"}})
	fmt.Println(fmt.Sprintf("%d animals were uodated", count))

	// Delete all hedgehogs
	count, _ = animals.DeleteAll(bson.M{"race": "Hedgehog"})
	fmt.Println(fmt.Sprintf("%d animals were uodated", count))

	// Drop all the collection
	// Be careful when using this
	err = animals.Drop()
	if err != nil { /* Handle error */
	}
}
