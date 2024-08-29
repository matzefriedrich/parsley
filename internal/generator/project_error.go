package generator

import "errors"

type ProjectError struct {
	msg   string
	cause error
}

const (
	errorCannotReadModFile             = "failed to read the go.mod file"
	errorCannotParseModFile            = "failed to parse go.mod file"
	errorFailedToAddRequiredDependency = "failed to add the required dependency"
	errorCannotFormatModFile           = "failed to format go.mod file"
	errorCannotWriteModFile            = "failed to write the go.mod file"
)

var (
	ErrCannotReadModFile             = errors.New(errorCannotReadModFile)
	ErrCannotParseModFile            = errors.New(errorCannotParseModFile)
	ErrFailedToAddRequiredDependency = errors.New(errorFailedToAddRequiredDependency)
	ErrCannotFormatModFile           = errors.New(errorCannotFormatModFile)
	ErrCannotWriteModFile            = errors.New(errorCannotWriteModFile)
)

func newProjectError(msg string, cause error) error {
	return ProjectError{
		msg:   msg,
		cause: cause,
	}
}

func (e ProjectError) Error() string {
	return e.msg
}

func (e ProjectError) Unwrap() error {
	return e.cause
}

func (e ProjectError) Is(err error) bool {
	return e.Error() == err.Error()
}
