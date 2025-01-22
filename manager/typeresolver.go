package manager

import (
	"github.com/mariotoffia/godeviceshadow/model"
	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
)

// ResolveType will first try _name_ (if any) and then the `id` (if provided) to resolve the type using the registered type resolvers.
//
// It will start with name since this is a explicit resolve operation. It will then use the `model.model.TypeRegistryResolver` and _id_
// to resolve. If that fails, it will use `persistencemodel.ID.Name` and the `model.TypeRegistry` to get at
// a last resort.
//
// NOTE: It is possible to provide with a set of _id_(s).
func (mgr *Manager) ResolveType(name string, id ...persistencemodel.ID) (model.TypeEntry, bool) {
	if name != "" && mgr.typeRegistry != nil {
		if t, ok := mgr.typeRegistry.Get(name); ok {
			return t, true
		}
	}

	if len(id) == 0 {
		return model.TypeEntry{}, false
	}

	if mgr.typeRegistryResolver != nil {
		for _, iid := range id {
			if t, ok := mgr.typeRegistryResolver.ResolveByID(iid.ID, iid.Name); ok {
				return t, true
			}
		}
	}

	if mgr.typeRegistry != nil {
		for _, iid := range id {
			if t, ok := mgr.typeRegistry.Get(iid.Name); ok {
				return t, true
			}
		}
	}

	return model.TypeEntry{}, false
}
