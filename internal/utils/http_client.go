package utils

import "net/http"

// HttpClient provides an interface for making HTTP requests, allowing for different implementations such as mocks or wrappers.
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}
