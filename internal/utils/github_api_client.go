package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type GithubRelease struct {
	Id          uint64    `json:"id"`
	TagName     string    `json:"tag_name"`
	Name        string    `json:"name"`
	HtmlUrl     string    `json:"html_url"`
	PublishedAt time.Time `json:"published_at"`
}

func (r GithubRelease) TryParseVersionFromTag() (*VersionInfo, error) {
	version := r.TagName
	return tryParseVersionInfo(version)
}

type githubApiClient struct {
	httpClient HttpClient
	options    HttpClientOptions
}

type GitHubClient interface {
	QueryLatestReleaseTag(ctx context.Context) (*GithubRelease, error)
}

var _ GitHubClient = (*githubApiClient)(nil)

type HttpClientOptions struct {
	RequestTimeout time.Duration
}

type HttpClientOptionsFunc func(*HttpClientOptions)

func NewGitHubApiClient(httpClient HttpClient, config ...HttpClientOptionsFunc) GitHubClient {
	options := HttpClientOptions{
		RequestTimeout: 5 * time.Second,
	}
	for _, optionsFunc := range config {
		optionsFunc(&options)
	}
	return &githubApiClient{
		httpClient: httpClient,
		options:    options,
	}
}

// QueryLatestReleaseTag Queries the latest version from the GitHub releases endpoint and compares it against the current application version.
func (c *githubApiClient) QueryLatestReleaseTag(ctx context.Context) (*GithubRelease, error) {

	const owner = "matzefriedrich"
	const repo = "parsley"

	requestCtx, cancel := context.WithTimeout(ctx, c.options.RequestTimeout)
	defer cancel()

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases", owner, repo)
	request, err := http.NewRequestWithContext(requestCtx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	response, requestErr := c.httpClient.Do(request)
	if requestErr != nil {
		return nil, requestErr
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch latest release: %s", response.Status)
	}

	var releases []GithubRelease
	if unmarshalErr := json.NewDecoder(response.Body).Decode(&releases); unmarshalErr != nil {
		return nil, err
	}

	if len(releases) > 0 {
		latestRelease := releases[0]
		return &latestRelease, nil
	}

	return nil, errors.New("failed to retrieve release information")
}
