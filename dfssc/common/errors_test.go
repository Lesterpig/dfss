package common

import (
	"testing"

	"dfss/dfssp/api"
	"github.com/bmizerany/assert"
)

func TestEvaluateErrorCodeResponse(t *testing.T) {

	success := &api.ErrorCode{
		Code:    api.ErrorCode_SUCCESS,
		Message: "Useless message",
	}

	err := EvaluateErrorCodeResponse(success)
	assert.Equal(t, nil, err)

	warning := &api.ErrorCode{
		Code:    api.ErrorCode_WARNING,
		Message: "Useful message",
	}

	err = EvaluateErrorCodeResponse(warning)
	assert.Equal(t, "Operation succeeded with a warning message: Useful message", err.Error())

	badauth := &api.ErrorCode{
		Code:    api.ErrorCode_BADAUTH,
		Message: "Useless message",
	}

	err = EvaluateErrorCodeResponse(badauth)
	assert.Equal(t, "Authentication error", err.Error())

	other := &api.ErrorCode{
		Code: api.ErrorCode_INTERR,
	}

	err = EvaluateErrorCodeResponse(other)
	assert.Equal(t, "Received error code INTERR", err.Error())

	otherWithMessage := &api.ErrorCode{
		Code:    api.ErrorCode_INVARG,
		Message: "Invalid mail",
	}

	err = EvaluateErrorCodeResponse(otherWithMessage)
	assert.Equal(t, "Invalid mail", err.Error())
}
