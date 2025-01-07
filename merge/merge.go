package merge

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/mariotoffia/godeviceshadow/model"
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
	// Mode specified how the merge should be performed regarding deletions.
	Mode MergeMode
	// DoOverrideWithEmpty when set to `true` it will check if the override value is "empty"
	// i.e. zero value zero len array, map, slice, nil pointer, etc. If it is `true` and the
	// override value is "empty" it will override the base value. If set to `false` it will
	// keep the base value.
	DoOverrideWithEmpty bool
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

	if err := opts.Loggers.NotifyPrepare(); err != nil {
		return oldModel, err
	}

	mergedVal, err := mergeRecursive(oldVal, newVal, MergeObject{MergeOptions: opts})

	if err2 := opts.Loggers.NotifyPost(err); err2 != nil {
		return oldModel, err2
	}

	var zero T

	if err != nil {
		return zero, err
	}

	res, _ := mergedVal.Interface().(T)

	return res, nil
}

// mergeRecursive will try to merge base with override recursively.
func mergeRecursive(base, override reflect.Value, obj MergeObject) (reflect.Value, error) {
	if !override.IsValid() {
		return reflect.Value{}, fmt.Errorf("both base: '%T' and override: '%T' must be valid", base.Interface(), override.Interface())
	}

	baseVal := base
	overrideVal := override

	// If either is an interface, unwrap to get the concrete type.
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
	baseValTS, baseOk := unwrapValueAndTimestamp(baseVal)
	overrideValTS, overrideOk := unwrapValueAndTimestamp(overrideVal)

	if baseOk && overrideOk {
		oldTS := baseValTS.GetTimestamp()
		newTS := overrideValTS.GetTimestamp()

		switch {
		case newTS.After(oldTS):
			obj.Loggers.NotifyManaged(obj.CurrentPath, model.MergeOperationUpdate, baseValTS, overrideValTS, oldTS, newTS)

			return override, nil // override newer -> replace
		default:
			obj.Loggers.NotifyManaged(obj.CurrentPath, model.MergeOperationNotChanged, baseValTS, overrideValTS, oldTS, newTS)

			return base, nil // override less or equal -> no update -> keep old
		}
	}

	switch baseVal.Kind() {
	case reflect.Struct:
		return mergeStruct(baseVal, overrideVal, obj)
	case reflect.Map:
		return mergeMap(baseVal, overrideVal, obj)
	case reflect.Slice, reflect.Array:
		return mergeSlice(baseVal, overrideVal, obj)
	// Basic types (int, string, ...)
	default:
		if obj.DoOverrideWithEmpty && isEmptyValue(overrideVal) {
			return override, nil
		}

		bv := baseVal.Interface()
		ov := overrideVal.Interface()

		if bv != ov {
			obj.Loggers.NotifyPlain(obj.CurrentPath, model.MergeOperationUpdate, bv, ov)

			return override, nil // not equal -> override
		} else {
			obj.Loggers.NotifyPlain(obj.CurrentPath, model.MergeOperationNotChanged, bv, ov)

			return base, nil // equal -> keep old
		}
	}
}

func mergeSlice(baseVal, overrideVal reflect.Value, opts MergeObject) (reflect.Value, error) {
	if baseVal.IsNil() {
		notifyRecursive(overrideVal, model.MergeOperationAdd, opts)

		return overrideVal, nil
	}

	basePath := opts.CurrentPath

	// override is nil -> remove or keep
	if overrideVal.IsNil() {
		switch opts.Mode {
		case ClientIsMaster:
			for i := 0; i < baseVal.Len(); i++ {
				opts.CurrentPath = fmt.Sprintf("%s.%d", basePath, i)

				notifyRecursive(baseVal.Index(i), model.MergeOperationRemove, opts)
			}
			return overrideVal, nil
		case ServerIsMaster:
			for i := 0; i < baseVal.Len(); i++ {
				opts.CurrentPath = fmt.Sprintf("%s.%d", basePath, i)

				notifyRecursive(baseVal.Index(i), model.MergeOperationNotChanged, opts)
			}
			return baseVal, nil
		}
	}

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

		opts.CurrentPath = fmt.Sprintf("%s.%d", basePath, i)
		mergedElem, err := mergeRecursive(baseElem, ovElem, opts)

		if err != nil {
			return reflect.Value{}, err
		}

		result = reflect.Append(result, mergedElem)
	}

	// new slice is longer -> add extra elements in override
	if ovLen > minLen {
		for i := minLen; i < ovLen; i++ {
			opts.CurrentPath = fmt.Sprintf("%s.%d", basePath, i)
			result = reflect.Append(result, overrideVal.Index(i))

			notifyRecursive(overrideVal.Index(i), model.MergeOperationAdd, opts)
		}
	}

	// old slice is longer -> remove or keep
	if baseLen > minLen {
		if opts.Mode == ServerIsMaster {
			// ServerIsMaster -> keep
			for i := minLen; i < baseLen; i++ {
				opts.CurrentPath = fmt.Sprintf("%s.%d", basePath, i)
				result = reflect.Append(result, baseVal.Index(i))

				notifyRecursive(baseVal.Index(i), model.MergeOperationNotChanged, opts)
			}
		} else /*ClientIsMaster*/ {
			for i := minLen; i < baseLen; i++ {
				opts.CurrentPath = fmt.Sprintf("%s.%d", basePath, i)

				notifyRecursive(baseVal.Index(i), model.MergeOperationRemove, opts)
			}
		}
	}

	return result, nil
}

