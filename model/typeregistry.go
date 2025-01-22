package model

import "reflect"

// TypeEntry is a entry returned by the `TypeRegistry`.
type TypeEntry struct {
	// Model is the model type.
	Model reflect.Type
	// Name is a arbitrary name that the model is registered with.
	Name string
	// Meta is a optional metadata key=value attached when registered.
	Meta map[string]string
}

type TypeRegistry interface {
	// Register a type with a arbitrary name.
	//
	// NOTE: This is a go routine safe operation.
	//
	// _meta_ is an optional metadata key=value that can be used for lookups.
	//
	// If _name_ is empty, it will register with the types {pkg path}.{type name}.
	Register(name string, t any, meta ...map[string]string)
	// Get will is a simple get a type by _name_. It returns, `true` if found, otherwise `false`.
	//
	// NOTE: This is a go routine safe operation.
	Get(name string) (TypeEntry, bool)
}

// TypeRegistryResolver is when a `TypeRegistry` has the ability to perform more advanced
// lookups and resolution of types. The types do not even have to be registered in the
// `TypeRegistry` but can be resolved by other means.
type TypeRegistryResolver interface {
	// ResolveByID is a more advanced lookup of a type by `id` and `name`. Where the _id_ is
	// the same as `persistencemodel.PersistenceID.ID` and `name` is the same as `persistencemodel.PersistenceID.Name`.
	//
	// This allows for model type lookups based on the `persistencemodel.PersistenceID`.
	//
	// NOTE: This is a go routine safe operation.
	//
	// This may be simulated in `TypeRegistry.Register` where name is ID+Name and can be resolved by `Get` by the
	// same way. This function do not rely on any prior registration.
	ResolveByID(id, name string) (TypeEntry, bool)
}

type TypeRegistryResolverImpl struct {
	f func(id, name string) (TypeEntry, bool)
}

func NewResolveFunc(f func(id, name string) (TypeEntry, bool)) *TypeRegistryResolverImpl {
	return &TypeRegistryResolverImpl{f: f}
}

func (r *TypeRegistryResolverImpl) ResolveByID(id, name string) (TypeEntry, bool) {
	return r.f(id, name)
}
