package merge

import (
	"fmt"
	"reflect"
	"time"

	"github.com/mariotoffia/godeviceshadow/model"
)

// notifyRecursive will recursively notify the leafs of the _op_ operation.
//
// NOTE: _op_ may only be `model.MergeOperationAdd`, `model.MergeOperationRemove`,  and
// `model.MergeOperationNotChanged` nothing else.
func notifyRecursive(val reflect.Value, op model.MergeOperation, obj MergeObject) {
	if len(obj.Loggers) == 0 {
		return
	}

	if vt, ok := unwrapValueAndTimestamp(val); ok {
		switch op {
		case model.MergeOperationAdd:
			obj.Loggers.NotifyManaged(obj.CurrentPath, op, nil, vt, time.Time{}, vt.GetTimestamp())
		case model.MergeOperationRemove:
			obj.Loggers.NotifyManaged(obj.CurrentPath, op, vt, nil, vt.GetTimestamp(), time.Time{})
		case model.MergeOperationNotChanged:
			obj.Loggers.NotifyManaged(obj.CurrentPath, op, vt, vt, vt.GetTimestamp(), vt.GetTimestamp())
		}

		return
	}

	val = unwrapReflectValue(val)
	basePath := obj.CurrentPath

	switch val.Kind() {
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			field := val.Type().Field(i)

			if field.PkgPath != "" {
				continue // Unexported field -> skip
			}

			obj.CurrentPath = concatPath(basePath, getJSONTag(field))

			notifyRecursive(val.Field(i), op, obj)
		}
	case reflect.Map:
		for _, key := range val.MapKeys() {
			obj.CurrentPath = concatPath(basePath, key.String())

			notifyRecursive(val.MapIndex(key), op, obj)
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			obj.CurrentPath = fmt.Sprintf("%s.%d", basePath, i)

			notifyRecursive(val.Index(i), op, obj)
		}
	default:
		if !val.IsValid() {
			return // Do not notify nil values
		}

		switch op {
		case model.MergeOperationAdd:
			obj.Loggers.NotifyPlain(obj.CurrentPath, op, nil, val.Interface())
		case model.MergeOperationRemove:
			obj.Loggers.NotifyPlain(obj.CurrentPath, op, val.Interface(), nil)
		case model.MergeOperationNotChanged:
			obj.Loggers.NotifyPlain(obj.CurrentPath, op, val.Interface(), val.Interface())
		}
	}
}
