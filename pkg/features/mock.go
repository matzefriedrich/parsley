package features

type MockBase struct {
	functions map[string]MockFunction
}

type MockFunction struct {
	name        string
	signature   string
	tracedCalls []methodCall
}

type methodCall struct {
	args []any
}

func NewMockBase() MockBase {
	return MockBase{
		functions: make(map[string]MockFunction),
	}
}

func (m *MockBase) AddFunction(name string, signature string) {
	m.functions[name] = MockFunction{
		name:        name,
		signature:   signature,
		tracedCalls: make([]methodCall, 0),
	}
}

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

