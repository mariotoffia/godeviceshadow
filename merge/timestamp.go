package merge

import (
	"reflect"
	"time"
)

// ValueAndTimestamp is the interface that fields must implement if they
// support timestamp-based merging.
type ValueAndTimestamp interface {
	GetTimestamp() time.Time
	SetTimestamp(t time.Time)
}

// mergeTimestamped handles merging two values that implement ValueAndTimestamp.
func mergeTimestamped(base, override reflect.Value, opts MergeOptions) reflect.Value {
	baseTS, baseOk := unwrapValueAndTimestamp(base)
	overrideTS, overrideOk := unwrapValueAndTimestamp(override)

	// If either is not valid or doesn’t implement the interface, fallback
	if !baseOk && !overrideOk {
		// fallback to normal merging
		return mergeRecursive(base, override, opts)
	}
	if !baseOk {
		// old missing => new is present => use override
		return override
	}
	if !overrideOk {
		// new missing => remove or keep old depending on the mode
		if opts.Mode == ClientIsMaster {
			return reflect.Value{}
		}
		return base
	}

	oldTS := baseTS.GetTimestamp()
	newTS := overrideTS.GetTimestamp()

	switch {
	case newTS.After(oldTS):
		// use override
		return override
	case oldTS.After(newTS):
		// keep old
		return base
	default:
		// equal => no update => keep old
		return base
	}
}

func isValueAndTimestamp(rv reflect.Value) bool {
	// Attempt to unwrap to a ValueAndTimestamp, ignoring if .Addr() can panic
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
		// We can try the interface() directly
		if vt, ok := rv.Interface().(ValueAndTimestamp); ok {
			return vt, true
		}
		// Also unwrap one more level
		rv = rv.Elem()
	}
	if !rv.IsValid() {
		return nil, false
	}

	// If we can directly address it, do so
	if rv.CanAddr() {
		if vt, ok := rv.Addr().Interface().(ValueAndTimestamp); ok {
			return vt, true
		}
	} else {
		// We can't address it => make a new pointer to copy
		newPtr := reflect.New(rv.Type())
		newPtr.Elem().Set(rv)

		v := newPtr.Interface()

		if vt, ok := v.(ValueAndTimestamp); ok {
			return vt, true
		}
	}

	return nil, false
}
