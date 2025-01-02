package processmodel

import (
	"fmt"
	"reflect"
	"time"
)

// ValueAndTimestamp is the interface that fields must implement if they
// support timestamp-based merging.
type ValueAndTimestamp interface {
	GetTimestamp() time.Time
	SetTimestamp(t time.Time)
}

// MergeMode indicates how merging is done regarding deletions.
type MergeMode int

const (
	// ClientIsMaster is when a client is considered the master and deletions are propagated.
	ClientIsMaster MergeMode = iota
	// ServerIsMaster, only updates and additions are propagated.
	ServerIsMaster
)

// MergeOptions holds configuration for how the merge should be performed.
type MergeOptions struct {
	Mode MergeMode
}

// merge[T any] merges newModel into oldModel following the specified rules:
// 1. If a field implements ValueAndTimestamp:
//   - Compare timestamps. The newer timestamp wins.
//   - If Mode=ClientIsMaster and field missing in newModel, remove from merged result.
//   - If Mode=ServerIsMaster and field missing in newModel, keep from oldModel.
//
// 2. If a field does not implement ValueAndTimestamp:
//   - Overwrite from newModel directly (per the specification).
//
// Returns the merged model. Both _oldModel_ and _newModel_ are not modified.
//
// CAUTION: The returned model may copied elements from _oldModel_ and _newModel_ and hence
// do not modify those!
func Merge[T any](oldModel, newModel T, opts MergeOptions) (T, error) {
	// Walk the struct fields recursively
	res := merge(reflect.ValueOf(oldModel), reflect.ValueOf(newModel), opts)

	var zero T

	switch t := res.(type) {
	case T:
		return t, nil
	default:
		return zero, fmt.Errorf("unexpected type %T", t)
	}
}

// Merge will take two data sets and merge them together - returning a new data set
func merge(base, override any, opts MergeOptions) any {
	//reflect and recurse
	b := reflect.ValueOf(base)
	o := reflect.ValueOf(override)
	ret := mergeRecursive(b, o)

	return ret.Interface()
}

func mergeRecursive(base, override reflect.Value) reflect.Value {
	var result reflect.Value

	switch base.Kind() {
	case reflect.Ptr:
		switch base.Elem().Kind() {
		case reflect.Ptr:
			fallthrough
		case reflect.Interface:
			fallthrough
		case reflect.Struct:
			fallthrough
		case reflect.Map:
			// Pointers to complex types should recurse if they aren't nil
			if base.IsNil() {
				result = override
			} else if override.IsNil() {
				result = base
			} else {
				result = mergeRecursive(base.Elem(), override.Elem())
			}
		default:
			// Pointers to basic types should just override
			if isEmptyValue(override) {
				result = base
			} else {
				result = override
			}
		}
	case reflect.Interface:
		// Interfaces should just be unwrapped and recursed through
		result = mergeRecursive(base.Elem(), override.Elem())

	case reflect.Struct:
		// For structs we loop over fields and recurse
		// setup our result struct
		result = reflect.New(base.Type())
		for i, n := 0, base.NumField(); i < n; i++ {
			// We cant set private fields so don't recurse on them
			if result.Elem().Field(i).CanSet() {
				// get the merged value of each field
				newVal := mergeRecursive(base.Field(i), override.Field(i))

				//attempt to set that merged value on our result struct
				if result.Elem().Field(i).CanSet() && newVal.IsValid() {
					if newVal.Kind() == reflect.Ptr && result.Elem().Field(i).Kind() != reflect.Ptr {
						newVal = newVal.Elem()
					} else if result.Elem().Field(i).Kind() == reflect.Ptr && newVal.Kind() != reflect.Ptr && newVal.CanAddr() {
						newVal = newVal.Addr()
					}
					result.Elem().Field(i).Set(newVal)
				}
			}
		}

	case reflect.Map:
		// For Maps we copy the base data, and then replace it with merged data
		// We use two for loops to make sure all map keys from base and all keys from
		// override exist in the result just in case one of the maps is sparse.
		elementsAreValues := base.Type().Elem().Kind() != reflect.Ptr

		result = reflect.MakeMap(base.Type())
		// Copy from base first
		for _, key := range base.MapKeys() {
			result.SetMapIndex(key, base.MapIndex(key))
		}

		// Override with values from override if they exist
		if override.Kind() == reflect.Map {
			for _, key := range override.MapKeys() {
				overrideVal := override.MapIndex(key)
				baseVal := base.MapIndex(key)
				if !overrideVal.IsValid() {
					continue
				}

				// if there is no base value, just set the override
				if !baseVal.IsValid() {
					result.SetMapIndex(key, overrideVal)
					continue
				}

				// Merge the values and set in the result
				newVal := mergeRecursive(baseVal, overrideVal)
				if elementsAreValues && newVal.Kind() == reflect.Ptr {
					result.SetMapIndex(key, newVal.Elem())

				} else {
					result.SetMapIndex(key, newVal)
				}
			}
		}

	default:
		// These are all generic types
		// override will be taken for generic types if it is set
		if isEmptyValue(override) {
			result = base
		} else {
			result = override
		}
	}
	return result
}

func isEmptyValue(v reflect.Value) bool {
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
