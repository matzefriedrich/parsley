package utils

import (
	"context"
	"github.com/matzefriedrich/parsley/internal/utils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func Test_GitHubApiClient_QueryLatestReleaseTag(t *testing.T) {

	// Arrange
	client := &http.Client{}
	sut := utils.NewGitHubApiClient(client, func(options *utils.HttpClientOptions) {

	})

	// Act
	tag, err := sut.QueryLatestReleaseTag(context.Background())

	// Assert
	assert.Nil(t, err)
	assert.NotEmpty(t, tag)
}
