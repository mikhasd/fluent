package array

// Map creates a new array populated with the result of calling the provided
// `mapper` function on every element of the input array.
func Map[I any, O any](in []I, mapper func(I) O) []O {
	out := make([]O, len(in))
	for i, val := range in {
		out[i] = mapper(val)
	}
	return out
}

// Filter creates a new array with all elements that pass the provided
// `condition` function.
func Filter[T any](in []T, condition func(T) bool) []T {
	if len(in) == 0 {
		return []T{}
	}
	var out []T
	for _, val := range in {
		if condition(val) {
			out = append(out, val)
		}
	}
	return out
}
