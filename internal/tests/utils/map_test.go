package utils

import (
	"github.com/matzefriedrich/parsley/internal/utils"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func Test_Map_convert_slice(t *testing.T) {
	// Arrange
	source := []string{"1", "2", "3"}
	// Act
	actual := utils.Map(source, func(s string) int {
		value, _ := strconv.Atoi(s)
		return value
	})
	// Assert
	assert.Equal(t, []int{1, 2, 3}, actual)
}
