package mempersistence

import (
	"sync"

	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
)

type Partition map[string]*modelEntry

type Store struct {
	partitions map[string]Partition
	mu         sync.RWMutex
}

// Persistence is a in memory persistence, that stores the model without cloning.
type Persistence struct {
	store Store
	opt   PersistenceOpts
}

type modelEntry struct {
	modelType   persistencemodel.ModelType
	reported    any
	desired     any
	version     int64
	timestamp   int64
	clientToken string
}

type PersistenceOpts struct {
	// Separation is the model separation strategy. Default is `CombinedModels`.
	Separation persistencemodel.ModelSeparation
}

// New creates a new instance of InMemoryReadonlyPersistence.
func New(opts ...PersistenceOpts) *Persistence {
	var opt PersistenceOpts

	if len(opts) > 0 {
		opt = opts[0]
	}

	if opt.Separation == 0 {
		opt.Separation = persistencemodel.CombinedModels
	}

	return &Persistence{
		opt:   opt,
		store: Store{partitions: map[string]Partition{}},
	}
}
