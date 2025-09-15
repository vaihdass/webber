package errh

import (
	"encoding/json"
	"errors"
	"fmt"

	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type TypedGRPCStatus struct {
	status  *status.Status
	errType string
}

func NewTypedGRPCStatus(st *status.Status, errType string) (*TypedGRPCStatus, error) {
	st, err := setGRPCErrorType(st, errType)
	if err != nil {
		return nil, fmt.Errorf("errh.NewTypedGRPCStatus: %w", err)
	}

	return &TypedGRPCStatus{
		status:  st,
		errType: errType,
	}, nil
}

func ExtractNewTypedGRPCStatus(st *status.Status) (*TypedGRPCStatus, error) {
	errType, ok := getGRPCErrorType(st)
	if !ok {
		return nil, errors.New("error type not found")
	}

	return &TypedGRPCStatus{
		status:  st,
		errType: errType,
	}, nil
}

func (s *TypedGRPCStatus) GRPCStatus() *status.Status {
	return s.status
}

func (s *TypedGRPCStatus) Type() string {
	return s.errType
}

func (s *TypedGRPCStatus) Error() string {
	return s.status.Message()
}

func (s *TypedGRPCStatus) Unwrap() error {
	return s.status.Err()
}

type errorType struct {
	ErrorType string `json:"xerr_error_type,omitempty"`
}

func setGRPCErrorType(st *status.Status, errType string) (*status.Status, error) {
	if st == nil {
		return nil, errors.New("nil status")
	}

	bytes, err := json.Marshal(errorType{ErrorType: errType})
	if err != nil {
		return nil, err
	}

	return st.WithDetails(wrapperspb.String(string(bytes)))
}

func getGRPCErrorType(st *status.Status) (string, bool) {
	if st == nil {
		return "", false
	}

	var errType errorType
	var parsed bool

	d := st.Details()
	for i := 0; i < len(d); i++ {
		str, ok := d[i].(*wrapperspb.StringValue)
		if !ok {
			continue
		}

		if err := json.Unmarshal([]byte(str.GetValue()), &errType); err != nil {
			continue
		}

		parsed = true
		break
	}

	return errType.ErrorType, parsed
}
