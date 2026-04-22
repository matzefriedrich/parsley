package features

import (
	"errors"
	"reflect"
	"testing"

	"github.com/matzefriedrich/parsley/pkg/features"
	"github.com/stretchr/testify/assert"
)

func Test_NewInterceptorBase_verify_name_and_position_accessors(t *testing.T) {
	// Arrange
	name := "test-interceptor"
	position := 10

	// Act
	base := features.NewInterceptorBase(name, position)

	// Assert
	assert.Equal(t, name, base.Name())
	assert.Equal(t, position, base.Position())
}

func Test_NewReturnValueInfo_verify_accessors(t *testing.T) {
	// Arrange
	name := "result0"
	value := "world"
	valueType := reflect.TypeOf(value)

	// Act
	sut := features.NewReturnValueInfo(name, value, valueType)

	// Assert
	assert.Equal(t, name, sut.Name())
	assert.Equal(t, value, sut.Value())
	assert.Equal(t, valueType, sut.ValueType())

	s := sut.String()
	assert.Contains(t, s, name)
	assert.Contains(t, s, valueType.String())
	assert.Contains(t, s, value)
}

func Test_ReturnValueInfo_nil_value_does_not_panic(t *testing.T) {
	// Arrange
	name := "result0"
	var value interface{} = nil
	errorType := reflect.TypeOf((*error)(nil)).Elem()

	// Act
	actual := features.NewReturnValueInfo(name, value, errorType)

	// Assert
	assert.Equal(t, name, actual.Name())
	assert.Nil(t, actual.Value())
	assert.Equal(t, errorType, actual.ValueType())

	s := actual.String()
	assert.Contains(t, s, name)
	assert.Contains(t, s, errorType.String())
	assert.Contains(t, s, "nil")
}

type captureInterceptor struct {
	features.InterceptorBase
	capturedParameters   []features.ParameterInfo
	capturedReturnValues []features.ReturnValueInfo
	capturedError        error
}

func (c *captureInterceptor) Enter(_ any, _ string, parameters []features.ParameterInfo) {
	c.capturedParameters = parameters
}

func (c *captureInterceptor) Exit(_ any, _ string, returnValues []features.ReturnValueInfo) {
	c.capturedReturnValues = returnValues
}

func (c *captureInterceptor) OnError(_ any, _ string, err error) {
	c.capturedError = err
}

func newCaptureInterceptor() *captureInterceptor {
	return &captureInterceptor{
		InterceptorBase: features.NewInterceptorBase("capture", 0),
	}
}

func Test_ProxyBase_Interception(t *testing.T) {
	// Arrange
	target := &greeter{}
	interceptor := newCaptureInterceptor()
	sut := features.NewProxyBase[Greeter](target, []features.MethodInterceptor{interceptor})

	methodName := "SayHello"
	parameterNames := []string{"name", "polite"}
	parameters := map[string]interface{}{
		"name":   "John",
		"polite": true,
	}
	resultNames := []string{"result0", "result1"}

	callContext := features.NewMethodCallContext(methodName, parameterNames, parameters, resultNames...)

	// Act
	sut.InvokeEnterMethodInterceptors(callContext)

	// Note: Simulate target call
	res0, res1 := target.SayHello("John", true)

	sut.InvokeMethodErrorInterceptors(callContext, res0, res1)
	sut.InvokeExitMethodInterceptors(callContext)

	// Assert
	assert.Len(t, interceptor.capturedParameters, 2)

	firstCapturedParameter := interceptor.capturedParameters[0]
	assert.Equal(t, "name", firstCapturedParameter.Name())
	assert.Equal(t, "John", firstCapturedParameter.Value())

	secondCapturedParameter := interceptor.capturedParameters[1]
	assert.Equal(t, "polite", secondCapturedParameter.Name())
	assert.Equal(t, true, secondCapturedParameter.Value())

	s := firstCapturedParameter.String()
	assert.Contains(t, s, "name")
	assert.Contains(t, s, "string")
	assert.Contains(t, s, "John")

	assert.Len(t, interceptor.capturedReturnValues, 2)

	firstReturnValue := interceptor.capturedReturnValues[0]
	assert.Equal(t, "result0", firstReturnValue.Name())
	assert.Equal(t, "Hello John", firstReturnValue.Value())

	secondReturnValue := interceptor.capturedReturnValues[1]
	assert.Equal(t, "result1", secondReturnValue.Name())
	assert.Nil(t, secondReturnValue.Value())
}

