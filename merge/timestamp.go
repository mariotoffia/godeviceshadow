package merge

import (
	"reflect"
)

// mergeTimestamped handles merging two values that implement ValueAndTimestamp.
func mergeTimestamped(base, override reflect.Value, opts MergeObject) reflect.Value {
	baseTS, baseOk := unwrapValueAndTimestamp(base)
	overrideTS, overrideOk := unwrapValueAndTimestamp(override)

	if !baseOk && !overrideOk {
		// fallback to normal merging
		return mergeRecursive(base, override, opts)
	}

	if !baseOk {
		// old missing -> new is present -> use override

		return override
	}

	if !overrideOk {
		// new missing
		if opts.Mode == ClientIsMaster {
			return reflect.Value{} // remove
		}

		return base // keep old
	}

	oldTS := baseTS.GetTimestamp()
	newTS := overrideTS.GetTimestamp()

	switch {
	case newTS.After(oldTS):
		return override
	case oldTS.After(newTS):
		return base // keep old
	default:
		// equal => no update => keep old
		return base
	}
}

func isValueAndTimestamp(rv reflect.Value) bool {
	_, ok := unwrapValueAndTimestamp(rv)
	return ok
}

func unwrapValueAndTimestamp(rv reflect.Value) (ValueAndTimestamp, bool) {
	if !rv.IsValid() {
		return nil, false
	}

	// If rv is interface or pointer, unwrap once
	if rv.Kind() == reflect.Interface && !rv.IsNil() {
		rv = rv.Elem()
	}
	if rv.Kind() == reflect.Ptr && !rv.IsNil() {
		// Since Ptr
		if vt, ok := rv.Interface().(ValueAndTimestamp); ok {
			return vt, true
		}
		// unwrap one more level
		rv = rv.Elem()
	}
	if !rv.IsValid() {
		return nil, false
	}

	// can directly address it? -> do so
	if rv.CanAddr() {
		if vt, ok := rv.Addr().Interface().(ValueAndTimestamp); ok {
			return vt, true
		}
	} else {
		// can't address it -> make a new pointer to copy
		newPtr := reflect.New(rv.Type())
		newPtr.Elem().Set(rv)

		v := newPtr.Interface()

		if vt, ok := v.(ValueAndTimestamp); ok {
			return vt, true
		}
	}

	return nil, false
}
