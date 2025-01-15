package persistencemodel

import "fmt"

type PersistenceID struct {
	// Is a unique identifier e.g. MyCar 22 or a UUID.
	ID string
	// Name is the persistence name so it is possible to have multiple device shadows (model types).
	Name string
	// ModelType is the model type that this `PersistenceID` refers to.
	//
	// When in a read operation, if this is zero, it is assumed to be a combined model type.
	ModelType ModelType
}

func (pid PersistenceID) String() string {
	return fmt.Sprintf("%s#%s/%s", pid.ID, pid.Name, pid.ModelType.String())
}

func (pid PersistenceID) ToPersistenceID(modelType ModelType) PersistenceID {
	pid.ModelType = modelType
	return pid
}

// ModelType stipulates the model type in e.g. a `PersistenceID`.
type ModelType int

const (
	// ModelTypeReported is a the reported portion.
	ModelTypeReported ModelType = 1
	// ModelTypeDesired is the desired portion.
	ModelTypeDesired ModelType = 2
)

func (mt ModelType) String() string {
	switch mt {
	case ModelTypeReported:
		return "reported"
	case ModelTypeDesired:
		return "desired"
	}

	return fmt.Sprintf("model type id: %d", int(mt))
}

// ModelSeparation is typically used in persistence to determine how to store the models.
type ModelSeparation int

const (
	// SeparateModels specifies that the models (reported, desired) should be stored separately
	SeparateModels ModelSeparation = 1
	// Combined specifies that the models (reported, desired) should be stored together in a single DynamoDB item
	CombinedModels ModelSeparation = 2
)

func (ms ModelSeparation) String() string {
	switch ms {
	case SeparateModels:
		return "separate"
	case CombinedModels:
		return "combined"
	}

	return fmt.Sprintf("model separation id: %d", int(ms))
}
