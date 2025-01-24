package stdmgr

import (
	"github.com/mariotoffia/godeviceshadow/model"
	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
)

type builder struct {
	m *ManagerImpl
}

func New() *builder {
	return &builder{
		m: &ManagerImpl{},
	}
}

func (b *builder) Build() *ManagerImpl {
	sep := b.m.separation

	if sep == 0 {
		sep = persistencemodel.CombinedModels
	}

	return &ManagerImpl{
		persistence:            b.m.persistence,
		separation:             sep,
		typeRegistry:           b.m.typeRegistry,
		typeRegistryResolver:   b.m.typeRegistryResolver,
		reportedMergeLoggers:   b.m.reportedMergeLoggers,
		reportedDesiredLoggers: b.m.reportedDesiredLoggers,
		desiredMergeLoggers:    b.m.desiredMergeLoggers,
	}
}

func (b *builder) WithPersistence(persistence persistencemodel.Persistence) *builder {
	b.m.persistence = persistence
	return b
}

// WithSeparation will set the default separation to use. If not set, it will default to `CombinedModels`.
func (b *builder) WithSeparation(separation persistencemodel.ModelSeparation) *builder {
	b.m.separation = separation
	return b
}

// WithTypeRegistry will set the type registry to use. If the `TypeRegistryResolver` is set and it also implements the
// `TypeRegistry` interface, it will also set the `TypeRegistry` to the resolver.
func (b *builder) WithTypeRegistry(typeRegistry model.TypeRegistry) *builder {
	b.m.typeRegistry = typeRegistry

	if tr, ok := typeRegistry.(model.TypeRegistryResolver); ok && b.m.typeRegistryResolver == nil {
		b.m.typeRegistryResolver = tr
	}

	return b
}

// WithTypeRegistryResolver will set the type registry resolver to use. If it also implements the `TypeRegistry` interface
// and it has not yet been set, it will also set the `TypeRegistry` to the resolver.
func (b *builder) WithTypeRegistryResolver(typeRegistryResolver model.TypeRegistryResolver) *builder {
	b.m.typeRegistryResolver = typeRegistryResolver

	if tr, ok := typeRegistryResolver.(model.TypeRegistry); ok && b.m.typeRegistry == nil {
		b.m.typeRegistry = tr
	}

	return b
}

func (b *builder) WithReportLoggers(reportedLoggers ...model.CreatableMergeLogger) *builder {
	b.m.reportedMergeLoggers = reportedLoggers
	return b
}

// WithDesiredMergeLoggers is the `model.CreatableMergeLogger` instances that will be used when the `Desire` operation is invoked.
func (b *builder) WithDesiredMergeLoggers(desiredLoggers ...model.CreatableMergeLogger) *builder {
	b.m.desiredMergeLoggers = desiredLoggers
	return b
}

// WithReportDesiredLoggers will set the default desired loggers for the manager instance. If none is supplied in the `Report` operation
// those will be used.
func (b *builder) WithReportDesiredLoggers(desiredLoggers ...model.CreatableDesiredLogger) *builder {
	b.m.reportedDesiredLoggers = desiredLoggers
	return b
}
