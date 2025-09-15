package xerr

import (
	"errors"
)

func From(err error) (*Error, bool) {
	var e *Error

	if errors.As(err, &e) {
		return e, true
	}

	return nil, false
}

func HasType[T ~string](t T, err error) (*Error, bool) {
	e, ok := From(err)
	if !ok {
		return nil, false
	}

	return e, e.Type() == string(t)
}
