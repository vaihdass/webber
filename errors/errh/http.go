package errh

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/vaihdass/webber/errors/xerr"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *ErrorHandler) HandleHTTP(
	ctx context.Context, w http.ResponseWriter, r *http.Request,
	operation string, err error, options ...Option,
) {
	if err == nil {
		return
	}

	if operation != "" {
		err = fmt.Errorf("%s: %w", operation, err)
	}

	opts := configureOptions(options...)
	xErr, ok := xerr.From(err)

	// Fast path: untyped error (not xerr.Error)
	if !ok {
		code, msg := getUnexpectedHTTPError(err, h.notXerrFn, opts.fallbackMsg)
		setHTTPError(w, code, msg, xerr.UntypedErrType)
		return
	}

	// Happy path: all typed errors (xerr.Error)
	httpCode, grpcCode := getCodesByErrType(xErr.Type(), h.codes)
	logLvl := getLogLvlByErrType(xErr.Type(), h.logging)

	// Logging
	logValues := extractErrorValues(err, opts.values)
	log(ctx, h.logger, logLvl, grpcCode, err.Error(), xErr.Error(), xErr.Type(), logValues, opts.span)

	setHTTPError(w, httpCode, xErr.Error(), xErr.Type())
}

func getUnexpectedHTTPError(err error, notXerrFn NotXerrCallback, fallbackMsg string) (int, string) {
	var newErr error
	var modified bool

	if notXerrFn != nil {
		newErr, modified = notXerrFn(err)
	}

	if modified && newErr != nil {
		st := status.Convert(newErr)
		return runtime.HTTPStatusFromCode(st.Code()), st.Message()
	}

	if fallbackMsg == "" {
		fallbackMsg = defaultErrMsg
	}

	return defaultHTTPCode, fallbackMsg
}

func getCodesByErrType(errType string, cb CodeByErrorType) (int, codes.Code) {
	grpcCode := getCodeByErrType(errType, cb)

	// Default HTTP code for typed error without code configuration
	if grpcCode == codes.OK {
		return defaultHTTPCode, defaultGRPCCode
	}

	return runtime.HTTPStatusFromCode(grpcCode), grpcCode
}
