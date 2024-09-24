package internal

// Stack represents a last-in-first-out (LIFO) collection of items. It can be used to manage a set of elements with push and pop operations.
type Stack[TValue any] []TValue

// MakeStack creates and returns a new Stack initialized with the provided initial values. Useful for algorithm implementations needing LIFO structures.
func MakeStack[TValue any](initialValues ...TValue) Stack[TValue] {
	s := Stack[TValue]{}
	for _, v := range initialValues {
		s.Push(v)
	}
	return s
}

// Any determines if the stack contains at least one element, indicating it is not empty.
func (s *Stack[TValue]) Any() bool {
	return s.IsEmpty() == false
}

// IsEmpty checks whether the stack is empty, which can help determine if there are any elements left to process.
func (s *Stack[TValue]) IsEmpty() bool {
	return len(*s) == 0
}

// Push appends one or more values to the top of the stack, allowing dynamic and flexible addition of elements.
func (s *Stack[TValue]) Push(value ...TValue) {
	*s = append(*s, value...)
}

// Pop removes and returns the item at the top of the stack, allowing for a LIFO retrieval of elements.
func (s *Stack[TValue]) Pop() TValue {
	lastItemIndex := len(*s) - 1
	element := (*s)[lastItemIndex]
	*s = (*s)[:lastItemIndex]
	return element
}
