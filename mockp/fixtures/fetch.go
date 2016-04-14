package fixtures

import (
	"dfss/dfssp/api"
)

var FetchFixture map[string]*api.Contract = map[string]*api.Contract{
	"01": &api.Contract{
		ErrorCode: &api.ErrorCode{Code: api.ErrorCode_SUCCESS},
		Json:      []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f}, // Hello
	},
	"02": &api.Contract{
		ErrorCode: &api.ErrorCode{Code: api.ErrorCode_BADAUTH},
	},
	"03": &api.Contract{
		ErrorCode: &api.ErrorCode{Code: api.ErrorCode_INVARG},
	},
}
