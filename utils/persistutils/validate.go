package persistutils

import "github.com/mariotoffia/godeviceshadow/model/persistencemodel"

// ValidateModelSeparation will validate the _op_ and make sure it matches the _sep_ type.
func ValidateModelSeparation(op GroupedWriteOperation, sep persistencemodel.ModelSeparation) error {
	if sep == persistencemodel.CombinedModels && len(op.Operations) != 2 {
		return persistencemodel.Error400("both reported and desired models must be present when CombinedModels is set")
	}

	return nil
}
