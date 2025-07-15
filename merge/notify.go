package merge

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/mariotoffia/godeviceshadow/model"
)

// notifyRecursive will recursively notify the leafs of the _op_ operation.
//
// NOTE: _op_ may only be `model.MergeOperationAdd`, `model.MergeOperationRemove`,  and
// `model.MergeOperationNotChanged` nothing else.
func notifyRecursive(ctx context.Context, val reflect.Value, op model.MergeOperation, obj MergeObject) {
	if len(obj.Loggers) == 0 {
		return
	}

	if vt, ok := unwrapValueAndTimestamp(val); ok {
		switch op {
		case model.MergeOperationAdd:
			obj.Loggers.NotifyManaged(ctx, obj.CurrentPath, op, nil, vt, time.Time{}, vt.GetTimestamp())
		case model.MergeOperationRemove:
			obj.Loggers.NotifyManaged(ctx, obj.CurrentPath, op, vt, nil, vt.GetTimestamp(), time.Time{})
		case model.MergeOperationNotChanged:
			obj.Loggers.NotifyManaged(ctx, obj.CurrentPath, op, vt, vt, vt.GetTimestamp(), vt.GetTimestamp())
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

			tag := getJSONTag(field)

			if tag == "" {
				continue // No tag -> skip
			}

			obj.CurrentPath = concatPath(basePath, tag)

			notifyRecursive(ctx, val.Field(i), op, obj)
		}
	case reflect.Map:
		for _, key := range val.MapKeys() {
			obj.CurrentPath = concatPath(basePath, formatKey(key))

			notifyRecursive(ctx, val.MapIndex(key), op, obj)
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			obj.CurrentPath = fmt.Sprintf("%s.%d", basePath, i)

			notifyRecursive(ctx, val.Index(i), op, obj)
		}
	default:
		if !val.IsValid() {
			return // Do not notify nil values
		}

		switch op {
		case model.MergeOperationAdd:
			obj.Loggers.NotifyPlain(ctx, obj.CurrentPath, op, nil, val.Interface())
		case model.MergeOperationRemove:
			obj.Loggers.NotifyPlain(ctx, obj.CurrentPath, op, val.Interface(), nil)
		case model.MergeOperationNotChanged:
			obj.Loggers.NotifyPlain(ctx, obj.CurrentPath, op, val.Interface(), val.Interface())
		}
	}
}
