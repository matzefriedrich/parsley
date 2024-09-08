// Code generated by parsley-cli; DO NOT EDIT.
//
// This file was automatically generated and any changes to it will be overwritten.

package features

import (
	"github.com/matzefriedrich/parsley/pkg/features"
)

type greeterMock struct {
	features.MockBase
	SayHelloFunc   SayHelloFunc
	SayNothingFunc SayNothingFunc
}

type SayHelloFunc func(name string) (string, error)
type SayNothingFunc func()

func (m *greeterMock) SayHello(name string) (string, error) {
	m.TraceMethodCall("SayHello", name)
	return m.SayHelloFunc(name)
}

func (m *greeterMock) SayNothing() {
	m.TraceMethodCall("SayNothing")
	m.SayNothingFunc()
}

var _ Greeter = (*greeterMock)(nil)

// NewGreeterMock Creates a new configurable greeterMock object.
func NewGreeterMock() *greeterMock {
	mock := &greeterMock{
		MockBase: features.NewMockBase(),
		SayHelloFunc: func(name string) (string, error) {
			var result0 string
			var result1 error
			return result0, result1
		},
		SayNothingFunc: func() {},
	}
	mock.AddFunction("SayHello", "SayHello(name string) (string, error)")
	mock.AddFunction("SayNothing", "SayNothing()")
	return mock
}
