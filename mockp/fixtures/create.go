package fixtures

import (
	"dfss/dfssp/api"
)

var CreateFixture map[string]*api.ErrorCode = map[string]*api.ErrorCode{
	"success": &api.ErrorCode{
		Code: api.ErrorCode_SUCCESS,
	},
	"warning": &api.ErrorCode{
		Code:    api.ErrorCode_WARNING,
		Message: "Some users are not ready yet",
	},
}
