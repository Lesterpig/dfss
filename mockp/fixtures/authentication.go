package fixtures

import (
	pb "dfss/dfssp/api"
)

// AuthFixture holds the fixture for the Auth route
var AuthFixture = map[string]*pb.RegisteredUser{
	"default": &pb.RegisteredUser{ClientCert: "default"},
}
