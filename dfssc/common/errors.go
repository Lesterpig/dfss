package common

import (
	"errors"

	"dfss/dfssp/api"
)

// EvaluateErrorCodeResponse converts an ErrorCode to human-friendly error
func EvaluateErrorCodeResponse(code *api.ErrorCode) error {

	switch code.Code {
	case api.ErrorCode_SUCCESS:
		return nil
	case api.ErrorCode_WARNING:
		return errors.New("Operation succeeded with a warning message: " + code.Message)
	case api.ErrorCode_BADAUTH:
		return errors.New("Authentication error")
	}

	if len(code.Message) == 0 {
		return errors.New("Received error code " + (code.Code).String())
	}
	return errors.New("Received error code " + (code.Code).String() + ": " + code.Message)
}
