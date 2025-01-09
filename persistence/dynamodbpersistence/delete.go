package dynamodbpersistence

import (
	"context"

	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
)

func (p *Persistence) Delete(
	ctx context.Context,
	opt persistencemodel.WriteOptions,
	operations ...persistencemodel.WriteOperation,
) []persistencemodel.WriteResult {
	panic("TODO: Not Implemented")
}
