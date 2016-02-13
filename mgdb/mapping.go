package mgdb

import "reflect"

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
		metadata = newMetadata(factory, element, make(map[string]bool))
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

	// Nested holds types for nested structs if necessary
	Nested map[string]string
}

// NewMetadata instantiate the Metadata associated to the given struct
// It uses the `key` tag to do the mapping, more concrete
// examples are provided in the documentation
// Handles nested and recursive types
func newMetadata(factory *MetadataFactory, element interface{}, visited map[string]bool) *Metadata {
	m := make(map[string]string) // Go field name to mongo field name
	n := make(map[string]string) // Go field name to type name
	t := reflect.TypeOf(element)
	visited[t.String()] = true
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if tag := field.Tag.Get("key"); tag != "" {
			fieldType := field.Type
			m[field.Name] = tag
			kind := field.Type.Kind()
			// Handle nested arrays, slices, structs and pointers. Does not handle maps
			if kind == reflect.Struct {
				if _, ok := visited[fieldType.String()]; !ok {
					visited[fieldType.String()] = true
					factory.metadatas[fieldType.String()] = newMetadata(factory, reflect.New(fieldType).Elem().Interface(), visited)
				}
				n[field.Name] = fieldType.String()
			} else if kind == reflect.Array || kind == reflect.Slice || kind == reflect.Ptr {
				if _, ok := visited[fieldType.Elem().String()]; !ok {
					visited[fieldType.Elem().String()] = true
					factory.metadatas[fieldType.Elem().String()] = newMetadata(factory, reflect.Indirect(reflect.New(fieldType.Elem())).Interface(), visited)
				}
				n[field.Name] = fieldType.Elem().String()
			}
		}
	}
	return &Metadata{
		m,
		n,
	}
}
