package manager

import (
	"context"

	"github.com/mariotoffia/godeviceshadow/model/managermodel"
	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
)

// List will list the models. If no id or search expression is provided all models will be listed. It may be
// paged, thus check the `managermodel.ListResults.Token` to see if there's more to fetch.
func (mgr *Manager) List(ctx context.Context, options ...managermodel.ListOptions) (managermodel.ListResults, error) {
	var opt persistencemodel.ListOptions

	if len(options) > 0 {
		opt = persistencemodel.ListOptions{
			ID:         options[0].ID,
			SearchExpr: options[0].SearchExpr,
			Token:      options[0].Token,
		}
	}

	if results, err := mgr.persistence.List(ctx, opt); err != nil {
		return managermodel.ListResults{}, err
	} else {
		return managermodel.ListResults{
			Items: results.Items,
			Token: results.Token,
		}, nil
	}
}
