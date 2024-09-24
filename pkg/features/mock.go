package features

// MockBase is used as a foundational struct to track and manage mocked functions and their call history.
// It helps in testing by allowing function signature tracking and call verification.
type MockBase struct {
	functions map[string]MockFunction
}

// MockFunction provides a structure to represent a mocked function in test scenarios. It allows tracking its calls and signature.
type MockFunction struct {
	name        string
	signature   string
	tracedCalls []methodCall
}

// String returns the signature of the mocked function if it exists, otherwise it returns the function's name.
func (m MockFunction) String() string {
	if len(m.signature) > 0 {
		return m.signature
	}
	return m.name
}

type methodCall struct {
	args []any
}

// NewMockBase initializes and returns an instance of MockBase, ideal for setting up and using mock functions in tests.
func NewMockBase() MockBase {
	return MockBase{
		functions: make(map[string]MockFunction),
	}
}

// AddFunction adds a new mock function with the specified name and signature to the MockBase instance.
func (m *MockBase) AddFunction(name string, signature string) {
	m.functions[name] = MockFunction{
		name:        name,
		signature:   signature,
		tracedCalls: make([]methodCall, 0),
	}
}

// TraceMethodCall logs the invocation of a mocked function with specified arguments to facilitate function call tracking during testing. Before function calls can be tracked, the function must be registered with the MockBase instance; use AddFunction.
func (m *MockBase) TraceMethodCall(name string, arguments ...any) {
	function, found := m.functions[name]
	if found {
		call := methodCall{
			args: arguments,
		}
		function.tracedCalls = append(function.tracedCalls, call)
		m.functions[name] = function
	}
}
