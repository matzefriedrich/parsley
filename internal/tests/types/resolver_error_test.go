package types_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/matzefriedrich/parsley/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestResolverError_ErrorsIs(t *testing.T) {
	// Arrange
	err := types.NewResolverError(types.ErrorServiceTypeNotRegistered, types.ForServiceTypeByName("MyService"))

	// Act & Assert
	assert.True(t, errors.Is(err, types.ErrServiceTypeNotRegistered))
}

func TestResolverError_Format(t *testing.T) {
	// Arrange
	err := types.NewResolverError(types.ErrorServiceTypeNotRegistered, types.ForServiceTypeByName("MyService"))

	// Act
	msg := fmt.Sprintf("%v", err)

	// Assert
	assert.Equal(t, "service type is not registered: MyService", msg)
}

func TestResolverError_Getter(t *testing.T) {
	// Arrange
	err := types.NewResolverError(types.ErrorServiceTypeNotRegistered, types.ForServiceTypeByName("MyService"))

	// Act
	withType, ok := err.(types.ParsleyErrorWithServiceTypeName)

	// Assert
	assert.True(t, ok)
	assert.Equal(t, "MyService", withType.ServiceTypeName())
}
