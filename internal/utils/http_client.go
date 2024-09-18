package utils

import "net/http"

//go:generate parsley-cli generate mocks

//parsley:mock
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}
