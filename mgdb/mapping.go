package mgdb

import (
	"fmt"
	"reflect"
)

/****************
 MetadataFactory
****************/

// MetadataFactory is a factory of metadata for structs
// Metadata are stored in a map, indexed by the struct name
type MetadataFactory struct {
	metadatas map[string]*Metadata
}

// NewMetadataFactory instantiate a new empty factory
func NewMetadataFactory() *MetadataFactory {
	return &MetadataFactory{make(map[string]*Metadata)}
}

// Metadata get the Metadata associated to the struct
// When querying with a yet unknown struct, the associated Metadata is built
// If it is already known, just returns the stored Metadata
func (factory *MetadataFactory) Metadata(element interface{}) *Metadata {
	metadata, present := factory.metadatas[reflect.TypeOf(element).String()]
	if !present {
		metadata = newMetadata(element, make(map[string]bool))
		factory.metadatas[reflect.TypeOf(element).String()] = metadata
	}
	return metadata
}

// GetID gets the id field of an entity
func (factory *MetadataFactory) GetID(entity interface{}) interface{} {
	mapping := factory.Metadata(entity).Mapping
	for k, v := range mapping {
		if v == "_id" {
			return reflect.ValueOf(entity).FieldByName(k).Interface()
		}
	}
	return nil
}

/*********
 Metadata
*********/

// Metadata represents the metadata for a struct
type Metadata struct {

	// Mapping maps the go fields to the database fields
	Mapping map[string]string

	// Nested holds metadata for nested structs if necessary
	Nested map[string]*Metadata
}

// NewMetadata instantiate the Metadata associated to the given struct
// It uses the `key` tag to do the mapping, more concrete
// examples are provided in the documentation
// Handles nested and recursive types
func newMetadata(element interface{}, visited map[string]bool) *Metadata {
	m := make(map[string]string)
	n := make(map[string]*Metadata)
	t := reflect.TypeOf(element)
	fmt.Println("I build the Metadata for element " + t.String() + " of Kind " + t.Kind().String())
	visited[t.String()] = true
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if tag := field.Tag.Get("key"); tag != "" {
			fmt.Println("Processing field ", field.Name)
			m[field.Name] = tag
			if _, ok := visited[field.Type.String()]; !ok {
				kind := field.Type.Kind()
				// Handle nested types. Does not handle maps
				if kind == reflect.Struct {
					fmt.Println("I want the metadata of struct", reflect.New(field.Type).Elem().Interface())
					n[field.Name] = newMetadata(reflect.New(field.Type).Elem().Interface(), visited)
				} else if kind == reflect.Array || kind == reflect.Slice || kind == reflect.Ptr {
					fmt.Println("Studying type ", field.Type)
					fmt.Println("Field elem type is ", reflect.TypeOf(reflect.Indirect(reflect.New(field.Type.Elem())).Interface()))
					n[field.Name] = newMetadata(reflect.Indirect(reflect.New(field.Type.Elem())).Interface(), visited)
				}
			}
		}
	}
	return &Metadata{
		m,
		n,
	}
}
