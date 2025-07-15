# GoDeviceShadow AI Coding Assistant Guide

## Project Overview

GoDeviceShadow is a pluggable device shadow implementation with persistence and notification capabilities. It provides a modular architecture for managing reported and desired states of device shadows with configurable merging strategies.

## Core Architecture

### Key Components

1. **Merge System**: Core functionality for merging reported and desired states
   - Located in: `/merge`
   - Implements timestamp-based merging strategies
   - Supports two merging modes: `ClientIsMaster` and `ServerIsMaster`

2. **Manager**: Coordinates operations on device shadows
   - Located in: `/manager/stdmgr`
   - Built using a builder pattern
   - Configurable with persistence, loggers, and type registry

3. **Persistence Layer**: Pluggable storage for device shadows
   - Located in: `/persistence`
   - Built-in implementations: in-memory (`mempersistence`), DynamoDB (`dynamodbpersistence`)

4. **Notification System**: Configurable event notification
   - Located in: `/notify`
   - Supports selection filters for targeted notifications
   - Extensible with custom notification targets

5. **Loggers**: Track changes during merge operations
   - Located in: `/loggers`
   - Implementations: `changelogger`, `desirelogger`, etc.

### Data Flow

1. Client reports state via `Manager.Report()`
2. Manager merges with existing state using `merge` package
3. Loggers track changes during merge
4. Notification system notifies relevant targets
5. Client can set desired state via `Manager.Desire()`

## Key Interfaces and Models

### ValueAndTimestamp Interface

```go
type ValueAndTimestamp interface {
  GetTimestamp() time.Time
  GetValue() any
}
```

Fields implementing this interface participate in timestamp-based merging. The newer timestamp always wins.

### IdValueAndTimestamp Interface

```go
type IdValueAndTimestamp interface {
  ValueAndTimestamp
  GetID() string
}
```

Extends ValueAndTimestamp to enable merging slice elements by ID instead of position.

## Development Workflow

### Building the Project

```bash
# Run tests
make test

# Run integration tests
make integration-test

# Create a new version
make version v=vX.Y.Z
```

### Common Patterns

1. **Creating a Manager**:
   ```go
   mgr := stdmgr.New().
     WithPersistence(mempersistence.New()).
     WithReportedLoggers(changelogger.New()).
     WithTypeRegistryResolver(...).
     Build()
   ```

2. **Reporting State**:
   ```go
   mgr.Report(ctx, managermodel.ReportOperation{
     ClientID: "myClient",
     ID: persistencemodel.ID{ID: "deviceId", Name: "modelName"},
     Model: myModel,
   })
   ```

3. **Setting Desired State**:
   ```go
   mgr.Desire(ctx, managermodel.DesireOperation{
     ClientID: "myClient",
     ID: persistencemodel.ID{ID: "deviceId", Name: "modelName"},
     Model: myDesiredModel,
   })
   ```

4. **Creating Notification Selection**:
   ```go
   // Using DSL (experimental)
   sel, err := selectlang.ToSelection(`SELECT * FROM Notification WHERE...`)

   // Or programmatically
   notificationManager := notify.NewBuilder().
     TargetBuilder(...).
     WithSelection(...).
     Build().
     Build()
   ```

## Project Conventions

1. **Submodules**: External dependencies are isolated in submodules to keep the core framework dependency-free
2. **Type Registry**: Models are resolved through a registry to maintain type safety
3. **Builder Pattern**: Components are typically created using builders (e.g., `stdmgr.New()`, `notify.NewBuilder()`)
4. **Pluggable Design**: All components (persistence, notification, logging) follow interfaces for extensibility

## Common Gotchas

1. Always implement the `ValueAndTimestamp` interface for fields that should participate in timestamp-based merging
2. Pay attention to the merge mode (`ClientIsMaster` vs `ServerIsMaster`) as it affects how deletions are handled
3. Use `omitempty` in JSON struct tags for desired state fields to avoid overwriting with zero values
4. The selection language DSL is experimental and may change in future versions
