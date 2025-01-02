package merge

import "reflect"

// isEmptyValue checks if a reflect.Value is the zero value.
func isEmptyValue(v reflect.Value) bool {
	if !v.IsValid() {
		return true
	}
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.Interface, reflect.Pointer:
		return v.IsZero()
	}
	return false
}

func concatPath(path, name string) string {
	if path == "" {
		return name
	}
	return path + "." + name
}

// getJSONTag get the tag (name part only) from a struct field.
//
// If no _JSON_ field tag is present, the field name is returned.
func getJSONTag(field reflect.StructField) string {
	tag := field.Tag.Get("json")
	if tag == "" {
		return field.Name
	}

	// If the tag is "-", ignore it
	if tag == "-" {
		return ""
	}

	// If the tag contains a comma, only consider the first part
	if idx := 0; idx < len(tag) {
		if tag[idx] == ',' {
			return tag[:idx]
		}
	}

	return tag
}
