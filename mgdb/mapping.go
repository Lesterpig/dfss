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
		metadata = newMetadata(element, make(map[string]bool))
		factory.metadatas[reflect.TypeOf(element).String()] = metadata
	}
	return metadata
}

// ToMap uses the metadata associated to the struct to returns the map
// of the struct. Keys are the database fields and values are the values
// stored in the struct
func (factory *MetadataFactory) ToMap(element interface{}) map[string]interface{} {
	data := factory.Metadata(element)
	m := make(map[string]interface{})
	v := reflect.ValueOf(element)
	for key, value := range data.Mapping {
		fieldValue := v.FieldByName(key).Interface()
		if _, ok := data.Nested[key]; ok {
			m[value] = factory.ToMap(fieldValue)
		} else {
			m[value] = fieldValue
		}
	}
	return m
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
	visited[t.String()] = true
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if tag := field.Tag.Get("key"); tag != "" {
			m[field.Name] = tag
			if _, ok := visited[field.Type.String()]; !ok {
				kind := field.Type.Kind()
				// Handle nested types. Does not handle maps
				if kind == reflect.Struct {
					n[field.Name] = newMetadata(reflect.New(field.Type).Elem().Interface(), visited)
				} else if kind == reflect.Array || kind == reflect.Slice {
					n[field.Name] = newMetadata(reflect.New(field.Type.Elem()).Elem().Interface(), visited)
				}
			}
		}
	}
	return &Metadata{
		m,
		n,
	}
}
