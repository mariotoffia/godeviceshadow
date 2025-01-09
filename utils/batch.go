package utils

// ToBatch splits a slice into batches of the given size.
func ToBatch[T any](s []T, batchSize int) [][]T {
	var batches [][]T

	for batchSize < len(s) {
		s, batches = s[batchSize:], append(batches, s[0:batchSize:batchSize])
	}

	return append(batches, s)
}
