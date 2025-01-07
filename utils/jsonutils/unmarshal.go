package jsonutils

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// UnmarshalJSON will unmarshal a JSON into a model of type `t` and return the model as `any`.
func UnmarshalJSON(data []byte, t reflect.Type) (any, error) {
	v := reflect.New(t).Interface()

	if err := json.Unmarshal(data, v); err != nil {
		var (
			se *json.SyntaxError
			ue *json.UnmarshalTypeError
		)

		if errors.As(err, &se) {
			return nil, fmt.Errorf("%s", HighlightSyntaxError(data, se))
		} else if errors.As(err, &ue) {
			return nil, fmt.Errorf("unmarshal type error: %w", err)
		}

		return nil, err
	}

	return v, nil
}
