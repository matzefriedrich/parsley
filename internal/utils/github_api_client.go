package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	options HttpClientOptions
}

type HttpClientOptions struct {
	RequestTimeout time.Duration
}

type HttpClientOptionsFunc func(*HttpClientOptions)

func NewGitHubApiClient(config ...HttpClientOptionsFunc) *githubApiClient {
	options := HttpClientOptions{
		RequestTimeout: 5 * time.Second,
	}
	for _, optionsFunc := range config {
		optionsFunc(&options)
	}
	return &githubApiClient{
		options: options,
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

	client := &http.Client{}
	response, requestErr := client.Do(request)
	if requestErr != nil {
		return nil, requestErr
	}

	defer response.Body.Close()

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
