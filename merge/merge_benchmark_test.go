package merge_test

import (
	"context"
	"testing"
	"time"

	"github.com/mariotoffia/godeviceshadow/merge"
)

// BenchmarkMergeSimple benchmarks simple merging operations
func BenchmarkMergeSimple(b *testing.B) {
	now := time.Now().UTC()
	oneHourAgo := now.Add(-1 * time.Hour)

	oldDevice := Device{Name: "old",
		Circuits: []Circuit{
			{ID: 1, Sensors: []Sensor{
				{ID: 1, TimeStamp: oneHourAgo},
				{ID: 2, TimeStamp: oneHourAgo},
				{ID: 3, TimeStamp: oneHourAgo},
			},
			},
		},
	}

	newDevice := Device{Name: "new",
		Circuits: []Circuit{
			{ID: 1, Sensors: []Sensor{
				{ID: 1, TimeStamp: oneHourAgo.Add(-1 * time.Hour)},
				{ID: 2, TimeStamp: now},
				{ID: 3, TimeStamp: oneHourAgo},
			},
			},
		},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = merge.Merge(context.Background(), oldDevice, newDevice, merge.MergeOptions{
			Mode: merge.ClientIsMaster,
		})
	}
}

// BenchmarkMergeMap benchmarks map merging operations
func BenchmarkMergeMap(b *testing.B) {
	now := time.Now().UTC()

	oldMap := map[string]TimestampedMapVal{
		"keep1":   {Value: "old1", UpdatedAt: now.Add(-2 * time.Hour)},
		"keep2":   {Value: "old2", UpdatedAt: now.Add(-2 * time.Hour)},
		"remove1": {Value: "willRemove1", UpdatedAt: now.Add(-1 * time.Hour)},
		"remove2": {Value: "willRemove2", UpdatedAt: now.Add(-1 * time.Hour)},
	}

	newMap := map[string]TimestampedMapVal{
		"keep1":  {Value: "new1", UpdatedAt: now.Add(-1 * time.Hour)},
		"keep2":  {Value: "new2", UpdatedAt: now.Add(-1 * time.Hour)},
		"added1": {Value: "brandNew1", UpdatedAt: now},
		"added2": {Value: "brandNew2", UpdatedAt: now},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = merge.Merge(context.Background(), oldMap, newMap, merge.MergeOptions{
			Mode: merge.ClientIsMaster,
		})
	}
}

// BenchmarkMergeSliceByID benchmarks ID-based slice merging
func BenchmarkMergeSliceByID(b *testing.B) {
	now := time.Now().UTC()
	oneHourAgo := now.Add(-1 * time.Hour)

	// Create slices for ID-based merging
	oldSlice := []IdSensor{
		{ID: "temp1", TimeStamp: oneHourAgo, Value: 22.5},
		{ID: "temp2", TimeStamp: oneHourAgo, Value: 18.0},
		{ID: "temp3", TimeStamp: oneHourAgo, Value: 25.0},
	}

	newSlice := []IdSensor{
		{ID: "temp3", TimeStamp: now, Value: 26.0},        // Updated
		{ID: "temp4", TimeStamp: now, Value: 21.0},        // New
		{ID: "temp1", TimeStamp: oneHourAgo, Value: 22.5}, // Same
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = merge.Merge(context.Background(), oldSlice, newSlice, merge.MergeOptions{
			Mode:            merge.ClientIsMaster,
			MergeSlicesByID: true,
		})
	}
}

// BenchmarkCustomMerger benchmarks custom merger implementation
func BenchmarkCustomMerger(b *testing.B) {
	now := time.Now().UTC()
	oneHourAgo := now.Add(-1 * time.Hour)

	oldObj := &CustomMergeable{
		Name:      "Base",
		Value:     10,
		Timestamp: oneHourAgo,
	}

	newObj := &CustomMergeable{
		Name:      "Override",
		Value:     20,
		Timestamp: now,
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = merge.Merge(context.Background(), oldObj, newObj, merge.MergeOptions{
			Mode: merge.ClientIsMaster,
		})
	}
}
