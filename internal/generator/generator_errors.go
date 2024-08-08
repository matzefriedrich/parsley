package generator

import "errors"

const (
	ErrorCannotGenerateProxies    = "cannot generate proxies"
	ErrorCannotExecuteTemplate    = "cannot execute template"
	ErrorFailedToOpenTemplateFile = "failed to open template file"
	ErrorTemplateFileNotFound     = "template file not found"
)

var (
	ErrCannotGenerateProxies    = errors.New(ErrorCannotGenerateProxies)
	ErrCannotExecuteTemplate    = errors.New(ErrorCannotExecuteTemplate)
	ErrFailedToOpenTemplateFile = errors.New(ErrorFailedToOpenTemplateFile)
	ErrTemplateFileNotFound     = errors.New(ErrorTemplateFileNotFound)
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
