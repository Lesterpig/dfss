package mgdb

import (
	"reflect"
)

/****************
 MetadataFactory
****************/

// Represent a factory of metadata for structs
// Metadata are stored in a map, indexed by the struct name
type MetadataFactory struct {
	metadatas map[string]*Metadata
}

// Instantiate a new empty factory
func NewMetadataFactory() *MetadataFactory {
	return &MetadataFactory{make(map[string]*Metadata)}
}

// Get the Metadata associated to the struct
// When querying with a yet unknown struct, the associated Metadata is built
// If it is already known, just returns the stored Metadata
func (factory *MetadataFactory) Metadata(element interface{}) *Metadata {
	metadata, present := factory.metadatas[reflect.TypeOf(element).String()]
	if !present {
		metadata = NewMetadata(element)
		factory.metadatas[reflect.TypeOf(element).String()] = metadata
	}
	return metadata
}

// Using the metadata associated to the struct, returns the map
// of the struct. Keys are the database fields and values are the values
// stored in the struct
func (factory *MetadataFactory) ToMap(element interface{}) map[string]interface{} {
	data := factory.Metadata(element)
	m := make(map[string]interface{})
	v := reflect.ValueOf(element)
	for key, value := range data.Mapping {
		fieldValue := v.FieldByName(key).Interface()
		m[value] = fieldValue
	}
	return m
}

/*********
 Metadata
*********/

// Represents the metadata for a struct
type Metadata struct {

	// Map the go fields to the database fields
	Mapping map[string]string
}

// Instantiate the Metadata associated to the given struct
// It uses the `key` tag to do the mapping, more concrete
// examples are provided in the documentation
func NewMetadata(element interface{}) *Metadata {
	m := make(map[string]string)
	t := reflect.TypeOf(element)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if tag := field.Tag.Get("key"); tag != "" {
			m[field.Name] = tag
		}
	}
	return &Metadata{
		m,
	}
}
