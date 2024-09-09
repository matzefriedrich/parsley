package core

import (
	"github.com/matzefriedrich/parsley/internal/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Version_parse_version_from_github_release(t *testing.T) {

	// Arrange
	const version = "1.2.3"
	release := utils.GithubRelease{TagName: version}

	const expectedVersionString = "1.2.3"

	// Act
	actual, err := release.TryParseVersionFromTag()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedVersionString, actual.String())
}

func Test_Version_parse_prefixed_version_from_github_release(t *testing.T) {

	// Arrange
	const version = "v1.2.3"
	release := utils.GithubRelease{TagName: version}

	const expectedVersionString = "1.2.3"

	// Act
	actual, err := release.TryParseVersionFromTag()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedVersionString, actual.String())
}

func Test_Version_parse_prefixed_prerelease_version_from_github_release(t *testing.T) {

	// Arrange
	const version = "v1.2.3-alpha.1"
	release := utils.GithubRelease{TagName: version}

	const expectedVersionString = "1.2.3"

	// Act
	actual, err := release.TryParseVersionFromTag()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedVersionString, actual.String())
}
