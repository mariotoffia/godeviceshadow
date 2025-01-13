package persistutils

import "github.com/mariotoffia/godeviceshadow/model/persistencemodel"

type GroupedWriteOperation struct {
	// ID is from the `WriteOperation.ID.ID` field.
	ID string
	// Name is from the `WriteOperation.ID.Name` field.
	Name string
	// Operations is a slice of `WriteOperation` that is grouped by the `ID` and `Name`.
	Operations []persistencemodel.WriteOperation
	// ModelSeparation is the separation type for the model. If not set it will use the default in the `Persistor`.
	//
	// This is extracted from the `WriteOperation.Config` _separation_ key. All operations in same group must have the
	// same separation type (or nothing set).
	ModelSeparation persistencemodel.ModelSeparation
	// Error is set when something went wrong during the grouping and the whole grouping should be discarded.
	Error error
}

// GetByModelType will return the first `WriteOperation` that matches the `model` type. If not found it will return `nil`.
func (group *GroupedWriteOperation) GetByModelType(model persistencemodel.ModelType) *persistencemodel.WriteOperation {
	if group == nil {
		return nil
	}

	for _, v := range group.Operations {
		if v.ID.ModelType == model {
			return &v
		}
	}

	return nil
}

// Group will group all write operations based on the ID and Name of the model.
func Group(operations []persistencemodel.WriteOperation) []GroupedWriteOperation {
	operationsByModel := make(map[string]*GroupedWriteOperation, len(operations)/2)

	for _, op := range operations {
		key := op.ID.ID + op.ID.Name

		if _, ok := operationsByModel[key]; !ok {
			operationsByModel[key] = &GroupedWriteOperation{
				ID:         op.ID.ID,
				Name:       op.ID.Name,
				Operations: make([]persistencemodel.WriteOperation, 0, 2),
			}
		}

		opm := operationsByModel[key]
		if sep, ok := FromConfig[persistencemodel.ModelSeparation](op.AdditionalProperties, persistencemodel.ModelSeparationConfigKey); ok {
			if opm.ModelSeparation == 0 {
				opm.ModelSeparation = sep
			} else if opm.ModelSeparation != sep {
				opm.Error = persistencemodel.Error400("all operations in same group must have the same separation type")
			}
		}

		opm.Operations = append(opm.Operations, op)
	}

	groupedOperations := make([]GroupedWriteOperation, 0, len(operationsByModel))

	for _, op := range operationsByModel {
		groupedOperations = append(groupedOperations, *op)
	}

	return groupedOperations
}
