package generator

import "errors"

const (
	ErrorCannotExecuteTemplate      = "cannot execute template"
	ErrorCannotFormatGeneratedCode  = "cannot format generated code"
	ErrorCannotGenerateProxies      = "cannot generate proxies"
	ErrorFailedToOpenTemplateFile   = "failed to open template file"
	ErrorFailedToWriteGeneratedCode = "failed to write generated code"
	ErrorTemplateFileNotFound       = "template file not found"
)

var (
	ErrCannotExecuteTemplate      = errors.New(ErrorCannotExecuteTemplate)
	ErrCannotFormatGeneratedCode  = errors.New(ErrorCannotFormatGeneratedCode)
	ErrCannotGenerateProxies      = errors.New(ErrorCannotGenerateProxies)
	ErrFailedToOpenTemplateFile   = errors.New(ErrorFailedToOpenTemplateFile)
	ErrFailedToWriteGeneratedCode = errors.New(ErrorFailedToWriteGeneratedCode)
	ErrTemplateFileNotFound       = errors.New(ErrorTemplateFileNotFound)
)

type generatorError struct {
	msg   string
	cause error
}

func (g generatorError) Error() string {
	return g.msg
}

func (g generatorError) Unwrap() error {
	return g.cause
}

func (g generatorError) Is(err error) bool {
	return g.Error() == err.Error()
}

var _ error = &generatorError{}

func newGeneratorError(msg string, initializers ...func(error)) error {
	err := &generatorError{
		msg: msg,
	}
	for _, initializer := range initializers {
		initializer(err)
	}
	return err
}

func WithCause(err error) func(target error) {
	return func(target error) {
		var errWithCause *generatorError
		if errors.As(target, errWithCause) {
			errWithCause.cause = err
		}
	}
}
