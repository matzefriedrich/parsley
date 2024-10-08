package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_VersionInfo_tryParseVersionInfo_valid_semantic_version_expression_with_prefix_does_not_return_err(t *testing.T) {

	// Arrange
	const version = "v0.1.0"

	// Act
	actual, err := tryParseVersionInfo(version)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, "0.1.0", actual.String())
}

func Test_VersionInfo_tryParseVersionInfo_valid_semantic_version_expression_no_prefix_does_not_return_err(t *testing.T) {

	// Arrange
	const version = "0.1.0"

	// Act
	actual, err := tryParseVersionInfo(version)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, "0.1.0", actual.String())
}

func Test_VersionInfo_tryParseVersionInfo_invalid_semantic_version_expression_returns_err(t *testing.T) {

	// Arrange
	const version = "a.b.c"

	// Act
	actual, err := tryParseVersionInfo(version)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, actual)
}

func Test_VersionInfo_LessThan(t *testing.T) {

	// Arrange
	const a = "v0.1.0"
	const b = "v0.2.0"

	version, _ := tryParseVersionInfo(a)
	otherVersion, _ := tryParseVersionInfo(b)

	// Act
	actual := version.LessThan(*otherVersion)

	// Assert
	assert.True(t, actual)
}

func Test_VersionInfo_LessThanOrEqual(t *testing.T) {

	// Arrange
	const a = "v0.1.0"
	b := []string{"v0.2.0", "v0.1.0"}

	version, _ := tryParseVersionInfo(a)

	// Act
	for _, v := range b {
		otherVersion, _ := tryParseVersionInfo(v)

		actual := version.LessThanOrEqual(*otherVersion)

		// Assert
		assert.True(t, actual)
	}
}

func Test_VersionInfo_Equal(t *testing.T) {

	// Arrange
	const a = "v0.1.0"
	const b = "v0.1.0"

	version, _ := tryParseVersionInfo(a)
	otherVersion, _ := tryParseVersionInfo(b)

	// Act

	actual := version.Equal(*otherVersion)

	// Assert
	assert.True(t, actual)
}
