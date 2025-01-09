package persistencemodel

import "fmt"

type PersistenceID struct {
	// Is a unique identifier e.g. MyCar 22 or a UUID.
	ID string
	// Name is the persistence name so it is possible to have multiple device shadows (model types).
	Name string
	// ModelType is the model type that this `PersistenceID` refers to.
	ModelType ModelType
}

// ModelType stipulates the model type in e.g. a `PersistenceID`.
type ModelType int

const (
	// ModelTypeReported is a the reported portion.
	ModelTypeReported ModelType = 1
	// ModelTypeDesired is the desired portion.
	ModelTypeDesired ModelType = 2
)

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

// ModelSeparationConfigKey is the common  key to use when a `Write` request specifies the separation type in its `WriteOperation.Config`.
const ModelSeparationConfigKey = "separation"