func Test_ProxyBase_InvokeMethodErrorInterceptors_invokes_OnError(t *testing.T) {
	// Arrange
	target := &johnGreeter{} // From register_proxy_error_test.go
	interceptor := newCaptureInterceptor()
	sut := features.NewProxyBase[Greeter](target, []features.MethodInterceptor{interceptor})

	methodName := "SayHello"
	parameterNames := []string{"name", "polite"}
	parameters := map[string]interface{}{
		"name":   "Jane",
		"polite": false,
	}
	resultNames := []string{"result0", "result1"}

	callContext := features.NewMethodCallContext(methodName, parameterNames, parameters, resultNames...)

	// Act
	res0, res1 := target.SayHello("Jane", false) // Note: simulate target call
	sut.InvokeMethodErrorInterceptors(callContext, res0, res1)

	// Assert
	assert.NotNil(t, interceptor.capturedError)
	assert.Equal(t, "name is not John", interceptor.capturedError.Error())

	// Test proxyError behavior (Unwrap and Is)
	unwrapped := errors.Unwrap(interceptor.capturedError)
	assert.NotNil(t, unwrapped)
	assert.Equal(t, "name is not John", unwrapped.Error())

	assert.True(t, errors.Is(interceptor.capturedError, errors.New("name is not John")))
}

func Test_ProxyBase_InvokeEnterMethodInterceptors_reflects_type_for_nil_parameter_value(t *testing.T) {
	// Arrange
	target := &nilParamGreeter{}
	interceptor := newCaptureInterceptor()
	sut := features.NewProxyBase[NilParamRepro](target, []features.MethodInterceptor{interceptor})

	methodName := "SaySomething"
	parameterNames := []string{"err"}
	parameters := map[string]interface{}{
		"err": nil,
	}

	callContext := features.NewMethodCallContext(methodName, parameterNames, parameters)

	// Act
	sut.InvokeEnterMethodInterceptors(callContext)

	// Assert
	assert.Len(t, interceptor.capturedParameters, 1)

	firstCapturedParameter := interceptor.capturedParameters[0]
	assert.Equal(t, "err", firstCapturedParameter.Name())
	assert.Nil(t, firstCapturedParameter.Value())

	expectedParameterType := reflect.TypeOf((*error)(nil)).Elem()
	assert.Equal(t, expectedParameterType, firstCapturedParameter.ParameterType())
}

type nilParamGreeter struct{}

func (n *nilParamGreeter) SaySomething(_ error) {}

func Test_ProxyBase_InvokeEnterMethodInterceptors_invokes_interceptors_in_correct_order(t *testing.T) {
	// Arrange
	target := &greeter{}
	order := make([]string, 0)
	i1 := &orderInterceptor{InterceptorBase: features.NewInterceptorBase("i1", 10), order: &order}
	i2 := &orderInterceptor{InterceptorBase: features.NewInterceptorBase("i2", 5), order: &order}
	i3 := &orderInterceptor{InterceptorBase: features.NewInterceptorBase("i3", 15), order: &order}

	sut := features.NewProxyBase[Greeter](target, []features.MethodInterceptor{i1, i2, i3})
	callContext := features.NewMethodCallContext("SayNothing", []string{}, map[string]interface{}{})

	// Act
	sut.InvokeEnterMethodInterceptors(callContext)

	// Assert
	assert.Equal(t, []string{"i2", "i1", "i3"}, order)
}

type orderInterceptor struct {
	features.InterceptorBase
	order *[]string
}

func (o *orderInterceptor) Enter(_ any, _ string, _ []features.ParameterInfo) {
	*o.order = append(*o.order, o.Name())
}

func (o *orderInterceptor) Exit(_ any, _ string, _ []features.ReturnValueInfo) {}
func (o *orderInterceptor) OnError(_ any, _ string, _ error)                   {}

func Test_ProxyBase_InvokeEnterMethodInterceptors_does_not_panic_if_method_not_found(t *testing.T) {
	// Arrange
	target := &greeter{}
	interceptor := newCaptureInterceptor()
	sut := features.NewProxyBase[Greeter](target, []features.MethodInterceptor{interceptor})

	callContext := features.NewMethodCallContext("NonExistentMethod", []string{}, map[string]interface{}{})

	// Act
	sut.InvokeEnterMethodInterceptors(callContext)
	sut.InvokeMethodErrorInterceptors(callContext, "some result")

	// Assert
	assert.Empty(t, interceptor.capturedParameters)
}

func Test_ProxyBase_NewMethodCallContext_uses_default_return_parameter_names_if_not_given(t *testing.T) {
	// Arrange
	target := &greeter{}
	interceptor := newCaptureInterceptor()
	sut := features.NewProxyBase[Greeter](target, []features.MethodInterceptor{interceptor})

	methodName := "SayHello"
	parameterNames := []string{"name", "polite"}
	parameters := map[string]interface{}{"name": "John", "polite": true}

	// Note: create context with default result names
	callContext := features.NewMethodCallContext(methodName, parameterNames, parameters)

	// Act
	sut.InvokeMethodErrorInterceptors(callContext, "Hello", nil)
	sut.InvokeExitMethodInterceptors(callContext)

	// Assert
	assert.Len(t, interceptor.capturedReturnValues, 2)

	firstCapturedReturnValue := interceptor.capturedReturnValues[0]
	assert.Equal(t, "result0", firstCapturedReturnValue.Name())

	secondCapturedReturnValue := interceptor.capturedReturnValues[1]
	assert.Equal(t, "result1", secondCapturedReturnValue.Name())
}
