package features

//go:generate parsley-cli generate proxy

type Greeter interface {
	SayHello(name string) (string, error)
}
