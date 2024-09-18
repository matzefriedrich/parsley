package utils

import (
	"bytes"
	"context"
	"github.com/matzefriedrich/parsley/internal/tests/mocks"
	"github.com/matzefriedrich/parsley/internal/utils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func Test_GitHubApiClient_QueryLatestReleaseTag(t *testing.T) {

	// Arrange
	const data = "[ { \"id\": 123, \"tag_name\": \"v0.1.0\", \"name\": \"release v0.1.0\" } ]"
	buffer := bytes.Buffer{}
	buffer.WriteString(data)

	client := mocks.NewHttpClientMock()
	client.DoFunc = func(req *http.Request) (*http.Response, error) {
		body := mocks.NewHttpResponseMock(buffer)
		return &http.Response{
			Body:       body,
			StatusCode: http.StatusOK,
		}, nil
	}

	sut := utils.NewGitHubApiClient(client, func(options *utils.HttpClientOptions) {

	})

	// Act
	tag, err := sut.QueryLatestReleaseTag(context.Background())

	// Assert
	assert.Nil(t, err)
	assert.NotEmpty(t, tag)
}
