package dynamodbpersistence

import (
	"context"

	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
)

func (p *Persistence) List(
	ctx context.Context,
	opt persistencemodel.ListOptions,
) ([]persistencemodel.ListResult, error) {
	panic("TODO: Not Implemented")
}
