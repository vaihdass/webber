package buildin

// Ptr returns a pointer copy of value.
func Ptr[T any](value T) *T {
	return &value
}

// FromPtr returns the pointer value & true or pointer type zero value & false.
func FromPtr[T any](ptr *T) (T, bool) {
	if ptr == nil {
		var zero T
		return zero, false
	}

	return *ptr, true
}
