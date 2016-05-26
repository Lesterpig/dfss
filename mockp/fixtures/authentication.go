package fixtures

import (
	"dfss/dfssp/api"
)

// AuthFixture holds the fixture for the Auth route
var AuthFixture = map[string]*api.RegisteredUser{
	"default": &api.RegisteredUser{ClientCert: "default"},
}
