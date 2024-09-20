package utils

// Map applies a given function to each element of a slice and returns a new slice containing the results of the function.
func Map[S ~[]E, E any, V any](ts S, fn func(E) V) []V {
	result := make([]V, len(ts))
	for i, t := range ts {
		result[i] = fn(t)
	}
	return result
}
