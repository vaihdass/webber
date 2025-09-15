package errh

import (
	"context"
	"errors"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/vaihdass/webber/errors/xerr"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPCToHTTPMiddleware is an error handler for HTTP gateway, sets typed HTTP error.
func GRPCToHTTPMiddleware(
	_ context.Context, _ *runtime.ServeMux, _ runtime.Marshaler,
	w http.ResponseWriter, _ *http.Request, err error,
) {
	code, msg, errType := extractError(err)

	setHTTPError(w, runtime.HTTPStatusFromCode(code), msg, errType)
}

func extractError(err error) (codes.Code, string, string) {
	var terr *TypedGRPCStatus
	if errors.As(err, &terr) {
		return terr.GRPCStatus().Code(), terr.Error(), terr.Type()
	}

	st, ok := status.FromError(err)
	if !ok {
		return defaultGRPCCode, defaultErrMsg, xerr.UntypedErrType
	}

	return st.Code(), st.Message(), xerr.UntypedErrType
}
