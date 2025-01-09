package dynamodbpersistence

import (
	"context"

	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
)

func (p *Persistence) Read(
	ctx context.Context,
	opt persistencemodel.ReadOptions,
	operations ...persistencemodel.ReadOperation,
) []persistencemodel.ReadResult {
	panic("TODO: Not Implemented")
}
