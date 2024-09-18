package commands

import (
	"bytes"
	"github.com/matzefriedrich/parsley/internal/commands"
	"github.com/matzefriedrich/parsley/internal/tests/mocks"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func Test_NewVersionCommand_Execute_check_for_update(t *testing.T) {

	// Arrange
	const data = "[ { \"id\": 123, \"tag_name\": \"v115.1.3\", \"name\": \"release v0.9.3\" } ]"
	buffer := bytes.Buffer{}
	buffer.WriteString(data)

	httpClientMock := mocks.NewHttpClientMock()
	httpClientMock.DoFunc = func(req *http.Request) (*http.Response, error) {
		body := mocks.NewHttpResponseMock(buffer)
		return &http.Response{
			Body:       body,
			StatusCode: http.StatusOK,
		}, nil
	}

	sut := commands.NewVersionCommand(httpClientMock)
	sut.SetArgs([]string{"--check-update"})

	// Act
	err := sut.Execute()

	// Assert
	assert.NoError(t, err)
}
