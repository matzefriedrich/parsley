package core

import (
	"github.com/matzefriedrich/parsley/internal/core"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ServiceIdSequence_Next_produces_incr_values(t *testing.T) {

	// Arrange
	sut := core.NewServiceId(0)
	collectedIdentifiers := make([]uint64, 0)

	// Act
	for i := 0; i < 10; i++ {
		id := sut.Next()
		collectedIdentifiers = append(collectedIdentifiers, id)
	}

	// Assert
	expectedIdentifiers := []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	assert.Equal(t, expectedIdentifiers, collectedIdentifiers)
}
