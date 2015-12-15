package mgdb

import (
	"fmt"
	"testing"
)

type animal struct {
	Name string `key:"Name"`
	Race string `key:"Race"`
	Age  int    `key:"Age"`
}

type person struct {
	firstName string `key:"f"`
	lastName  string `key:"l"`
	address   string `key:"a"`
	ani       animal `key:"ani"`
}

type user struct {
	mail     string `key:"_id"`
	login    string `key:"login"`
	password string
}

var factory = NewMetadataFactory()

func TestPublicFields(t *testing.T) {
	a := animal{"pika", "mouse", 15}
	data := factory.Metadata(a)

	if data == nil {
		t.Fatal("Metadata is nil")
	}

	mapping := data.Mapping
	fmt.Println(mapping)
	if len(mapping) != 3 || mapping["Name"] != "Name" || mapping["Race"] != "Race" || mapping["Age"] != "Age" {
		t.Fatal("Mapping not correctly built")
	}
}

func TestPrivateFields(t *testing.T) {
	p := person{"Ariel", "Undomiel", "Fondcombe", animal{}}
	data := factory.Metadata(p)

	if data == nil {
		t.Fatal("Metadata is nil")
	}

	mapping := data.Mapping
	fmt.Println(mapping)
	if len(mapping) != 4 || mapping["firstName"] != "f" || mapping["lastName"] != "l" || mapping["address"] != "a" || mapping["ani"] != "ani" {
		t.Fatal("Mapping not correctly built")
	}
}

func TestUnmappedFields(t *testing.T) {
	u := user{"sonic@hedgehog.io", "chaos", "emerald"}
	data := factory.Metadata(u)

	if data == nil {
		t.Fatal("Metadata is nil")
	}

	mapping := data.Mapping
	fmt.Println(mapping)
	if len(mapping) != 2 || mapping["mail"] != "_id" || mapping["login"] != "login" {
		t.Fatal("Mapping not correctly built")
	}
}
