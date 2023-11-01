package main

func Sum[T int32 | float32 | string](slice []T) T {
	var res T
	for _, element := range slice {
		res += element
	}

	return res
}

func MapKeys[K comparable, V any](m map[K]V) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	return r
}
