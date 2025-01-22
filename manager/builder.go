package manager

import (
	"github.com/mariotoffia/godeviceshadow/model"
	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
)

type Builder struct {
	m *Manager
}

func New() *Builder {
	return &Builder{
		m: &Manager{},
	}
}

func (b *Builder) Build() *Manager {
	sep := b.m.separation

	if sep == 0 {
		sep = persistencemodel.CombinedModels
	}

	return &Manager{
		persistence:            b.m.persistence,
		separation:             sep,
		typeRegistry:           b.m.typeRegistry,
		typeRegistryResolver:   b.m.typeRegistryResolver,
		reportedMergeLoggers:   b.m.reportedMergeLoggers,
		reportedDesiredLoggers: b.m.reportedDesiredLoggers,
	}
}

func (b *Builder) WithPersistence(persistence persistencemodel.Persistence) *Builder {
	b.m.persistence = persistence
	return b
}

// WithSeparation will set the default separation to use. If not set, it will default to `CombinedModels`.
func (b *Builder) WithSeparation(separation persistencemodel.ModelSeparation) *Builder {
	b.m.separation = separation
	return b
}

// WithTypeRegistry will set the type registry to use. If the `TypeRegistryResolver` is set and it also implements the
// `TypeRegistry` interface, it will also set the `TypeRegistry` to the resolver.
func (b *Builder) WithTypeRegistry(typeRegistry model.TypeRegistry) *Builder {
	b.m.typeRegistry = typeRegistry

	if tr, ok := typeRegistry.(model.TypeRegistryResolver); ok && b.m.typeRegistryResolver == nil {
		b.m.typeRegistryResolver = tr
	}

	return b
}

// WithTypeRegistryResolver will set the type registry resolver to use. If it also implements the `TypeRegistry` interface
// and it has not yet been set, it will also set the `TypeRegistry` to the resolver.
func (b *Builder) WithTypeRegistryResolver(typeRegistryResolver model.TypeRegistryResolver) *Builder {
	b.m.typeRegistryResolver = typeRegistryResolver

	if tr, ok := typeRegistryResolver.(model.TypeRegistry); ok && b.m.typeRegistry == nil {
		b.m.typeRegistry = tr
	}

	return b
}

func (b *Builder) WithReportedLoggers(reportedLoggers []model.CreatableMergeLogger) *Builder {
	b.m.reportedMergeLoggers = reportedLoggers
	return b
}

// WithDesiredLoggers will set the default desired loggers for the manager instance. If none is supplied in the `Report` operation
// those will be used.
func (b *Builder) WithDesiredLoggers(desiredLoggers []model.CreatableDesiredLogger) *Builder {
	b.m.reportedDesiredLoggers = desiredLoggers
	return b
}
