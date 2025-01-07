package merge

import (
	"fmt"
	"reflect"
	"strings"
)

func formatKey(key reflect.Value) string {
	if !key.IsValid() {
		return "<invalid>"
	}

	switch key.Kind() {
	case reflect.String:
		return key.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", key.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("%d", key.Uint())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%f", key.Float())
	case reflect.Struct:
		var fields []string
		for i := 0; i < key.NumField(); i++ {
			field := key.Type().Field(i)
			value := key.Field(i).Interface()
			fields = append(fields, fmt.Sprintf("%s:%v", field.Name, value))
		}
		return fmt.Sprintf("{%s}", strings.Join(fields, ","))
	case reflect.Interface:
		// Recursively handle the underlying type of the interface
		return formatKey(key.Elem())
	default:
		return fmt.Sprintf("<unsupported key type: %s>", key.Type())
	}
}
