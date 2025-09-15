package xerr

const UntypedErrType = "untyped_error"

type Error struct {
	errType string
	message string
}

func New[T ~string](errorType T, message string) *Error {
	if string(errorType) == "" {
		errorType = UntypedErrType
	}

	return &Error{
		errType: string(errorType),
		message: message,
	}
}

func (e *Error) Error() string {
	return e.message
}

func (e *Error) Type() string {
	return e.errType
}
