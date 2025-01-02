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

func isPrimitiveOrPtrToPrimitive(t reflect.Type) bool {
	// If pointer, unwrap one level and check again
	if t.Kind() == reflect.Ptr {
		elem := t.Elem()
		return isPrimitiveKind(elem.Kind())
	}

	return isPrimitiveKind(t.Kind())
}

func isPrimitiveKind(k reflect.Kind) bool {
	switch k {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128,
		reflect.String:
		return true
	default:
		// Other kinds: Struct, Map, Slice, Array, Interface, Ptr, Chan, Func, UnsafePointer are non-primitive
		return false
	}
}

// unwrapToInterface safely unwraps pointer/interface and returns the underlying
// `any` if possible.
func unwrapToInterface(v reflect.Value) any {
	if !v.IsValid() {
		return nil
	}
	// Unwrap interface
	if v.Kind() == reflect.Interface && !v.IsNil() {
		v = v.Elem()
	}
	// Unwrap pointer
	if v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	if !v.IsValid() {
		return nil
	}
	return v.Interface()
}

// tryFieldByName attempts to get a named field from a struct by name (if it exists).
// If overrideVal is not a struct or doesn't have that field, returns zero Value.
func tryFieldByName(structVal reflect.Value, fieldName string) reflect.Value {
	if structVal.Kind() != reflect.Struct {
		return reflect.Value{}
	}
	f, ok := structVal.Type().FieldByName(fieldName)
	if !ok {
		return reflect.Value{}
	}
	return structVal.FieldByIndex(f.Index)
}
