package generator

import (
	"errors"
	"github.com/matzefriedrich/parsley/pkg/types"
)

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
	types.ParsleyError
}

var _ error = &generatorError{}

func newGeneratorError(msg string, initializers ...types.ParsleyErrorFunc) error {
	err := &generatorError{
		ParsleyError: types.ParsleyError{
			Msg: msg,
		},
	}
	for _, initializer := range initializers {
		initializer(&err.ParsleyError)
		initializer(err)
	}
	return err
}
