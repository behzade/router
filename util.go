package router

func insertToIndex[T any](array []T, index int, value T) []T {
	if len(array) == index { // nil or empty slice or after last element
		return append(array, value)
	}
	array = append(array[:index+1], array[index:]...) // index < len(a)
	array[index] = value
	return array
}

func keys[T comparable, R any](m map[T]R) []T {
	keys := make([]T, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}
