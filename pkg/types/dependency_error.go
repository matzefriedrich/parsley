package types

import "errors"

const (
	ErrorInstanceAlreadySet = "instance already set"
)

var (
	ErrInstanceAlreadySet = errors.New(ErrorInstanceAlreadySet)
)

type DependencyError struct {
	ParsleyError
}

func NewDependencyError(msg string) error {
	err := DependencyError{ParsleyError{Msg: msg}}
	return err
}
