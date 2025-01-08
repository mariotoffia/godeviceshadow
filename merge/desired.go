package merge

import (
	"fmt"
	"reflect"
)

// DesiredOptions holds configuration such as `MergeLoggers`.
type DesiredOptions struct {
	// Loggers will be notified on add, updated, remove, not-changed operations while merging.
	Loggers DesiredLoggers
}

type DesiredObject struct {
	DesiredOptions
	CurrentPath string
}

// Desired is a special merge where a report model is analyses if it matches the desired model. All matched
// are removed or set to `nil` in the desired model. The result is the desired model with the changes.
//
// It will only merge `model.ValueAndTimestamp` where `model.ValueAndTimestamp.GetValue()` are equal in
// the reported model. WHen that happens, it will remove the value from the desired model and report
// to `DesiredOptions.DesiredLoggers`.
func Desired[T any](reportedModel, desiredModel T, opts DesiredOptions) (T, error) {
	//
	reportedVal := reflect.ValueOf(reportedModel)
	desiredVal := reflect.ValueOf(desiredModel)

	mergedVal := desiredRecursive(reportedVal, desiredVal, DesiredObject{DesiredOptions: opts})

	if res, ok := mergedVal.Interface().(T); ok {
		return res, nil
	} else {
		var zero T

		return res, fmt.Errorf("desired merge failed, could not cast type: %T to %T", res, zero)
	}
}

func desiredRecursive(reportedVal, desiredVal reflect.Value, obj DesiredObject) reflect.Value {
	if !reportedVal.IsValid() || !desiredVal.IsValid() {
		return reflect.Value{}
	}

	// Only check IsNil for types where it is valid
	if (reportedVal.Kind() == reflect.Ptr || reportedVal.Kind() == reflect.Interface) && reportedVal.IsNil() {
		return reflect.Value{}
	}
	if (desiredVal.Kind() == reflect.Ptr || desiredVal.Kind() == reflect.Interface) && desiredVal.IsNil() {
		return reflect.Value{}
	}

	if !desiredVal.CanSet() {
		// If desiredVal is not addressable, create a new addressable copy
		desiredVal = makeAddressable(desiredVal)
	}

	if rvt, ok := unwrapValueAndTimestamp(reportedVal); ok {
		if dvt, ok := unwrapValueAndTimestamp(desiredVal); ok && rvt.GetValue() == dvt.GetValue() {
			obj.Loggers.NotifyAcknowledge(obj.CurrentPath, rvt)

			// Remove from desired model
			return reflect.Zero(desiredVal.Type())
		}

		return desiredVal
	}

	reportedVal = unwrapReflectValue(reportedVal)
	basePath := obj.CurrentPath

	switch reportedVal.Kind() {
	case reflect.Struct:
		for i := 0; i < reportedVal.NumField(); i++ {
			field := reportedVal.Type().Field(i)

			if field.PkgPath != "" {
				continue // Unexported field -> skip
			}

			tag := getJSONTag(field)

			if tag == "" {
				continue // No tag -> skip
			}

			obj.CurrentPath = concatPath(basePath, tag)

			if r := desiredRecursive(reportedVal.Field(i), desiredVal.Field(i), obj); r.IsValid() {
				desiredVal.Field(i).Set(r)
			}
		}
	case reflect.Map:
		for _, key := range reportedVal.MapKeys() {
			obj.CurrentPath = concatPath(basePath, formatKey(key))

			result := desiredRecursive(reportedVal.MapIndex(key), desiredVal.MapIndex(key), obj)

			if result.IsValid() && !result.IsZero() {
				// Update key with the new value
				desiredVal.SetMapIndex(key, result)
			} else {
				// Remove key from the map
				desiredVal.SetMapIndex(key, reflect.Value{}) // This deletes the key
			}
		}
	case reflect.Slice, reflect.Array:
		// This makes no sense - we need to introduce a `model.ValueIDAndTimeStamp` to match ID
		// as we do in a map. Do not use slices/arrays!
		minLen := reportedVal.Len()

		if minLen > desiredVal.Len() {
			minLen = desiredVal.Len()
		}

		for i := 0; i < minLen; i++ {
			obj.CurrentPath = fmt.Sprintf("%s.%d", basePath, i)

			if r := desiredRecursive(reportedVal.Index(i), desiredVal.Index(i), obj); r.IsValid() {
				desiredVal.Index(i).Set(r)
			}
		}
	}

	return desiredVal
}

func makeAddressable(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Ptr {
		return v
	}
	ptr := reflect.New(v.Type())
	ptr.Elem().Set(v)
	return ptr.Elem()
}
