package features

type MockBase struct {
	functions map[string]MockFunction
}

type MockFunction struct {
	name      string
	signature string
	calls     []methodCall
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
		name:      name,
		signature: signature,
		calls:     make([]methodCall, 0),
	}
}

func (m *MockBase) TraceMethodCall(name string, arguments ...any) {
	function, found := m.functions[name]
	if found {
		call := methodCall{
			args: arguments,
		}
		function.calls = append(function.calls, call)
	}
}
