package E

import (
	"errors"
)

const (
	UnknownErrCode = 500
	NormalErrCode  = 400
)

var _ error = (*ApplicationError)(nil)

type ApplicationError struct {
	Err  error
	Code int16
}

func (e *ApplicationError) Error() string {
	return e.Err.Error()
}

func PanicErr(err error) {
	if err == nil {
		return
	}
	if _, ok := err.(*ApplicationError); ok {
		panic(err)
	}
	e := &ApplicationError{
		Code: UnknownErrCode,
		Err:  err,
	}
	panic(e)
}

func Message(message string) *ApplicationError {
	e := &ApplicationError{
		Code: NormalErrCode,
		Err:  errors.New(message),
	}
	return e
}
