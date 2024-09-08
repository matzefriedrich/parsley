package types

type reflectionError struct {
	ParsleyError
}

func NewReflectionError(msg string, initializers ...ParsleyErrorFunc) error {
	err := &reflectionError{
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
