package merge_test

import (
	"context"
	"testing"
	"time"

	"github.com/mariotoffia/godeviceshadow/merge"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMapWithStructHolder(t *testing.T) {
	now := time.Now().UTC()
	oldMap := map[string]TimestampedMapVal{
		"keep":   {Value: "old", UpdatedAt: now.Add(-2 * time.Hour)},
		"remove": {Value: "willRemove", UpdatedAt: now.Add(-1 * time.Hour)},
	}
	newMap := map[string]TimestampedMapVal{
		"keep":  {Value: "new", UpdatedAt: now.Add(-1 * time.Hour)}, // old is older => replaced if newer
		"added": {Value: "brandNew", UpdatedAt: now},                // brand new key
	}

	type MapHolder struct {
		M map[string]TimestampedMapVal
	}

	oldObj := MapHolder{M: oldMap}
	newObj := MapHolder{M: newMap}

	merged, err := merge.Merge(context.Background(), oldObj, newObj, merge.MergeOptions{
		Mode: merge.ClientIsMaster,
	})
	require.NoError(t, err)
	require.NotNil(t, merged.M)

	require.Contains(t, merged.M, "keep")
	assert.Equal(t, "new", merged.M["keep"].Value, "Should override because new is strictly newer")

	require.Contains(t, merged.M, "added")
	assert.Equal(t, "brandNew", merged.M["added"].Value)

	require.NotContains(t, merged.M, "remove")

	mergedSM, err := merge.Merge(context.Background(), oldObj, newObj, merge.MergeOptions{
		Mode: merge.ServerIsMaster,
	})

	require.NoError(t, err)
	require.NotNil(t, mergedSM.M)

	require.Contains(t, mergedSM.M, "remove")
}

func TestDirectMapMerging(t *testing.T) {
	// Test merging maps directly (without a container struct)
	now := time.Now().UTC()
	oldMap := map[string]TimestampedMapVal{
		"keep":   {Value: "old", UpdatedAt: now.Add(-2 * time.Hour)},
		"remove": {Value: "willRemove", UpdatedAt: now.Add(-1 * time.Hour)},
	}
	newMap := map[string]TimestampedMapVal{
		"keep":  {Value: "new", UpdatedAt: now.Add(-1 * time.Hour)}, // old is older => replaced if newer
		"added": {Value: "brandNew", UpdatedAt: now},                // brand new key
	}

	merged, err := merge.Merge(context.Background(), oldMap, newMap, merge.MergeOptions{
		Mode: merge.ClientIsMaster,
	})
	require.NoError(t, err)
	require.NotNil(t, merged)

	require.Contains(t, merged, "keep")
	assert.Equal(t, "new", merged["keep"].Value, "Should override because new is strictly newer")

	require.Contains(t, merged, "added")
	assert.Equal(t, "brandNew", merged["added"].Value)

	require.NotContains(t, merged, "remove")

	mergedSM, err := merge.Merge(context.Background(), oldMap, newMap, merge.MergeOptions{
		Mode: merge.ServerIsMaster,
	})
	require.NoError(t, err)
	require.NotNil(t, mergedSM)
	require.Contains(t, mergedSM, "remove")
}

func TestMapWithIntegerKeys(t *testing.T) {
	// Test maps with integer keys
	type IntKeyMapHolder struct {
		M map[int]string
	}

	oldMap := IntKeyMapHolder{
		M: map[int]string{
			1: "one",
			2: "two",
		},
	}

	newMap := IntKeyMapHolder{
		M: map[int]string{
			1: "ONE",
			3: "THREE",
		},
	}

	merged, err := merge.Merge(context.Background(), oldMap, newMap, merge.MergeOptions{
		Mode: merge.ClientIsMaster,
	})

	require.NoError(t, err)
	require.NotNil(t, merged.M)

	require.Len(t, merged.M, 2, "With ClientIsMaster, should only have entries from client")
	assert.Equal(t, "ONE", merged.M[1])
	assert.Equal(t, "THREE", merged.M[3])
	assert.NotContains(t, merged.M, 2)

	// With ServerIsMaster
	mergedSM, err := merge.Merge(context.Background(), oldMap, newMap, merge.MergeOptions{
		Mode: merge.ServerIsMaster,
	})

	require.NoError(t, err)
	require.NotNil(t, mergedSM.M)

	require.Len(t, mergedSM.M, 3, "With ServerIsMaster, should have entries from both")
	assert.Equal(t, "ONE", mergedSM.M[1], "Value from client should still override for keys in both")
	assert.Equal(t, "two", mergedSM.M[2], "Should keep server-only keys")
	assert.Equal(t, "THREE", mergedSM.M[3], "Should add client-only keys")
}

func TestNestedMaps(t *testing.T) {
	// Test nested maps
	type NestedMapHolder struct {
		OuterMap map[string]map[string]string
	}

	oldObj := NestedMapHolder{
		OuterMap: map[string]map[string]string{
			"map1": {
				"key1": "old_value1",
				"key2": "old_value2",
			},
			"map2": {
				"keyA": "old_valueA",
			},
		},
	}

	newObj := NestedMapHolder{
		OuterMap: map[string]map[string]string{
			"map1": {
				"key1": "new_value1",
				"key3": "new_value3",
			},
			"map3": {
				"keyX": "new_valueX",
			},
		},
	}

	// With ClientIsMaster
	merged, err := merge.Merge(context.Background(), oldObj, newObj, merge.MergeOptions{
		Mode: merge.ClientIsMaster,
	})

	require.NoError(t, err)
	require.NotNil(t, merged.OuterMap)

	// Check map1
	require.Contains(t, merged.OuterMap, "map1")
	require.Len(t, merged.OuterMap["map1"], 2, "Should only have client keys for map1")
	assert.Equal(t, "new_value1", merged.OuterMap["map1"]["key1"])
	assert.Equal(t, "new_value3", merged.OuterMap["map1"]["key3"])
	assert.NotContains(t, merged.OuterMap["map1"], "key2")

	// Check map3 exists and map2 doesn't
	require.Contains(t, merged.OuterMap, "map3")
	require.NotContains(t, merged.OuterMap, "map2")

	// With ServerIsMaster
	mergedSM, err := merge.Merge(context.Background(), oldObj, newObj, merge.MergeOptions{
		Mode: merge.ServerIsMaster,
	})

	require.NoError(t, err)
	require.NotNil(t, mergedSM.OuterMap)

	// Should have all maps
	require.Contains(t, mergedSM.OuterMap, "map1")
	require.Contains(t, mergedSM.OuterMap, "map2")
	require.Contains(t, mergedSM.OuterMap, "map3")

	// map1 should have all keys
	require.Len(t, mergedSM.OuterMap["map1"], 3)
	assert.Equal(t, "new_value1", mergedSM.OuterMap["map1"]["key1"])
	assert.Equal(t, "old_value2", mergedSM.OuterMap["map1"]["key2"])
	assert.Equal(t, "new_value3", mergedSM.OuterMap["map1"]["key3"])
}
