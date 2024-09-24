package features

// ArgMatch is a function type used to match an argument against a certain condition during mock function verification.
type ArgMatch func(actual any) bool

// IsAny always returns true, enabling it to match any given argument during mock function verification.
func IsAny() ArgMatch {
	return func(actual any) bool {
		return true
	}
}

// Exact returns an ArgMatch that checks if a given argument is exactly equal to the specified expected value.
func Exact[T comparable](expected T) ArgMatch {
	return func(actual any) bool {
		value, compatible := actual.(T)
		if compatible && value == expected {
			return true
		}
		return false
	}
}

// TimesFunc is used to verify the number of times a mock function is called. It allows flexibility in call count assertions.
type TimesFunc func(times int) bool

// TimesOnce returns a TimesFunc that checks if the number of function calls equals one. It is useful for verifying single call assertions.
func TimesOnce() TimesFunc {
	return func(times int) bool {
		return times == 1
	}
}

// TimesAtLeastOnce returns a TimesFunc that verifies if a mock function is called at least once.
func TimesAtLeastOnce() TimesFunc {
	return func(times int) bool {
		return times >= 1
	}
}

// TimesExactly returns a TimesFunc that checks if the number of function calls is exactly equal to the specified value.
func TimesExactly(n int) TimesFunc {
	return func(times int) bool {
		return times == n
	}
}

// TimesNever returns a TimesFunc that ensures the function has never been called, providing a strict zero call condition.
func TimesNever() TimesFunc {
	return func(times int) bool {
		return times == 0
	}
}

// Verify checks if a mock function was called a specific number of times, optionally matching provided argument conditions.
func (m *MockBase) Verify(name string, times TimesFunc, matches ...ArgMatch) bool {
	function, found := m.functions[name]
	if found {
		if len(matches) > 0 {
			numMatches := 0
		callsLoop:
			for _, call := range function.tracedCalls {
				for i, arg := range call.args {
					if i < len(matches) {
						match := matches[i]
						if match(arg) == false {
							continue callsLoop
						}
					} else {
						break
					}
				}
				numMatches++
			}
			return times(numMatches)
		} else {
			numCalls := len(function.tracedCalls)
			return times(numCalls)
		}
	}
	return false
}
