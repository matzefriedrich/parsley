package utils

import (
	"github.com/matzefriedrich/parsley/internal/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_VersionInfo_IsLessThan_major(t *testing.T) {
	// Arrange
	sut := utils.VersionInfo{Major: 0, Minor: 9, Patch: 3}

	// Act
	actual := sut.LessThan(utils.VersionInfo{Major: 1})

	// Assert
	assert.True(t, actual)
}

func Test_VersionInfo_IsLessThan_major_minor(t *testing.T) {
	// Arrange
	sut := utils.VersionInfo{Major: 0, Minor: 9, Patch: 3}

	// Act
	actual := sut.LessThan(utils.VersionInfo{Major: 0, Minor: 10})

	// Assert
	assert.True(t, actual)
}

func Test_VersionInfo_IsLessThan_major_minor_patch(t *testing.T) {
	// Arrange
	sut := utils.VersionInfo{Major: 0, Minor: 9, Patch: 3}

	// Act
	actual := sut.LessThan(utils.VersionInfo{Major: 0, Minor: 9, Patch: 4})

	// Assert
	assert.True(t, actual)
}

func Test_VersionInfo_Equal(t *testing.T) {
	// Arrange
	sut := utils.VersionInfo{Major: 0, Minor: 9, Patch: 3}

	// Act
	actual := sut.Equal(utils.VersionInfo{Major: 0, Minor: 9, Patch: 3})

	// Assert
	assert.True(t, actual)
}

func Test_ApplicationVersion_has_valid_application_version_string(t *testing.T) {
	// Arrange

	// Act
	version, err := utils.ApplicationVersion()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, version)
}
