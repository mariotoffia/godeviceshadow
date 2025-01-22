package mempersistence

import (
	"context"

	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
)

// List lists models in the in-memory persistence. SearchExpr is not supported.
func (p *Persistence) List(
	ctx context.Context,
	opt persistencemodel.ListOptions,
) (persistencemodel.ListResults, error) {
	//
	if opt.SearchExpr != "" {
		return persistencemodel.ListResults{}, persistencemodel.Error400("SearchExpr is not supported")
	}

	if opt.Token != "" {
		return persistencemodel.ListResults{}, persistencemodel.Error400("Token is not supported")
	}

	return p.store.List(opt.ID), nil
}
