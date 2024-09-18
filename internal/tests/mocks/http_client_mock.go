package mocks

import (
	"bytes"
	"github.com/matzefriedrich/parsley/internal/utils"
	"io"
	"net/http"
)

type HttpClientMock struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (h *HttpClientMock) Do(req *http.Request) (*http.Response, error) {
	return h.DoFunc(req)
}

var _ utils.HttpClient = (*HttpClientMock)(nil)

func NewHttpClientMock() *HttpClientMock {
	return &HttpClientMock{}
}

type HttpResponseMock struct {
	responseBuffer bytes.Buffer
}

func NewHttpResponseMock(buffer bytes.Buffer) *HttpResponseMock {
	return &HttpResponseMock{
		responseBuffer: buffer,
	}
}

func (h *HttpResponseMock) Read(p []byte) (n int, err error) {
	return h.responseBuffer.Read(p)
}

func (h *HttpResponseMock) Close() error {
	return nil
}

var _ io.ReadCloser = (*HttpResponseMock)(nil)
