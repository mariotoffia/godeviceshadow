package merge

import (
	"fmt"
	"reflect"
)

// MergeMode indicates how merging is done regarding deletions.
type MergeMode int

const (
	// ClientIsMaster is when a client is considered the master
	// and deletions are propagated.
	ClientIsMaster MergeMode = 0
	// ServerIsMaster, only updates and additions are propagated.
	ServerIsMaster = 1
)

// MergeOptions holds configuration for how the merge should be performed.
type MergeOptions struct {
	Mode MergeMode
	// Loggers will be notified on add, updated, remove, not-changed operations while merging.
	Loggers MergeLoggers
}

type MergeObject struct {
	MergeOptions
	CurrentPath string
}

// Merge merges newModel into oldModel following the specified rules:
//
//  1. If a field implements ValueAndTimestamp:
//     - Compare timestamps. The newer timestamp wins.
//     - If Mode=ClientIsMaster and field missing in newModel, remove from merged result.
//     - If Mode=ServerIsMaster and field missing in newModel, keep from oldModel.
//     - If timestamps are equal => no update (keep old).
//
//  2. If a field does not implement ValueAndTimestamp:
//     - Overwrite from newModel if present.
//     - If absent in newModel: remove if ClientIsMaster, keep if ServerIsMaster.
//
// Returns the merged model. Neither _oldModel_ nor _newModel_ is modified.
func Merge[T any](oldModel, newModel T, opts MergeOptions) (T, error) {
	//
	oldVal := reflect.ValueOf(oldModel)
	newVal := reflect.ValueOf(newModel)

	mergedVal := mergeRecursive(oldVal, newVal, MergeObject{MergeOptions: opts})

	var zero T
	res, ok := mergedVal.Interface().(T)

	if !ok {
		return zero, fmt.Errorf("unexpected type %T after merge", mergedVal.Interface())
	}

	return res, nil
}

// mergeRecursive is the core that merges base with override respecting MergeOptions.
func mergeRecursive(base, override reflect.Value, opts MergeObject) reflect.Value {
	// Handle invalid (e.g., the new model missing that field).
	if !base.IsValid() {
		// base is missing => just return override as final (unless it's invalid too)
		return override
	}

	if !override.IsValid() {
		// override is missing => in ClientIsMaster remove, in ServerIsMaster keep
		switch opts.Mode {
		case ClientIsMaster:
			// Return zero or invalid. For struct fields, we'll skip it in the parent.
			return reflect.Value{}
		case ServerIsMaster:
			// Keep old
			return base
		}
	}

	// If both are interface or pointer, handle unwrapping
	baseVal := base
	overrideVal := override

	// If either is an interface, unwrap.
	if base.Kind() == reflect.Interface {
		baseVal = base.Elem()
	}

	if override.Kind() == reflect.Interface {
		overrideVal = override.Elem()
	}

	// is pointer -> unwrap to the element (unless nil)
	if baseVal.Kind() == reflect.Ptr && !baseVal.IsNil() {
		baseVal = baseVal.Elem()
	}

	if overrideVal.Kind() == reflect.Ptr && !overrideVal.IsNil() {
		overrideVal = overrideVal.Elem()
	}

	// implements ValueAndTimestamp -> handle timestamped merge
	if isValueAndTimestamp(baseVal) && isValueAndTimestamp(overrideVal) {
		return mergeTimestamped(base, override, opts)
	}

	switch baseVal.Kind() {
	case reflect.Struct:
		return mergeStruct(baseVal, overrideVal, opts)

	case reflect.Map:
		return mergeMap(baseVal, overrideVal, opts)

	case reflect.Slice, reflect.Array:
		return mergeSlice(baseVal, overrideVal, opts)

	// Basic types
	default:
		// override is empty -> remove or keep
		if isEmptyValue(overrideVal) {
			if opts.Mode == ClientIsMaster {
				// remove
				return reflect.Value{}
			}
			// keep old
			return base
		}

		// override is non-empty -> override
		return override
	}
}

