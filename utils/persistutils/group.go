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
	ModelSeparation persistencemodel.ModelSeparation
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

		operationsByModel[key].Operations = append(operationsByModel[key].Operations, op)
	}

	groupedOperations := make([]GroupedWriteOperation, 0, len(operationsByModel))

	for _, op := range operationsByModel {
		groupedOperations = append(groupedOperations, *op)
	}

	return groupedOperations
}
