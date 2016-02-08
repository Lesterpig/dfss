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

func checkMapping(t *testing.T, fields, expected []string, metadata *Metadata) {
	if len(fields) != len(metadata.Mapping) {
		t.Error("Invalid length of mapping")
	}
	for idx, val := range fields {
		if v, _ := metadata.Mapping[val]; v != expected[idx] {
			t.Error("Expected %s to be mapped into %s, got %s", val, expected[idx], v)
		}
	}
}

func TestPublicFields(t *testing.T) {
	a := animal{"pika", "mouse", 15}
	data := factory.Metadata(a)

	if data == nil {
		t.Fatal("Metadata is nil")
	}

	fmt.Println(data.Mapping)
	checkMapping(t, []string{"Name", "Race", "Age"}, []string{"Name", "Race", "Age"}, data)
}

func TestPrivateFields(t *testing.T) {
	p := person{"Ariel", "Undomiel", "Fondcombe", animal{}}
	data := factory.Metadata(p)

	if data == nil {
		t.Fatal("Metadata is nil")
	}

	fmt.Println(data.Mapping)
	checkMapping(t, []string{"firstName", "lastName", "address", "ani"}, []string{"f", "l", "a", "ani"}, data)
}

func TestUnmappedFields(t *testing.T) {
	u := user{"sonic@hedgehog.io", "chaos", "emerald"}
	data := factory.Metadata(u)

	if data == nil {
		t.Fatal("Metadata is nil")
	}

	fmt.Println(data.Mapping)
	checkMapping(t, []string{"mail", "login"}, []string{"_id", "login"}, data)
}

func TestNestedTypes(t *testing.T) {
	// Struct to check right mapping with nested types
	type node struct {
		u       user      `key:"u"`
		ptr     *node     `key:"ptr"`
		animals [3]animal `key:"animals"`
		persons []person  `key:"persons"`
	}

	f := NewMetadataFactory()

	var animals [3]animal
	persons := make([]person, 1)
	ptr := &node{}
	user := user{}

	animals[0] = animal{"carapuce", "turtle", 10}
	animals[1] = animal{"salameche", "lizard", 10}
	persons = append(persons, person{"Adam", "Douglas", "Earth", animals[0]})
	n := node{
		user,
		ptr,
		animals,
		persons,
	}

	metadata := f.Metadata(n)
	nestedNode := metadata.Nested
	if len(nestedNode) != 3 {
		t.Error("Expected only 3 nested entities")
	}
	metaUser, ok := nestedNode["u"]
	if !ok {
		t.Error("Expected mapping of User entity")
	}
	metaAnimals, ok := nestedNode["animals"]
	if !ok {
		t.Error("Expected mapping of Animal entity")
	}
	metaPersons, ok := nestedNode["persons"]
	if !ok {
		t.Error("Expected mapping of Person entity")
	}

	checkMapping(t, []string{"mail", "login"}, []string{"_id", "login"}, metaUser)
	checkMapping(t, []string{"Name", "Race", "Age"}, []string{"Name", "Race", "Age"}, metaAnimals)
	checkMapping(t, []string{"firstName", "lastName", "address", "ani"}, []string{"f", "l", "a", "ani"}, metaPersons)

}
