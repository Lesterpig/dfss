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
		animals [3]animal `key:"animals"`
		persons []person  `key:"persons"`
		ptr     *node     `key:"ptr"`
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
		animals,
		persons,
		ptr,
	}

	metadata := f.Metadata(n)
	nestedNode := metadata.Nested
	fmt.Printf("%s\n", nestedNode)
	// We check that all the nested types are there
	if len(nestedNode) != 4 {
		t.Error("Expected nested entities")
	}
	typeUser, ok := nestedNode["u"]
	if !ok {
		t.Error("Expected a nested user")
	}
	typeAnimal, ok := nestedNode["animals"]
	if !ok {
		t.Error("Expected a nested animal")
	}
	typePerson, ok := nestedNode["persons"]
	if !ok {
		t.Error("Expected a nested person")
	}
	typeNode, ok := nestedNode["ptr"]
	if !ok {
		t.Error("Expected a nested node")
	}

	// We check that the mapping has been done
	checkMapping(t, []string{"mail", "login"}, []string{"_id", "login"}, f.metadatas[typeUser])
	checkMapping(t, []string{"Name", "Race", "Age"}, []string{"Name", "Race", "Age"}, f.metadatas[typeAnimal])
	checkMapping(t, []string{"firstName", "lastName", "address", "ani"}, []string{"f", "l", "a", "ani"}, f.metadatas[typePerson])
	checkMapping(t, []string{"u", "animals", "persons", "ptr"}, []string{"u", "animals", "persons", "ptr"}, f.metadatas[typeNode])

}
