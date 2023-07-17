package itemHandler

import (
	"errors"
)

const (
	MethodNotAllowed        = "request method not allowed"
	BadJSON                 = "bad json data"
	BadQuery                = "bad url query data"
	NotFound                = "entity not found"
	SomeJSONFiledsEmpty     = "some fields empty"
	UnexpectedInternalError = "unexpected internal error"
)

func unwrapErr(err error, target string) bool {
	for i := 0; i < 4; i++ {
		if err == nil {
			break
		}
		err = errors.Unwrap(err)
		if err.Error() == target {
			return true
		}
	}
	return false
}
