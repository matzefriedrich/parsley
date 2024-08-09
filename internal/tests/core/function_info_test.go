package core

import (
	"errors"
	"fmt"
	"github.com/matzefriedrich/parsley/internal/core"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
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

func Test_FunctionInfo_ReflectFunctionInfoFrom_function_with_multiple_return_values_returns_error(t *testing.T) {

	// Arrange
	f := func() (some, error) {
		return nil, errors.New("this is an error")
	}

	// Act
	_, err := core.ReflectFunctionInfoFrom(reflect.ValueOf(f))

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, core.ErrReturnTypeHasToHaveExactlyOnReturnValue)
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
