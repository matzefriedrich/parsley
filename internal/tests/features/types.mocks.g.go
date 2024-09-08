package features

type greeterMock struct {
	sayHello   SayHelloFunc
	sayNothing SayNothingFunc
}

type SayHelloFunc func(name string) (string, error)
type SayNothingFunc func()

func (m *greeterMock) SayHello(name string) (string, error) {
	return m.sayHello(name)
}

func (m *greeterMock) SayNothing() {
	m.sayNothing()
}

var _ Greeter = (*greeterMock)(nil)

// NewGreeterMock Creates a new configurable greeterMock object.
func NewGreeterMock(
	sayHello SayHelloFunc,
	sayNothing SayNothingFunc,
) *greeterMock {
	return &greeterMock{
		sayHello:   sayHello,
		sayNothing: sayNothing,
	}
}
