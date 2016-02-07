// Package fixtures contains the responses to the gRpc requests
package fixtures

import (
	pb "dfss/dfssp/api"
)

// RegisterFixture holds the fixture for the Register route
var RegisterFixture = map[string]*pb.ErrorCode{
	"dfss@success.io": &pb.ErrorCode{
		Code:    pb.ErrorCode_SUCCESS,
		Message: "SUCCESS",
	},
	"dfss@invarg.io": &pb.ErrorCode{
		Code:    pb.ErrorCode_INVARG,
		Message: "INVARG",
	},
	"dfss@badauth.io": &pb.ErrorCode{
		Code:    pb.ErrorCode_BADAUTH,
		Message: "BADAUTH",
	},
	"dfss@warning.io": &pb.ErrorCode{
		Code:    pb.ErrorCode_WARNING,
		Message: "WARNING",
	},
	"dfss@interr.io": &pb.ErrorCode{
		Code:    pb.ErrorCode_INTERR,
		Message: "INTERR",
	},
	"default": &pb.ErrorCode{
		Code:    pb.ErrorCode_INTERR,
		Message: "INTERR",
	},
}
