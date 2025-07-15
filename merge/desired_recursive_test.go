package merge_test

import (
	"context"
	"testing"

	"github.com/mariotoffia/godeviceshadow/merge"
	"github.com/stretchr/testify/assert"
)

func TestDesiredRecursiveEdgeCases(t *testing.T) {
	t.Run("Map_With_Nil_Values", func(t *testing.T) {
		// Create maps with nil values
		type TestMap struct {
			Data map[string]interface{} `json:"data"`
		}

		reported := TestMap{
			Data: map[string]interface{}{
				"key1": nil,
				"key2": "value2",
			},
		}

		desired := TestMap{
			Data: map[string]interface{}{
				"key1": nil,
				"key2": "value2",
			},
		}

		// This should not panic
		result, err := merge.Desired(context.Background(), reported, desired, merge.DesiredOptions{})
		assert.NoError(t, err)
		assert.NotNil(t, result.Data)
	})
}
