package dynamodbpersistence_test

import (
	"context"
	"testing"

	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
	"github.com/mariotoffia/godeviceshadow/persistence/dynamodbpersistence"
	"github.com/mariotoffia/godeviceshadow/persistence/dynamodbpersistence/dynamodbutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListEmpty(t *testing.T) {
	ctx := context.TODO()
	res := dynamodbutils.NewTestTableResource(ctx, "go-deviceshadow-test")

	defer res.Dispose(ctx, false /*delete*/)

	persistence, err := dynamodbpersistence.New(ctx, dynamodbpersistence.Config{
		Table:  res.Table,
		Client: res.Client,
	})

	require.NoError(t, err)

	results, err := persistence.List(ctx, persistencemodel.ListOptions{
		ID: "test",
	})

	require.NoError(t, err)
	assert.Len(t, results, 0)
}
