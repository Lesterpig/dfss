package mgdb

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"os"
	"testing"
)

type card struct {
	Id    bson.ObjectId `key:"_id" bson:"_id"`
	Value string        `key:"value" bson:"value"`
	Color string        `key:"color" bson:"color"`
}

type hand struct {
	Id      bson.ObjectId `key:"_id" bson:"_id"`
	CardOne card          `key:"card_one" bson:"card_one"`
	CardTwo card          `key:"card_two" bson:"card_two"`
}

var manager *MongoManager
var err error

func TestMain(m *testing.M) {
	// Setup
	manager, err = NewManager("localhost", "admin", "admin", "demo", "demo", 27017)

	// Run
	code := m.Run()

	// Teardown
	// The collection is created automatically on
	// first connection, that's why we do not recreate it manually
	manager.collection.DropCollection()
	manager.Close()

	os.Exit(code)
}

func TestConnection(t *testing.T) {
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
	ok, err := manager.Insert(c)
	if !ok {
		fmt.Println("A problem occurred during insert : ", err)
		return c, false
	}
	return c, true
}

func TestInsert(t *testing.T) {
	_, ok := helperInsert("five", "Hearts")
	if !ok {
		t.Fatal("Couldn't insert")
	}
}

func TestFindById(t *testing.T) {
	c, ok := helperInsert("king", "Spades")
	if !ok {
		t.Fatal("Couldn't insert")
	}

	res := card{}
	err := manager.FindById(c, &res)
	if err != nil {
		t.Fatal("Couldn't fetch the card : ", err)
	}
	fmt.Println("Fetched the card : ", res)
}

func TestUpdateById(t *testing.T) {
	// Create and insert new card
	c, ok := helperInsert("Ace", "Diamonds")
	if !ok {
		t.Fatal("Couldn't insert")
	}

	// Update and persist the card
	c.Value = "Jack"
	c.Color = ""
	ok, err := manager.UpdateById(c)
	if !ok {
		t.Fatal("Couldn't update the card : ", err)
	}
	fmt.Println("Updated to : ", c)

	// Assert the changes have been persisted
	res := card{}
	err = manager.FindById(c, &res)
	if err != nil {
		t.Fatal("Couldn't fetch the previously updated card")
	}
	fmt.Println("Fetched the card : ", res)

	if c.Id != res.Id || c.Color != res.Color || c.Value != res.Value {
		t.Fatal(fmt.Sprintf("Updated card with %v and fetched %v", c, res))
	}
}

func TestUpdateByIdNestedTypes(t *testing.T) {
	// Create and insert a hand
	c1 := card{bson.NewObjectId(), "Ace", "Spades"}
	c2 := card{bson.NewObjectId(), "Ace", "Hearts"}
	h := hand{bson.NewObjectId(), c1, c2}

	ok, err := manager.Insert(h)
	if !ok {
		t.Fatal("Couldn't insert hand :", err)
	}
	fmt.Println("Hand is : ", h)

	// Update the hand and persist the changes
	h.CardOne.Value = "Three"
	h.CardOne.Color = "Clubs"
	h.CardTwo.Value = "Six"
	h.CardTwo.Color = "Diamonds"
	ok, err = manager.UpdateById(h)
	if !ok {
		t.Fatal("An error occured while updating the hand :", err)
	}
	fmt.Println("Update hand to : ", h)

	// Find the hand and assert the changes were made
	res := hand{}
	err = manager.FindById(h, &res)
	if err != nil {
		t.Fatal("Couldn't fetch the previously update hand")
	}
	fmt.Println("Fetched hand : ", res)

	if h.CardOne.Value != res.CardOne.Value || h.CardTwo.Value != res.CardTwo.Value || h.CardOne.Color != res.CardOne.Color || h.CardTwo.Color != res.CardTwo.Color {
		t.Fatal(fmt.Sprintf("Fetched card is %v; expected %v", res, h))
	}
	fmt.Println("Update was successful")
}

func TestDeleteById(t *testing.T) {
	c, ok := helperInsert("Three", "Hearts")
	if !ok {
		t.Fatal("Couldn't insert")
	}

	ok, err := manager.DeleteById(c)
	if !ok {
		t.Fatal("Couldn't remove the card : ", err)
	}
	fmt.Println("Removed the card")
}
