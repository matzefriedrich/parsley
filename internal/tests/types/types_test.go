package types

import (
	"github.com/matzefriedrich/parsley/pkg/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ServiceKey_value_matches_empty_string(t *testing.T) {
	// Arrange
	sut := types.NewServiceKey("key")

	// Act
	actual := sut.String()

	// Assert
	assert.Equal(t, "key", actual)
}
