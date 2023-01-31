package main

func Sum[T int32 | float32 | string](slice []T) T {
	var res T
	for _, element := range slice {
		res += element
	}

	return res
}
