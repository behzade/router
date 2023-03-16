package router

func insertToIndex[T any](array []T, index int, value T) []T {
    if len(array) == index { // nil or empty slice or after last element
        return append(array, value)
    }
    array = append(array[:index+1], array[index:]...) // index < len(a)
    array[index] = value
    return array
}
