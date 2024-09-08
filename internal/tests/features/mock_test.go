package features

import (
	"fmt"
	"github.com/matzefriedrich/parsley/pkg/features"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GreeterMock_SayHello(t *testing.T) {

	// Arrange
	mock := NewGreeterMock()
	mock.SayHelloFunc = func(name string) (string, error) {
		return fmt.Sprintf("Hi, %s", name), nil
	}

	const expectedName = "John"

	// Act
	_, _ = mock.SayHello("Max")
	actual, err := mock.SayHello(expectedName)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "Hi, John", actual)

	assert.True(t, mock.Verify(FunctionSayHello, features.TimesOnce(), features.Exact(expectedName)))
	assert.True(t, mock.Verify(FunctionSayHello, features.TimesNever(), features.Exact("Jane")))
	assert.True(t, mock.Verify(FunctionSayHello, features.TimesExactly(2)))
}
