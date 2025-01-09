package persistutils

// FromConfig will return the value of the key from the config map if it exists.
//
// It must also match the type of T. If not found or type mismatch it will return the zero value of T and false.
func FromConfig[T any](config map[string]any, key string) (T, bool) {
	if config != nil {
		if v, ok := config[key]; ok {
			if vt, ok := v.(T); ok {
				return vt, true
			}
		}
	}

	var zero T

	return zero, false
}
