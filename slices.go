package goutils

/* Remove first element of slice */
func DelFirst[T any](s *[]T) (T, bool) {
	if len(*s) == 0 {
		var zero T
		return zero, false
	}
	first := (*s)[0]
	*s = (*s)[1:]
	return first, true
}

/* Remove last element of slice */
func DelLast[T any](s *[]T) (T, bool) {
	if len(*s) == 0 {
		var zero T
		return zero, false
	}
	last := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return last, true
}

func FindIndex[T any](slice []T, pass func(T) bool) int {
	for i, item := range slice {
		if pass(item) {
			return i
		}
	}

	return -1
}

func Find[T any](slice []T, pass func(T, int) bool) T {
	var t T
	for i, item := range slice {
		if pass(item, i) {
			return item
		}
	}

	return t
}

func Some[T any](slice []T, pass func(T, int) bool) bool {
	for i, item := range slice {
		if pass(item, i) {
			return true
		}
	}

	return false
}

func All[T any](slice []T, pass func(T, int) bool) bool {
	for i, item := range slice {
		if !pass(item, i) {
			return false
		}
	}

	return true
}

func Map[T any, K any](slice []T, transform func(T, int) K) []K {
	result := make([]K, len(slice))
	for i, item := range slice {
		result[i] = transform(item, i)
	}
	return result
}

func Filter[T any](slice []T, keep func(T, int) bool) []T {
	var result []T
	for i, item := range slice {
		if keep(item, i) {
			result = append(result, item)
		}
	}
	return result
}

/* Removes falsy elements */
func Compact[T comparable](slice []T) []T {
	var zero T
	var result []T
	for _, item := range slice {
		if item != zero {
			result = append(result, item)
		}
	}
	return result
}

func Reverse[T any](slice []T) []T {
	result := make([]T, len(slice))
	for i, v := range slice {
		result[len(slice)-1-i] = v
	}
	return result
}

func ForEach[T any](slice []T, fn func(item T, index int)) {
	for i, item := range slice {
		fn(item, i)
	}
}
