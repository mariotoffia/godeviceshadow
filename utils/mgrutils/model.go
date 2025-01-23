package mgrutils

import "github.com/mariotoffia/godeviceshadow/model/managermodel"

// ReadResultToModel will try to extract the model from a `managermodel.ReadOperationResult`. If
// it succeeds, it will return it and `true`. Otherwise, it will return _"default"_ and `false`.
func ReadResultToModel[T any](rr *managermodel.ReadOperationResult) (T, bool) {
	var zero T

	if rr.Error != nil || rr.Model == nil {
		return zero, false
	}

	if m, ok := rr.Model.(T); ok {
		return m, true
	}

	return zero, false
}
