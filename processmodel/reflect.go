package processmodel

import "reflect"

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
