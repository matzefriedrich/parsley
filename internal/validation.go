package internal

import "reflect"

func IsNil[T any](instance T) bool {
	// Use reflection to check if instance is nil
	val := reflect.ValueOf(instance)
	switch val.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return val.IsNil()
	case reflect.Invalid:
		return true
	default:
		return false
	}
}
