package features

//go:generate parsley-cli generate proxy
//go:generate parsley-cli generate mocks

type Greeter interface {
	SayHello(name string, polite bool) (string, error)
	SayNothing()
}
