package manager

import (
	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
	"github.com/mariotoffia/godeviceshadow/types"
)

type Manager struct {
	// Persistence to use for _CRUD_ operations
	Persistence  persistencemodel.Persistence
	TypeRegistry *types.TypeRegistryImpl
}
