package internal

type Stack[TValue any] []TValue

func MakeStack[TValue any](initialValues ...TValue) Stack[TValue] {
	s := Stack[TValue]{}
	for _, v := range initialValues {
		s.Push(v)
	}
	return s
}

func (s *Stack[TValue]) emptyValue() TValue {
	var value TValue
	return value
}

func (s *Stack[TValue]) Any() bool {
	return s.IsEmpty() == false
}

func (s *Stack[TValue]) IsEmpty() bool {
	return len(*s) == 0
}

func (s *Stack[TValue]) Push(value ...TValue) {
	*s = append(*s, value...)
}

func (s *Stack[TValue]) TryPop() (TValue, bool) {
	if s.IsEmpty() {
		return s.emptyValue(), false
	} else {
		return s.Pop(), true
	}
}

func (s *Stack[TValue]) Pop() TValue {
	lastItemIndex := len(*s) - 1
	element := (*s)[lastItemIndex]
	*s = (*s)[:lastItemIndex]
	return element
}

func (s *Stack[TValue]) Peek() TValue {
	lastItemIndex := len(*s) - 1
	element := (*s)[lastItemIndex]
	return element
}

func (s *Stack[TValue]) TryPeek() (TValue, bool) {
	if s.IsEmpty() {
		return s.emptyValue(), false
	} else {
		return s.Peek(), true
	}
}
