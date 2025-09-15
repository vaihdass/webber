package errh

import (
	"errors"

	"github.com/vaihdass/webber/errors/xerr"
)

type wrapper struct {
	xErr error
	err  error
}

func (w *wrapper) Error() string {
	return w.err.Error()
}

func (w *wrapper) Unwrap() []error {
	return []error{w.xErr, w.err}
}

func newWrapper(xErr, err error) *wrapper {
	return &wrapper{xErr: xErr, err: err}
}

func TryRewrapTypedErr(err error, newMsg string) error {
	var ev *errorValues
	if errors.As(err, &ev) {
		ev.error = xerr.New(ev.error.Type(), newMsg)
		return newWrapper(ev, err)
	}

	xErr, ok := xerr.From(err)
	if !ok {
		return err
	}

	xErr = xerr.New(xErr.Type(), newMsg)

	return newWrapper(xErr, err)
}
