package merge_test

import (
	"testing"
	"time"

	"github.com/mariotoffia/godeviceshadow/merge"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Sensor struct {
	ID        int
	TimeStamp time.Time
}

type Circuit struct {
	ID      int
	Sensors []Sensor
}

type Device struct {
	Name     string
	Circuits []Circuit
}

func (s *Sensor) GetTimestamp() time.Time {
	return s.TimeStamp
}
func (s *Sensor) SetTimestamp(t time.Time) {
	s.TimeStamp = t
}

// For map testing
type TimestampedMapVal struct {
	Value     string
	UpdatedAt time.Time
}

func (tmv *TimestampedMapVal) GetTimestamp() time.Time {
	return tmv.UpdatedAt
}

func (tmv *TimestampedMapVal) SetTimestamp(t time.Time) {
	tmv.UpdatedAt = t
}

func TestMergeOneNewerInSlice(t *testing.T) {
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
	mergedCircuit, err := merge.Merge(oldDevice, newDevice, merge.MergeOptions{
		Mode: merge.ClientIsMaster,
	})

	require.NoError(t, err)

	require.Equal(t, 1, len(mergedCircuit.Circuits))
	require.Equal(t, 3, len(mergedCircuit.Circuits[0].Sensors))
	assert.Equal(t, "new", mergedCircuit.Name, "Should override name from new device")
	assert.Equal(t, oneHourAgo, mergedCircuit.Circuits[0].Sensors[0].TimeStamp, "since new ts is older than old ts")
	assert.Equal(t, now, mergedCircuit.Circuits[0].Sensors[1].TimeStamp, "since new ts is newer than old ts")
	assert.Equal(t, oneHourAgo, mergedCircuit.Circuits[0].Sensors[2].TimeStamp, "since new and old ts are the same")
}

func TestDeleteWhenClientIsMaster(t *testing.T) {
	// oldDevice has 2 circuits, newDevice has 1 => the second circuit should be removed
	oldDevice := Device{
		Name: "device-old",
		Circuits: []Circuit{
			{ID: 1},
			{ID: 2},
		},
	}
	newDevice := Device{
		Name: "device-new",
		Circuits: []Circuit{
			{ID: 1}, // Only one circuit remains
		},
	}

	merged, err := merge.Merge(oldDevice, newDevice, merge.MergeOptions{
		Mode: merge.ClientIsMaster,
	})
	require.NoError(t, err)

	require.Len(t, merged.Circuits, 1, "Should remove the second circuit because client only has one")
	assert.Equal(t, 1, merged.Circuits[0].ID)
}

func TestRetainWhenServerIsMaster(t *testing.T) {
	// oldDevice has 2 circuits, newDevice has 1 => the second circuit should remain
	oldDevice := Device{
		Name: "device-old",
		Circuits: []Circuit{
			{ID: 1},
			{ID: 2},
		},
	}
	newDevice := Device{
		Name: "device-new",
		Circuits: []Circuit{
			{ID: 1}, // doesn't list circuit ID=2
		},
	}

	merged, err := merge.Merge(oldDevice, newDevice, merge.MergeOptions{
		Mode: merge.ServerIsMaster,
	})
	require.NoError(t, err)

	require.Len(t, merged.Circuits, 2, "Should keep second circuit because server is master")
	assert.Equal(t, 1, merged.Circuits[0].ID)
	assert.Equal(t, 2, merged.Circuits[1].ID)
}

func TestMergeDifferentLengthSlices(t *testing.T) {
	now := time.Now().UTC()
	oneHourAgo := now.Add(-1 * time.Hour)
	twoHoursAgo := now.Add(-2 * time.Hour)

	oldDevice := Device{
		Name: "OldDevice",
		Circuits: []Circuit{
			{ID: 1, Sensors: []Sensor{
				{ID: 10, TimeStamp: twoHoursAgo},
				{ID: 11, TimeStamp: twoHoursAgo},
				{ID: 12, TimeStamp: twoHoursAgo},
			}},
			{ID: 2, Sensors: []Sensor{
				{ID: 20, TimeStamp: twoHoursAgo},
			}},
		},
	}

	newDevice := Device{
		Name: "NewDevice",
		Circuits: []Circuit{
			{ID: 1, Sensors: []Sensor{
				{ID: 10, TimeStamp: oneHourAgo},  // new is newer => override
				{ID: 11, TimeStamp: twoHoursAgo}, // same => no update
				{ID: 12, TimeStamp: twoHoursAgo}, // same => no update
				{ID: 13, TimeStamp: now},         // brand new sensor
			}},
			{ID: 2, Sensors: []Sensor{}}, // explicitly empty
			{ID: 3, Sensors: []Sensor{ // entirely new circuit
				{ID: 30, TimeStamp: now},
			}},
		},
	}

	merged, err := merge.Merge(oldDevice, newDevice, merge.MergeOptions{
		Mode: merge.ClientIsMaster,
	})
	require.NoError(t, err)

	require.Len(t, merged.Circuits, 3)

	// Check circuit #1 sensors
	circuit1 := merged.Circuits[0]
	require.Len(t, circuit1.Sensors, 4, "We have an extra sensor ID=13 from new device")
	assert.Equal(t, oneHourAgo, circuit1.Sensors[0].TimeStamp, "ID=10 => updated (was twoHoursAgo, new is oneHourAgo => newer => override)")
	assert.Equal(t, twoHoursAgo, circuit1.Sensors[1].TimeStamp, "ID=11 => timestamps equal => keep old => but they're the same time, so no update anyway")
	assert.Equal(t, twoHoursAgo, circuit1.Sensors[2].TimeStamp, "ID=12 => same time => no update => keep old is effectively same")
	assert.Equal(t, now, circuit1.Sensors[3].TimeStamp, "ID=13 => brand new sensor")

	// Check circuit #2
	circuit2 := merged.Circuits[1]
	assert.Equal(t, 2, circuit2.ID)
	require.Len(t, circuit2.Sensors, 0)

	require.Len(t, merged.Circuits, 3)
	circuit3 := merged.Circuits[2]
	require.Len(t, circuit3.Sensors, 1, "One new sensor")
	assert.Equal(t, 30, circuit3.Sensors[0].ID)
	assert.Equal(t, now, circuit3.Sensors[0].TimeStamp)
}

func TestMergeMaps(t *testing.T) {
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

	merged, err := merge.Merge(oldObj, newObj, merge.MergeOptions{
		Mode: merge.ClientIsMaster,
	})
	require.NoError(t, err)
	require.NotNil(t, merged.M)

	require.Contains(t, merged.M, "keep")
	assert.Equal(t, "new", merged.M["keep"].Value, "Should override because new is strictly newer")

	require.Contains(t, merged.M, "added")
	assert.Equal(t, "brandNew", merged.M["added"].Value)

	require.NotContains(t, merged.M, "remove")

	mergedSM, err := merge.Merge(oldObj, newObj, merge.MergeOptions{
		Mode: merge.ServerIsMaster,
	})

	require.NoError(t, err)
	require.NotNil(t, mergedSM.M)

	require.Contains(t, mergedSM.M, "remove")
}

func TestMergeMaps2(t *testing.T) {
	now := time.Now().UTC()
	oldMap := map[string]TimestampedMapVal{
		"keep":   {Value: "old", UpdatedAt: now.Add(-2 * time.Hour)},
		"remove": {Value: "willRemove", UpdatedAt: now.Add(-1 * time.Hour)},
	}
	newMap := map[string]TimestampedMapVal{
		"keep":  {Value: "new", UpdatedAt: now.Add(-1 * time.Hour)}, // old is older => replaced if newer
		"added": {Value: "brandNew", UpdatedAt: now},                // brand new key
	}

	merged, err := merge.Merge(oldMap, newMap, merge.MergeOptions{
		Mode: merge.ClientIsMaster,
	})
	require.NoError(t, err)
	require.NotNil(t, merged)

	require.Contains(t, merged, "keep")
	assert.Equal(t, "new", merged["keep"].Value, "Should override because new is strictly newer")

	require.Contains(t, merged, "added")
	assert.Equal(t, "brandNew", merged["added"].Value)

	require.NotContains(t, merged, "remove")

	mergedSM, err := merge.Merge(oldMap, newMap, merge.MergeOptions{
		Mode: merge.ServerIsMaster,
	})
	require.NoError(t, err)
	require.NotNil(t, mergedSM)
	require.Contains(t, mergedSM, "remove")
}

func TestEmptySlices(t *testing.T) {
	oldDevice := Device{
		Name:     "OldCorner",
		Circuits: []Circuit{},
	}
	newDevice := Device{
		Name:     "NewCorner",
		Circuits: nil, // explicitly nil
	}

	merged, err := merge.Merge(oldDevice, newDevice, merge.MergeOptions{
		Mode: merge.ClientIsMaster,
	})
	require.NoError(t, err)

	// new device slice is nil => in ClientIsMaster => remove old => result should have no circuits
	require.Empty(t, merged.Circuits, "Should remove old slices since new is nil and client is master")

	// If server is master, we keep the old even if new is nil
	mergedSM, err := merge.Merge(oldDevice, newDevice, merge.MergeOptions{
		Mode: merge.ServerIsMaster,
	})
	require.NoError(t, err)
	require.Len(t, mergedSM.Circuits, 0, "But old was empty anyway, so effectively the same in this corner case.")
}

func TestEqualTimestampNoUpdate(t *testing.T) {
	t0 := time.Now().UTC()

	oldDevice := Device{
		Name: "OldEq",
		Circuits: []Circuit{
			{
				ID: 1,
				Sensors: []Sensor{
					{ID: 10, TimeStamp: t0},
				},
			},
		},
	}
	newDevice := Device{
		Name: "NewEq",
		Circuits: []Circuit{
			{
				ID: 1,
				Sensors: []Sensor{
					{ID: 10, TimeStamp: t0}, // same time
				},
			},
		},
	}

	merged, err := merge.Merge(oldDevice, newDevice, merge.MergeOptions{
		Mode:                merge.ClientIsMaster,
		DoOverrideWithEmpty: true,
	})
	require.NoError(t, err)

	require.Len(t, merged.Circuits, 1)
	require.Len(t, merged.Circuits[0].Sensors, 1)
	assert.Equal(t, t0, merged.Circuits[0].Sensors[0].TimeStamp, "Timestamps are equal => no update => keep old (though they look the same in practice).")
}

func BenchmarkMergeOneNewerInSlice(t *testing.B) {
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

	t.ResetTimer()

	for i := 0; i < t.N; i++ {
		merge.Merge(oldDevice, newDevice, merge.MergeOptions{
			Mode: merge.ClientIsMaster,
		})
	}
}
