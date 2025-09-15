package errh

import (
	"log/slog"
	"net/http"

	grpc "google.golang.org/grpc/codes"
)

const (
	defaultErrMsg   = "Unexpected internal error"
	defaultGRPCCode = grpc.Internal
	defaultHTTPCode = http.StatusInternalServerError
)

type CodeByErrorType func(errorType string) grpc.Code

type LoggingByErrorType func(errorType string) LoggingLevel

type NotXerrCallback func(error) (error, bool)

type ErrorHandler struct {
	codes     CodeByErrorType
	notXerrFn NotXerrCallback

	logging LoggingByErrorType
	logger  *slog.Logger
}

func NewErrorHandler(
	logger *slog.Logger,
	codes CodeByErrorType,
	logging LoggingByErrorType,
	notXerrFn NotXerrCallback,
) *ErrorHandler {
	return &ErrorHandler{
		codes:     codes,
		logging:   logging,
		notXerrFn: notXerrFn,
		logger:    logger,
	}
}
