package core

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/matzefriedrich/parsley/internal/core"
	"github.com/stretchr/testify/assert"
)

type expectedFunctionInfo struct {
	functionName             string
	returnTypeName           string
	numParameters            int
	formattedSignatureString string
}

func Test_FunctionInfo_ReflectFunctionInfoFrom_variable_returns_error(t *testing.T) {

	// Arrange
	s := "i am not function"

	// Act
	_, err := core.ReflectFunctionInfoFrom(reflect.ValueOf(s))

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, core.ErrNotAFunction)
}

func Test_FunctionInfo_ReflectFunctionInfoFrom_function_with_three_return_values_returns_error(t *testing.T) {

	// Arrange
	f := func() (some, error, error) {
		return nil, nil, nil
	}

	// Act
	_, err := core.ReflectFunctionInfoFrom(reflect.ValueOf(f))

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, core.ErrReturnTypeHasToHaveExactlyOnReturnValue)
}

func Test_FunctionInfo_ReflectFunctionInfoFrom_function_with_two_return_values_including_error_succeeds(t *testing.T) {

	// Arrange
	f := func() (some, error) {
		return nil, nil
	}

	// Act
	info, err := core.ReflectFunctionInfoFrom(reflect.ValueOf(f))

	// Assert
	assert.NoError(t, err)
	assert.True(t, info.HasErrorReturn())
}

func Test_FunctionInfo_ReflectFunctionInfoFrom_local_function_returning_an_interface(t *testing.T) {

	// Arrange
	_ = func() {}                                      // this will become "func1"
	localFunctionReturningAnInterface := func() some { // and this "func2"; there is no way to get a better name
		return nil
	}

	const expectedFunctionPackageName = "github.com/matzefriedrich/parsley/internal/tests/core.Test_FunctionInfo_ReflectFunctionInfoFrom_local_function_returning_an_interface"
	expectedAnonymousFunctionName := fmt.Sprintf("%s.func2", expectedFunctionPackageName)
	expected := expectedFunctionInfo{
		functionName:             expectedAnonymousFunctionName,
		formattedSignatureString: fmt.Sprintf("%s() core.some", expectedAnonymousFunctionName),
		returnTypeName:           "some",
		numParameters:            0,
	}

	// Act, Assert
	reflectFunctionInfoFromTestHelper(t, localFunctionReturningAnInterface, expected)
}

func Test_FunctionInfo_ReflectFunctionInfoFrom_named_function_returning_an_interface(t *testing.T) {

	// Arrange
	const expectedFunctionPackageName = "github.com/matzefriedrich/parsley/internal/tests/core"
	expectedFunctionName := fmt.Sprintf("%s.functionReturningAnInterface", expectedFunctionPackageName)
	expected := expectedFunctionInfo{
		functionName:             expectedFunctionName,
		formattedSignatureString: fmt.Sprintf("%s() core.some", expectedFunctionName),
		returnTypeName:           "some",
		numParameters:            0,
	}

	// Act, Assert
	reflectFunctionInfoFromTestHelper(t, functionReturningAnInterface, expected)
}

func Test_FunctionInfo_ReflectFunctionInfoFrom_named_function_with_parameters_returning_an_interface(t *testing.T) {

	// Arrange
	const expectedFunctionPackageName = "github.com/matzefriedrich/parsley/internal/tests/core"
	expectedFunctionName := fmt.Sprintf("%s.functionWithParametersReturningAnInterface", expectedFunctionPackageName)
	expected := expectedFunctionInfo{
		functionName:             expectedFunctionName,
		formattedSignatureString: fmt.Sprintf("%s(interface {}) core.some", expectedFunctionName),
		returnTypeName:           "some",
		numParameters:            1,
	}

	// Act, Assert
	reflectFunctionInfoFromTestHelper(t, functionWithParametersReturningAnInterface, expected)
}

func Test_ReflectFunctionInfoFrom_verify_activator_function_parameters_and_return_values(t *testing.T) {

	t.Run("Standard function", func(t *testing.T) {
		// Arrange
		f := func(a dummy, b dummy) dummy { return dummy{} }

		// Act
		info, err := core.ReflectFunctionInfoFrom(reflect.ValueOf(f))

		// Assert
		assert.NoError(t, err)
		assert.False(t, info.HasErrorReturn())
		assert.False(t, info.HasContextParameter())
		assert.Equal(t, "dummy", info.ReturnType().Name())
		assert.Len(t, info.Parameters(), 2)
		assert.Contains(t, info.String(), "(core.dummy,core.dummy) core.dummy")
	})

	t.Run("Function with error return", func(t *testing.T) {
		// Arrange
		f := func() (dummy, error) { return dummy{}, nil }

		// Act
		info, err := core.ReflectFunctionInfoFrom(reflect.ValueOf(f))

		// Assert
		assert.NoError(t, err)
		assert.True(t, info.HasErrorReturn())
		assert.Equal(t, "dummy", info.ReturnType().Name())
		assert.Contains(t, info.String(), "() core.dummy")
	})

	t.Run("Function with context parameter", func(t *testing.T) {
		// Arrange
		f := func(ctx context.Context, s dummy) dummy { return dummy{} }

		// Act
		info, err := core.ReflectFunctionInfoFrom(reflect.ValueOf(f))

		// Assert
		assert.NoError(t, err)
		assert.True(t, info.HasContextParameter())
		assert.Len(t, info.Parameters(), 2)
		assert.Contains(t, info.String(), "(context.Context,core.dummy) core.dummy")
	})

	t.Run("Function with context and error", func(t *testing.T) {
		// Arrange
		f := func(ctx context.Context) (dummy, error) { return dummy{}, nil }

		// Act
		info, err := core.ReflectFunctionInfoFrom(reflect.ValueOf(f))

		// Arrange
		assert.NoError(t, err)
		assert.True(t, info.HasContextParameter())
		assert.True(t, info.HasErrorReturn())
		assert.Contains(t, info.String(), "(context.Context) core.dummy")
	})
}

func functionReturningAnInterface() some {
	return nil
}

func functionWithParametersReturningAnInterface(dependency any) some {
	return nil
}

func reflectFunctionInfoFromTestHelper(t *testing.T, target any, expected expectedFunctionInfo) {

	// Arrange
	f := reflect.ValueOf(target)

	// Act
	actual, err := core.ReflectFunctionInfoFrom(f)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, actual)

	assert.Equal(t, expected.functionName, actual.Name())

	returnType := actual.ReturnType()
	assert.NotNil(t, returnType)
	assert.Equal(t, expected.returnTypeName, returnType.Name())

	parameterTypes := actual.Parameters()
	assert.NotNil(t, parameterTypes)
	assert.Equal(t, expected.numParameters, len(parameterTypes))

	assert.Equal(t, expected.formattedSignatureString, actual.String())
}

type some interface {
}

type dummy struct{}
