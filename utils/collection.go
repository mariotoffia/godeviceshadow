package utils

// Find will find the first element in the slice that matches the predicate.
// If found it will return the element and `true`, otherwise the zero value of `T` and `false`.
func Find[T any](s []T, f func(T) bool) (T, bool) {
	for _, v := range s {
		if f(v) {
			return v, true
		}
	}

	var zero T

	return zero, false
}

// FindPtr is same as `Find` but will use pointer to minimize copying.
func FindPtr[T any](s []T, f func(*T) bool) (*T, bool) {
	for i := range s {
		v := &s[i]

		if f(v) {
			return v, true
		}
	}

	return nil, false
}
