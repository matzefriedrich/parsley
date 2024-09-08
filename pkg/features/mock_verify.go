package features

type ArgMatch func(actual any) bool

func IsAny() ArgMatch {
	return func(actual any) bool {
		return true
	}
}

func Exact[T comparable](expected T) ArgMatch {
	return func(actual any) bool {
		value, compatible := actual.(T)
		if compatible && value == expected {
			return true
		}
		return false
	}
}

type TimesFunc func(times int) bool

func TimesOnce() TimesFunc {
	return func(times int) bool {
		return times == 1
	}
}

func TimesAtLeastOnce() TimesFunc {
	return func(times int) bool {
		return times >= 1
	}
}

func TimesExactly(n int) TimesFunc {
	return func(times int) bool {
		return times == n
	}
}

func TimesNever() TimesFunc {
	return func(times int) bool {
		return times == 0
	}
}

func (m *MockBase) Verify(name string, times TimesFunc, matches ...ArgMatch) bool {
	function, found := m.functions[name]
	if found {
		if len(matches) > 0 {
			numMatches := 0
		callsLoop:
			for _, call := range function.tracedCalls {
				for _, arg := range call.args {
					for _, doesMatch := range matches {
						if doesMatch(arg) == false {
							continue callsLoop
						}
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
