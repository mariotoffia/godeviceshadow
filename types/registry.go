package types

import (
	"fmt"
	"sync"

	"github.com/mariotoffia/godeviceshadow/model"
)

// TypeRegistryImpl implements both `model.TypeRegistry` and `model.TypeRegistryResolver`.
type TypeRegistryImpl struct {
	types     map[string]model.TypeEntry
	ids       map[string]model.TypeEntry
	id_names  map[string]model.TypeEntry
	resolvers []model.TypeRegistryResolver
	mtx       *sync.RWMutex
}

func NewRegistry() *TypeRegistryImpl {
	return &TypeRegistryImpl{
		types:    map[string]model.TypeEntry{},
		ids:      map[string]model.TypeEntry{},
		id_names: map[string]model.TypeEntry{},
		mtx:      &sync.RWMutex{},
	}
}

// Register implements the `model.TypeRegistry` interface. If name is empty, it will
// use _'{pkg name}.{type name}'_ as the name.
func (r *TypeRegistryImpl) Register(name string, t any, meta ...map[string]string) {
	te := toEntry(t, name, meta)

	r.mtx.Lock()

	r.types[te.Name] = te

	r.mtx.Unlock()
}

// Get implements the `model.TypeRegistry` interface.
func (r *TypeRegistryImpl) Get(name string) (model.TypeEntry, bool) {
	r.mtx.RLock()

	t, ok := r.types[name]

	r.mtx.RUnlock()

	return t, ok
}

// ResolveByID implements the `model.TypeRegistryResolver` interface where it will invoke all registered resolvers
// and then perform it's internal lookups (if all resolvers failed).
//
// The internal lookups are performed in the following order:
//
// 1. `id+name` lookup first - i.e. most narrow lookup.
// 2. `name` lookup second
// 3. `id` lookup last
//
// If no lookup is successful, `false` is returned.
func (r *TypeRegistryImpl) ResolveByID(id, name string) (model.TypeEntry, bool) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	for _, resolver := range r.resolvers {
		if te, ok := resolver.ResolveByID(id, name); ok {
			return te, true
		}
	}

	if te, ok := r.id_names[id+name]; ok {
		return te, true
	}

	if te, ok := r.types[name]; ok {
		return te, true
	}

	if te, ok := r.ids[id]; ok {
		return te, true
	}

	return model.TypeEntry{}, false
}

// RegisterIDandNameLookup register for the most narrow lookup id+name where both must match. Thus both _id_ and _name_
// must be a non empty string. If is was already registered, it will return an error.
//
// This will ensure a unique lookup for a specific type. It will not register the type for the _id_ and _name_ lookups.
func (r *TypeRegistryImpl) RegisterIDandNameLookup(id, name string, t any, meta ...map[string]string) error {
	if id == "" || name == "" {
		return fmt.Errorf("both id and name must be non empty strings")
	}

	te := toEntry(t, name, meta)

	r.mtx.Lock()
	defer r.mtx.Unlock()

	if _, ok := r.id_names[id+name]; ok {
		return fmt.Errorf("id+name already registered: %s#%s", id, name)
	}

	r.id_names[id+name] = te

	return nil
}

// RegisterIDLookup register a type with a id. If it was already registered, it will return an error.
func (r *TypeRegistryImpl) RegisterIDLookup(id string, t any, meta ...map[string]string) error {
	if id == "" {
		return fmt.Errorf("id must be a non empty string")
	}

	te := toEntry(t, "", meta)

	r.mtx.Lock()
	defer r.mtx.Unlock()

	if _, ok := r.ids[id]; ok {
		return fmt.Errorf("id already registered: %s", id)
	}

	r.ids[id] = te

	return nil
}

// RegisterNameLookup register a type with a name. If it was already registered, it will return an error.
func (r *TypeRegistryImpl) RegisterNameLookup(name string, t any, meta ...map[string]string) error {
	if name == "" {
		return fmt.Errorf("name must be a non empty string")
	}

	te := toEntry(t, name, meta)

	r.mtx.Lock()
	defer r.mtx.Unlock()

	if _, ok := r.types[name]; ok {
		return fmt.Errorf("name already registered: %s", name)
	}

	r.types[name] = te

	return nil
}

// RegisterResolver registers a resolver that will be invoked in the `ResolveByID` function.
//
// It is possible to register multiple resolvers and they will be invoked in the order they are registered.
// First hit will be returned.
func (r *TypeRegistryImpl) RegisterResolver(resolver model.TypeRegistryResolver) {
	r.mtx.Lock()

	r.resolvers = append(r.resolvers, resolver)

	r.mtx.Unlock()
}
