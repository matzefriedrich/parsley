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
	mock.SayHelloFunc = func(name string, polite bool) (string, error) {
		return fmt.Sprintf("Hi, %s", name), nil
	}

	const expectedName = "John"

	// Act
	_, _ = mock.SayHello("Max", false)
	actual, err := mock.SayHello(expectedName, true)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "Hi, John", actual)

	// Verify that John has been greeted, regardless of the politeness flag
	assert.True(t, mock.Verify(FunctionSayHello, features.TimesOnce(), features.Exact(expectedName)))

	// Verify that John has not been greeted non-politely
	assert.True(t, mock.Verify(FunctionSayHello, features.TimesNever(), features.Exact(expectedName), features.Exact(false)))

	// Verify that somebody was greeted politely
	assert.True(t, mock.Verify(FunctionSayHello, features.TimesOnce(), features.IsAny(), features.Exact(true)))

	// Verify that Jane was not greeted
	assert.True(t, mock.Verify(FunctionSayHello, features.TimesNever(), features.Exact("Jane")))

	// Verify that SayHello was called twice in total, regardless of the call parameters
	assert.True(t, mock.Verify(FunctionSayHello, features.TimesExactly(2)))
	assert.True(t, mock.Verify(FunctionSayHello, features.TimesExactly(2), features.IsAny()))
	assert.True(t, mock.Verify(FunctionSayHello, features.TimesExactly(2), features.IsAny(), features.IsAny()))
}
