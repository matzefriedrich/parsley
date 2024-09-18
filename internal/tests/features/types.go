package features

import "net/http"

//go:generate parsley-cli generate proxy
//go:generate parsley-cli generate mocks

//parsley:mock
type Greeter interface {
	SayHello(name string, polite bool) (string, error)
	SayNothing()
}

//parsley:mock
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}
