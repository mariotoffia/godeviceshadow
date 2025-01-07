package persistencemodel

type PersistenceID struct {
	// Is a unique identifier e.g. MyCar 22 or a UUID.
	ID string
	// Name is the persistence name so it is possible to have multiple device shadows.
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
