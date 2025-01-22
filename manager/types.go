package manager

import (
	"github.com/mariotoffia/godeviceshadow/model"
	"github.com/mariotoffia/godeviceshadow/model/managermodel"
	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
)

type Manager struct {
	// persistence to use for _CRUD_ operations
	persistence persistencemodel.Persistence
	// typeRegistry is the type registry that is used when no `TypeRegistryResolver` is
	// registered. It will use it to resolve the types based on the name only or if passed
	// as argument on read operation.
	typeRegistry model.TypeRegistry
	// typeRegistryResolver is primarily used to resolve a type based on it's id and name.
	typeRegistryResolver model.TypeRegistryResolver
	// reportedMergeLoggers is a slice of instantiator to produce MergeLogger(s) of which is all applied when merge operations when
	// user invokes the `Report` function.
	reportedMergeLoggers []model.CreatableMergeLogger
	// reportedDesiredLoggers is a slice of instantiator to produce DesiredLogger(s) of which is all applied when merge operations when
	// user invokes the `Report` function.
	reportedDesiredLoggers []model.CreatableDesiredLogger
	// separation is the default separation.
	separation persistencemodel.ModelSeparation
}

type groupedPersistenceResult struct {
	id            persistencemodel.ID
	reported      *persistencemodel.ReadResult
	desired       *persistencemodel.ReadResult
	queueReported any
	queueDesired  any
	op            *managermodel.ReportOperation
}
