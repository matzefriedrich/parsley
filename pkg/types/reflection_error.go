package types

// ReflectionError represents an error specifically related to reflection operations, extending ParsleyError.
type ReflectionError struct {
	ParsleyError
}

// NewReflectionError creates a new ReflectionError with a specified message and optional initializers.
func NewReflectionError(msg string, initializers ...ParsleyErrorFunc) error {
	err := &ReflectionError{
		ParsleyError: ParsleyError{
			Msg: msg,
		},
	}
	for _, f := range initializers {
		f(&err.ParsleyError)
		f(err)
	}
	return err
}
