package dynamodbpersistence

import (
	"context"
	"fmt"
	"time"

	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
	"github.com/mariotoffia/godeviceshadow/utils/persistutils"
)

// writeCombined will combine both reported and desired models into a single persistence object and write it to the
// persistence layer. It will only do one conditional check on version and write the increased version in same operation.
func (p *Persistence) writeCombined(
	ctx context.Context,
	group persistutils.GroupedWriteOperation,
) []persistencemodel.WriteResult {
	// Determine PK and SK for combined models
	pk := fmt.Sprintf("DS#%s", group.ID)
	sk := fmt.Sprintf("DSC#%s", group.Name)

	reported := group.GetByModelType(persistencemodel.ModelTypeReported)
	desired := group.GetByModelType(persistencemodel.ModelTypeDesired)

	// Create PersistenceObject with both Desired and Reported
	now := time.Now().UTC().UnixNano()

	obj := PersistenceObject{
		Version:     reported.Version + 1,
		TimeStamp:   now,
		ClientToken: reported.ClientID,
		Desired:     desired.Model,
		Reported:    reported.Model,
	}

	// Perform conditional write
	err := p.dynamoDbPut(ctx, pk, sk, obj, reported.Version)

	version := obj.Version

	if err != nil {
		version = reported.Version
	}

	return []persistencemodel.WriteResult{
		{
			ID:        reported.ID,
			Version:   version,
			TimeStamp: obj.TimeStamp,
			Error:     conditionalWriteErrorFixup(err),
		},
		{
			ID:        desired.ID,
			Version:   version,
			TimeStamp: obj.TimeStamp,
			Error:     conditionalWriteErrorFixup(err),
		},
	}
}
