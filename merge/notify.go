package merge

import (
	"fmt"
	"reflect"

	"github.com/mariotoffia/godeviceshadow/model"
)

func notifyRecursive(oldVal, newVal reflect.Value, obj MergeObject) error {
	if !oldVal.IsValid() && !newVal.IsValid() {
		return nil // Both are invalid, nothing to notify.
	}

	if !oldVal.IsValid() {
		obj.Loggers.NotifyPlain(obj.CurrentPath, model.MergeOperationAdd, nil, newVal.Interface())
		return nil
	}

	if !newVal.IsValid() {
		obj.Loggers.NotifyPlain(obj.CurrentPath, model.MergeOperationRemove, oldVal.Interface(), nil)
		return nil
	}

	// Unwrap pointers and interfaces
	oldVal = unwrapReflectValue(oldVal)
	newVal = unwrapReflectValue(newVal)

	if oldVal.Kind() != newVal.Kind() {
		obj.Loggers.NotifyPlain(obj.CurrentPath, model.MergeOperationUpdate, oldVal.Interface(), newVal.Interface())
		return nil
	}

	switch oldVal.Kind() {
	case reflect.Struct:
		return notifyStruct(oldVal, newVal, obj)
	case reflect.Map:
		return notifyMap(oldVal, newVal, obj)
	case reflect.Slice, reflect.Array:
		return notifySlice(oldVal, newVal, obj)
	default:
		if !reflect.DeepEqual(oldVal.Interface(), newVal.Interface()) {
			obj.Loggers.NotifyPlain(obj.CurrentPath, model.MergeOperationUpdate, oldVal.Interface(), newVal.Interface())
		} else {
			obj.Loggers.NotifyPlain(obj.CurrentPath, model.MergeOperationNotChanged, oldVal.Interface(), newVal.Interface())
		}
	}

	return nil
}

func notifyStruct(oldVal, newVal reflect.Value, obj MergeObject) error {
	basePath := obj.CurrentPath

	for i := 0; i < oldVal.NumField(); i++ {
		field := oldVal.Type().Field(i)

		if field.PkgPath != "" {
			continue // Unexported field -> skip
		}

		obj.CurrentPath = concatPath(basePath, getJSONTag(field))

		if err := notifyRecursive(oldVal.Field(i), newVal.Field(i), obj); err != nil {
			return err
		}
	}
	return nil
}

// Notify added elements in a map.
func notifyMapAdditions(newVal reflect.Value, obj MergeObject) {
	basePath := obj.CurrentPath

	for _, key := range newVal.MapKeys() {
		obj.CurrentPath = concatPath(basePath, key.String())

		notifyRecursive(reflect.Value{}, newVal.MapIndex(key), obj)
	}
}

// Notify removed elements in a map.
func notifyMapRemovals(oldVal reflect.Value, obj MergeObject) {
	basePath := obj.CurrentPath

	for _, key := range oldVal.MapKeys() {
		obj.CurrentPath = concatPath(basePath, key.String())

		notifyRecursive(oldVal.MapIndex(key), reflect.Value{}, obj)
	}
}

// Notify unchanged elements in a map.
func notifyMapNoChange(oldVal reflect.Value, obj MergeObject) {
	basePath := obj.CurrentPath

	for _, key := range oldVal.MapKeys() {
		obj.CurrentPath = concatPath(basePath, key.String())

		notifyRecursive(oldVal.MapIndex(key), oldVal.MapIndex(key), obj)
	}
}

// Notify added elements in a slice.
func notifySliceAdditions(slice reflect.Value, startIndex int, obj MergeObject) {
	basePath := obj.CurrentPath

	for i := startIndex; i < slice.Len(); i++ {
		obj.CurrentPath = fmt.Sprintf("%s.%d", basePath, i)

		notifyRecursive(reflect.Value{}, slice.Index(i), obj)
	}
}

// Notify removed elements in a slice.
func notifySliceRemovals(slice reflect.Value, startIndex int, obj MergeObject) {
	basePath := obj.CurrentPath

	for i := startIndex; i < slice.Len(); i++ {
		obj.CurrentPath = fmt.Sprintf("%s.%d", basePath, i)

		notifyRecursive(slice.Index(i), reflect.Value{}, obj)
	}
}

// Notify unchanged elements in a slice.
func notifySliceNoChange(slice reflect.Value, startIndex int, obj MergeObject) {
	basePath := obj.CurrentPath

	for i := startIndex; i < slice.Len(); i++ {
		obj.CurrentPath = fmt.Sprintf("%s.%d", basePath, i)

		notifyRecursive(slice.Index(i), slice.Index(i), obj)
	}
}

// Recursive notifier adjusted for maps to call notification helpers.
func notifyMap(oldVal, newVal reflect.Value, obj MergeObject) error {
	oldKeys := make(map[reflect.Value]struct{}, oldVal.Len())

	for _, key := range oldVal.MapKeys() {
		oldKeys[key] = struct{}{}
	}

	// Notify for common keys.
	basePath := obj.CurrentPath

	for _, key := range newVal.MapKeys() {
		obj.CurrentPath = concatPath(basePath, key.String())

		if _, exists := oldKeys[key]; exists {
			delete(oldKeys, key)

			if err := notifyRecursive(oldVal.MapIndex(key), newVal.MapIndex(key), obj); err != nil {
				return err
			}
		} else {
			obj.Loggers.NotifyPlain(obj.CurrentPath, model.MergeOperationAdd, nil, newVal.MapIndex(key).Interface())
		}
	}

	// Notify for keys only in the old map.
	for key := range oldKeys {
		obj.CurrentPath = concatPath(basePath, key.String())

		obj.Loggers.NotifyPlain(obj.CurrentPath, model.MergeOperationRemove, oldVal.MapIndex(key).Interface(), nil)
	}

	return nil
}

// Recursive notifier adjusted for slices to call notification helpers.
func notifySlice(oldVal, newVal reflect.Value, obj MergeObject) error {
	minLen := oldVal.Len()

	if newVal.Len() < minLen {
		minLen = newVal.Len()
	}

	// Notify for common indices.
	for i := 0; i < minLen; i++ {
		path := fmt.Sprintf("%s.%d", obj.CurrentPath, i)
		obj.CurrentPath = path
		if err := notifyRecursive(oldVal.Index(i), newVal.Index(i), obj); err != nil {
			return err
		}
	}

	// Notify for remaining elements.
	if oldVal.Len() > minLen {
		notifySliceRemovals(oldVal, minLen, obj)
	}

	if newVal.Len() > minLen {
		notifySliceAdditions(newVal, minLen, obj)
	}

	return nil
}