// mergeStruct merges two struct values (non-timestamped case).
func mergeStruct(baseVal, overrideVal reflect.Value, opts MergeObject) (reflect.Value, error) {
	if !baseVal.IsValid() || !overrideVal.IsValid() {
		return reflect.Value{}, fmt.Errorf("both base: '%T' and override: '%T' must be valid", baseVal.Interface(), overrideVal.Interface())
	}

	result := reflect.New(baseVal.Type()).Elem()
	numFields := baseVal.NumField()
	basePath := opts.CurrentPath

	for i := 0; i < numFields; i++ {
		fieldValue := baseVal.Field(i)
		fieldType := baseVal.Type().Field(i)

		if fieldType.PkgPath != "" {
			continue // Unexported field -> skip
		}

		if !result.Field(i).CanSet() {
			continue
		}

		opts.CurrentPath = concatPath(basePath, getJSONTag(fieldType))
		merged, err := mergeRecursive(fieldValue, overrideVal.Field(i), opts)

		if err != nil {
			return reflect.Value{}, err
		}

		result.Field(i).Set(merged)
	}

	return result, nil
}

// mergeMap merges two map values (non-timestamped case).
func mergeMap(baseVal, overrideVal reflect.Value, opts MergeObject) (reflect.Value, error) {
	if baseVal.IsNil() {
		notifyRecursive(overrideVal, model.MergeOperationAdd, opts)

		return overrideVal, nil
	}

	basePath := opts.CurrentPath

	if overrideVal.IsNil() {
		switch opts.Mode {
		case ClientIsMaster:
			for _, key := range baseVal.MapKeys() {
				opts.CurrentPath = concatPath(basePath, formatKey(key))

				notifyRecursive(baseVal.MapIndex(key), model.MergeOperationRemove, opts)
			}

			return overrideVal, nil
		case ServerIsMaster:
			for _, key := range baseVal.MapKeys() {
				opts.CurrentPath = concatPath(basePath, formatKey(key))

				notifyRecursive(baseVal.MapIndex(key), model.MergeOperationNotChanged, opts)
			}

			return baseVal, nil
		}
	}

	result := reflect.MakeMap(baseVal.Type())

	// Base keys
	baseKeys := make(map[string]reflect.Value, baseVal.Len())

	for _, key := range baseVal.MapKeys() {
		baseKeys[formatKey(key)] = key
	}

	for _, key := range overrideVal.MapKeys() {
		overrideVal := overrideVal.MapIndex(key)
		baseValForKey := baseVal.MapIndex(key)

		opts.CurrentPath = concatPath(basePath, formatKey(key))

		if !baseValForKey.IsValid() {
			result.SetMapIndex(key, overrideVal) // add

			notifyRecursive(overrideVal, model.MergeOperationAdd, opts)

			continue
		}

		// Remove key from base
		delete(baseKeys, formatKey(key))

		// Merge recursively
		merged, err := mergeRecursive(baseValForKey, overrideVal, opts)

		if err != nil {
			return reflect.Value{}, err
		}

		result.SetMapIndex(key, merged) // Already notified in (mergeRecursive)
	}

	// keys in base (but not in override)
	for k, v := range baseKeys {
		opts.CurrentPath = concatPath(basePath, k)

		if opts.Mode == ServerIsMaster {
			result.SetMapIndex(v, baseVal.MapIndex(v)) // keep

			notifyRecursive(baseVal.MapIndex(v), model.MergeOperationNotChanged, opts)
		} else /*ClientIsMaster*/ {
			notifyRecursive(baseVal.MapIndex(v), model.MergeOperationRemove, opts)
		}
	}

	return result, nil
}

func unwrapValueAndTimestamp(rv reflect.Value) (model.ValueAndTimestamp, bool) {
	// Must be a pointer or interface with non nil value
	if !rv.IsValid() || rv.Kind() == reflect.Invalid || (rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface) && rv.IsNil() {
		return nil, false
	}

	// Unwrap pointers and interfaces recursively
	rv = unwrapReflectValue(rv)

	if !rv.IsValid() {
		return nil, false
	}

	// Try to convert the value directly
	if vt, ok := toValueAndTimestamp(rv); ok {
		return vt, true
	}

	if !rv.CanAddr() {
		copy := reflect.New(rv.Type())
		copy.Elem().Set(rv)
		rv = copy
	} else {
		rv = rv.Addr()
	}

	// Attempt conversion on the copied addressable value
	if vt, ok := toValueAndTimestamp(rv); ok {
		return vt, true
	}

	return nil, false
}

// unwrapReflectValue unwraps pointers and interfaces recursively so it returns the
// first non-pointer/interface value.
func unwrapReflectValue(rv reflect.Value) reflect.Value {
	for rv.IsValid() && (rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface) {
		if rv.IsNil() {
			return reflect.Value{}
		}

		rv = rv.Elem()
	}

	return rv
}

// toValueAndTimestamp attempts to convert a reflect.Value to ValueAndTimestamp.
func toValueAndTimestamp(rv reflect.Value) (model.ValueAndTimestamp, bool) {
	//
	if rv.Kind() == reflect.Interface || rv.Kind() == reflect.Ptr {
		if i, ok := rv.Interface().(model.ValueAndTimestamp); ok {
			return i, true
		}
	}

	return nil, false
}

// isEmptyValue checks if a reflect.Value is valid or "the zero value" (len of zero).
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

	// comma -> ignore the rest
	idx := strings.Index(tag, ",")
	if idx != -1 {
		return tag[:idx]
	}

	return tag
}
