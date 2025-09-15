package buildin

// Ptr returns a pointer copy of value.
func Ptr[T any](x T) *T {
	return &x
}

// FromPtr returns the pointer value & true or pointer type zero value & false.
func FromPtr[T any](x *T) (T, bool) {
	if x == nil {
		var zero T
		return zero, false
	}

	return *x, true
}