func mergeSlice(baseVal, overrideVal reflect.Value, opts MergeObject) reflect.Value {
	// base or invalid or nil -> override wins
	if !baseVal.IsValid() || baseVal.IsNil() {
		return overrideVal
	}

	// override is invalid or nil -> remove or keep
	if !overrideVal.IsValid() || overrideVal.IsNil() {
		switch opts.Mode {
		case ClientIsMaster:
			// remove entirely
			return reflect.Value{}
		case ServerIsMaster:
			// keep old
			return baseVal
		}
	}

	// Both are valid slices
	baseLen := baseVal.Len()
	ovLen := overrideVal.Len()

	// Merge each element up to the min of both lengths
	minLen := baseLen
	maxLen := ovLen
	if ovLen < minLen {
		minLen = ovLen
		maxLen = baseLen
	}

	// Create a new slice of the same type as baseVal
	result := reflect.MakeSlice(baseVal.Type(), 0, maxLen)

	for i := 0; i < minLen; i++ {
		baseElem := baseVal.Index(i)
		ovElem := overrideVal.Index(i)
		mergedElem := mergeRecursive(baseElem, ovElem, opts)

		if mergedElem.IsValid() {
			// "replace"
			result = reflect.Append(result, mergedElem)
		} else {
			if opts.Mode == ServerIsMaster {
				// ServerIsMaster -> keep old
				result = reflect.Append(result, baseElem)
			}
		}
	}

	// new slice is longer -> add extra elements in override
	if ovLen > minLen {
		for i := minLen; i < ovLen; i++ {
			result = reflect.Append(result, overrideVal.Index(i))
		}
	}

	// old slice is longer -> remove or keep
	if baseLen > minLen {
		// If ClientIsMaster -> (effectively) remove
		if opts.Mode == ServerIsMaster {
			// ServerIsMaster -> keep
			for i := minLen; i < baseLen; i++ {
				result = reflect.Append(result, baseVal.Index(i))
			}
		}
	}

	return result
}

// mergeStruct merges two struct values (non-timestamped case).
func mergeStruct(baseVal, overrideVal reflect.Value, opts MergeObject) reflect.Value {
	if !baseVal.IsValid() {
		return overrideVal
	}

	if !overrideVal.IsValid() {
		switch opts.Mode {
		case ClientIsMaster:
			return reflect.Value{}
		case ServerIsMaster:
			return baseVal
		}
	}

	result := reflect.New(baseVal.Type()).Elem() // same struct type

	nFields := baseVal.NumField()

	for i := 0; i < nFields; i++ {
		baseFieldValue := baseVal.Field(i)
		baseField := baseVal.Type().Field(i)

		opts.CurrentPath = concatPath(opts.CurrentPath, getJSONTag(baseField))

		var overrideField reflect.Value

		// get the override field if it exists
		if i < overrideVal.NumField() && baseField.Name == overrideVal.Type().Field(i).Name {
			overrideField = overrideVal.Field(i)
		} else {
			// if struct layouts differ -> skip it
			continue
		}

		if !result.Field(i).CanSet() {
			continue
		}

		merged := mergeRecursive(baseFieldValue, overrideField, opts)
		// If merged is invalid -> (effectively) remove
		if merged.IsValid() {
			// Ensure correct type/pointer shape
			if result.Field(i).Kind() == reflect.Ptr && merged.Kind() != reflect.Ptr && merged.CanAddr() {
				merged = merged.Addr()
			}

			result.Field(i).Set(merged)
		}
	}

	return result
}

// mergeMap merges two map values (non-timestamped case).
func mergeMap(baseVal, overrideVal reflect.Value, opts MergeObject) reflect.Value {
	if !baseVal.IsValid() || baseVal.IsNil() {
		// missing base -> override
		return overrideVal
	}

	if !overrideVal.IsValid() || overrideVal.IsNil() {
		// override missing
		switch opts.Mode {
		case ClientIsMaster:
			return reflect.Value{} // remove
		case ServerIsMaster:
			return baseVal // keep old
		}
	}

	// Both are maps
	result := reflect.MakeMap(baseVal.Type())
	// Copy base first
	for _, key := range baseVal.MapKeys() {
		result.SetMapIndex(key, baseVal.MapIndex(key))
	}

	for _, key := range overrideVal.MapKeys() {
		ovVal := overrideVal.MapIndex(key)
		baseValForKey := baseVal.MapIndex(key)

		if !baseValForKey.IsValid() {
			// Key is new -> add/override
			result.SetMapIndex(key, ovVal)
			continue
		}

		// Merge recursively
		merged := mergeRecursive(baseValForKey, ovVal, opts)
		if merged.IsValid() {
			result.SetMapIndex(key, merged)
		} else {
			if opts.Mode == ClientIsMaster {
				result.SetMapIndex(key, reflect.Value{}) // remove
			} else {
				result.SetMapIndex(key, baseValForKey) // keep old
			}
		}
	}

	// remove keys from result if they are missing in overrideVal (ClientIsMaster only)
	if opts.Mode == ClientIsMaster {
		for _, key := range baseVal.MapKeys() {
			if !overrideVal.MapIndex(key).IsValid() {
				result.SetMapIndex(key, reflect.Value{}) // remove
			}
		}
	}

	return result
}
