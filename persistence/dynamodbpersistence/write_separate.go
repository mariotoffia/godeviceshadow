package dynamodbpersistence

import (
	"context"
	"time"

	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
	"github.com/mariotoffia/godeviceshadow/utils/persistutils"
)

// writeSeparate will write reported and desired models separately to the persistence layer. It will use a transaction
// if both models are present, otherwise it will use a single conditional write.
func (p *Persistence) writeSeparate(
	ctx context.Context,
	group persistutils.GroupedWriteOperation,
) []persistencemodel.WriteResult {
	reported := group.GetByModelType(persistencemodel.ModelTypeReported)
	desired := group.GetByModelType(persistencemodel.ModelTypeDesired)

	if reported != nil && desired != nil {
		// Both reported and desired models present, use a transaction
		return p.writeTransactional(ctx, group, reported, desired)
	}

	// Only one model is present -> use "plain" conditional writes
	results := make([]persistencemodel.WriteResult, 0, len(group.Operations))
	if reported != nil {
		results = append(results, p.writeSingle(ctx, reported, persistencemodel.ModelTypeReported))
	}
	if desired != nil {
		results = append(results, p.writeSingle(ctx, desired, persistencemodel.ModelTypeDesired))
	}

	return results
}

// writeSingle performs a single conditional put for either reported or desired models.
func (p *Persistence) writeSingle(
	ctx context.Context,
	op *persistencemodel.WriteOperation,
	modelType persistencemodel.ModelType,
) persistencemodel.WriteResult {
	//
	pk := toPartitionKey(op.ID)
	sk := toSortKey(op.ID, modelType)
	now := time.Now().UTC().UnixNano()

	obj := PersistenceObject{
		Version:     op.Version + 1,
		TimeStamp:   now,
		ClientToken: op.ClientID,
	}

	if modelType == persistencemodel.ModelTypeReported {
		obj.Reported = op.Model
	} else {
		obj.Desired = op.Model
	}

	err := p.dynamoDbPut(ctx, pk, sk, obj, op.Version)

	version := obj.Version

	if err != nil {
		version = op.Version
	}

	return persistencemodel.WriteResult{
		ID:        op.ID,
		Version:   version,
		TimeStamp: obj.TimeStamp,
		Error:     conditionalWriteErrorFixup(err),
	}
}
