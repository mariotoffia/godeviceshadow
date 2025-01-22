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
}

type modelEntry struct {
	reported    any
	desired     any
	modelType   persistencemodel.ModelType // when zero it is combined
	version     int64
	timestamp   int64
	clientToken string
}

// New creates a new instance of InMemoryReadonlyPersistence.
func New() *Persistence {
	return &Persistence{
		store: Store{partitions: map[string]Partition{}},
	}
}
