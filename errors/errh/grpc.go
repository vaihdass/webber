package errh

import (
	"context"
	"fmt"

	grpc "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vaihdass/webber/errors/xerr"
)

func (h *ErrorHandler) Handle(ctx context.Context, operation string, err error, options ...Option) error {
	if err == nil {
		return nil
	}

	if operation != "" {
		err = fmt.Errorf("%s: %w", operation, err)
	}

	opts := configureOptions(options...)
	xErr, ok := xerr.From(err)

	// fast path: untyped error (not xerr.Error)
	if !ok {
		return getUnexpectedError(err, h.notXerrFn, opts.fallbackMsg)
	}

	// happy path: all typed errors (xerr.Error)
	code := getCodeByErrType(xErr.Type(), h.codes)
	logLvl := getLogLvlByErrType(xErr.Type(), h.logging)

	// default GRPC code for typed error without code configuration
	if code == grpc.OK {
		code = defaultGRPCCode
	}

	// logging
	logValues := extractErrorValues(err, opts.values)
	log(ctx, h.logger, logLvl, code, err.Error(), xErr.Error(), xErr.Type(), logValues, opts.span)

	// create typed GRPC status
	st := status.New(code, xErr.Error())

	typed, err := NewTypedGRPCStatus(st, xErr.Type())
	if err != nil {
		return st.Err()
	}

	return typed
}

func getUnexpectedError(err error, notXerrFn NotXerrCallback, fallbackMsg string) error {
	var newErr error
	var modified bool

	if notXerrFn != nil {
		newErr, modified = notXerrFn(err)
	}

	if modified && newErr != nil {
		return newErr
	}

	if fallbackMsg == "" {
		fallbackMsg = defaultErrMsg
	}

	st := status.New(defaultGRPCCode, fallbackMsg)

	typed, err := NewTypedGRPCStatus(st, xerr.UntypedErrType)
	if err != nil {
		return st.Err()
	}

	return typed
}

func getCodeByErrType(errType string, cb CodeByErrorType) grpc.Code {
	if cb == nil {
		return grpc.OK
	}

	return cb(errType)
}

func getLogLvlByErrType(errType string, cb LoggingByErrorType) LoggingLevel {
	if cb == nil {
		return UnknownLogging
	}

	return cb(errType)
}
