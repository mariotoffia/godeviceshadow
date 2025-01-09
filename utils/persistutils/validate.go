package persistutils

import "github.com/mariotoffia/godeviceshadow/model/persistencemodel"

// Validate will validate the _op_ and make sure it matches the _sep_ type.
//
// If any error is already set in _op_, it will be returned as is.
func Validate(op GroupedWriteOperation, sep persistencemodel.ModelSeparation) error {
	if op.Error != nil {
		return op.Error
	}

	if sep == persistencemodel.CombinedModels && len(op.Operations) != 2 {
		return persistencemodel.Error400("both reported and desired models must be present when CombinedModels is set")
	}

	return nil
}
