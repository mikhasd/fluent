package array

func Map[I any, O any](in []I, mapper func(I) O) []O {
	out := make([]O, len(in))
	for i, val := range in {
		out[i] = mapper(val)
	}
	return out
}

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
