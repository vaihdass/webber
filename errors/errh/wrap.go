package errh

import (
	"errors"
	"fmt"

	"github.com/vaihdass/webber/errors/xerr"
)

func Wrap(operation string, xErr, cause error, kvs ...any) error {
	if len(kvs) > 0 {
		xErr = handleKV(xErr, kvs)
	}

	if cause == nil {
		return fmt.Errorf("%s: %w", operation, xErr)
	}

	return fmt.Errorf("%s: %w: %w", operation, xErr, cause)
}

func handleKV(err error, kvs []any) error {
	xErr, ok := xerr.From(err)
	if !ok {
		return err
	}

	// fill with map[key]value if xerr.Error (additional values for logs & traces)
	return newErrorValues(xErr, kvs)
}

type errorValues struct {
	error  *xerr.Error
	values []any
}

func (e errorValues) Error() string {
	return e.error.Error()
}

func (e errorValues) Unwrap() error {
	return e.error
}

func newErrorValues(err *xerr.Error, kvs []any) error {
	if len(kvs) == 0 || len(kvs)%2 != 0 {
		return err
	}

	return &errorValues{
		error:  err,
		values: kvs,
	}
}

func extractErrorValues(e error, kvs []any) []any {
	var err *errorValues
	if !errors.As(e, &err) {
		return kvs
	}

	return append(kvs, err.values...)
}
