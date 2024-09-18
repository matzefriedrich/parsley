package utils

import (
	"bytes"
	"context"
	"github.com/matzefriedrich/parsley/internal/utils"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

func Test_GitHubApiClient_QueryLatestReleaseTag(t *testing.T) {

	// Arrange
	const data = "[ { \"id\": 123, \"tag_name\": \"v0.1.0\", \"name\": \"release v0.1.0\" } ]"
	buffer := bytes.Buffer{}
	buffer.WriteString(data)

	body := &httpResponseMock{responseBuffer: buffer}
	client := &httpClientMock{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				Body:       body,
				StatusCode: http.StatusOK,
			}, nil
		},
	}
	sut := utils.NewGitHubApiClient(client, func(options *utils.HttpClientOptions) {

	})

	// Act
	tag, err := sut.QueryLatestReleaseTag(context.Background())

	// Assert
	assert.Nil(t, err)
	assert.NotEmpty(t, tag)
}

type httpClientMock struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (h *httpClientMock) Do(req *http.Request) (*http.Response, error) {
	return h.DoFunc(req)
}

var _ utils.HttpClient = (*httpClientMock)(nil)

type httpResponseMock struct {
	responseBuffer bytes.Buffer
}

func (h *httpResponseMock) Read(p []byte) (n int, err error) {
	return h.responseBuffer.Read(p)
}

func (h *httpResponseMock) Close() error {
	return nil
}

var _ io.ReadCloser = (*httpResponseMock)(nil)
