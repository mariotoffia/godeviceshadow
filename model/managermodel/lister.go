package managermodel

import (
	"context"

	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
)

type ListOptions struct {
	// ID is the id of the models to list. Under a ID there may be many named models.
	// If omitted and the SearchExpr is omitted, all IDs and their named models will be returned.
	ID string
	// SearchExpr is a search expression to use to filter the IDs. This may not be supported by
	// the `Persistence`.
	SearchExpr string
	// Token is a ID that the `Persistence` did return when there's a additional page to be fetched of results.
	Token string
}

// ListResults is a list of results returned from the `List` operation.
type ListResults struct {
	// Items are the list of results.
	Items []persistencemodel.ListResult
	// Token is set to "something" when there's a additional page to be fetched of results.
	Token string
}

// Lister is a interface that can list models.
type Lister interface {
	// List will list the models. If no id or search expression is provided all models will be listed. It may be
	// paged, thus check the `ListResults.Token` to see if there's more to fetch. If id is provided, only the
	// named models under that id will be listed.
	//
	// NOTE: Some persistence may not support search expressions and thus return an error.
	List(ctx context.Context, options ...ListOptions) (ListResults, error)
}
