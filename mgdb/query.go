package mgdb

import "reflect"

// Selector allow a user to build a selector like an struct
type Selector struct {
	factory  *MetadataFactory       // The MetadataFactory for getting nested metadata
	metadata *Metadata              // The Metadata for the current entity
	maps     map[string]interface{} // Contains the query selectors
	childs   map[string]*Selector   // Contains the nested type
	parent   *Selector              // Parent selector in the hierarchy
}

// NewSelector creates and return a new selector
func NewSelector(factory *MetadataFactory, metadata *Metadata, selector *Selector) *Selector {
	return &Selector{
		factory,
		metadata,
		make(map[string]interface{}),
		make(map[string]*Selector),
		selector,
	}
}

// AddChild adds a child element to a selector and returns it
func (s *Selector) AddChild(name string) *Selector {
	nestedType, ok := s.metadata.Nested[name]
	mapped, _ := s.metadata.Mapping[name]
	if ok {
		selector := NewSelector(s.factory, factory.Metadata(reflect.New(nestedType).Interface()), s)
		s.maps[mapped] = selector.childs
		s.childs[mapped] = selector
		return selector
	}
	return nil
}

// Parent go up for one level in the tree
func (s *Selector) Parent() *Selector {
	return s.parent
}

// Equal query
func (s *Selector) Equal(value interface{}) {
	s.maps["$eq"] = value
}

// Greater query
func (s *Selector) Greater(value interface{}) {
	s.maps["$gt"] = value
}

// Lower query
func (s *Selector) Lower(value interface{}) {
	s.maps["$lt"] = value
}

// Greater or equal query
func (s *Selector) GreaterOrEqual(value interface{}) {
	s.maps["$gte"] = value
}

// Lower or equal query
func (s *Selector) LowerOrEqual(value interface{}) {
	s.maps["$lte"] = value
}

// Not equal query
func (s *Selector) NotEqual(value interface{}) {
	s.maps["$ne"] = value
}

// In query (parameter must be an array)
func (s *Selector) In(value interface{}) {
	s.maps["$in"] = value
}

// Not in query (parameter must be an arry)
func (s *Selector) NotIn(value interface{}) {
	s.maps["$nin"] = value
}
