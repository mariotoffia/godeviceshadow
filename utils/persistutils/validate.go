package persistutils

import "github.com/mariotoffia/godeviceshadow/model/persistencemodel"

// Validate will validate the _op_ and make sure it matches the _separation_.
//
// If any error is already set in _op_, it will be returned as is.
func Validate(op GroupedWriteOperation) error {
	if op.Error != nil {
		return op.Error
	}

	if op.ModelSeparation == persistencemodel.CombinedModels {
		if len(op.Operations) != 2 {
			return persistencemodel.Error400("both reported and desired models must be present when CombinedModels is set")
		}

		if op.Operations[0].ID.ModelType == op.Operations[1].ID.ModelType {
			return persistencemodel.Error400("both reported and desired models must be of different type")
		}

		if op.Operations[0].ID.ModelType == 0 {
			return persistencemodel.Error400("both reported and desired models must have the correct model type")
		}

		if op.GetByModelType(persistencemodel.ModelTypeDesired) == nil {
			return persistencemodel.Error400("desired model must be present")
		}

		if op.GetByModelType(persistencemodel.ModelTypeReported) == nil {
			return persistencemodel.Error400("reported model must be present")
		}
	}

	return nil
}
