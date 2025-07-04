package examples

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/mariotoffia/godeviceshadow/merge"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDesiredReported(t *testing.T) {
	// Initial device shadow state of the reported (e.g. from db)
	reported := HomeTemperatureHub{
		MetaInfo: &MetaInfo{
			TimeZone: "Europe/Stockholm",
			Owner:    "mariotoffia",
		},
		ClimateSensors: &ClimateSensors{
			Outdoor: map[string]OutdoorTemperatureSensor{
				"garden": {
					Direction:   DirectionNorth,
					Temperature: -27.0,
					UpdatedAt:   parse("2023-01-01T12:00:00+01:00"),
					Humidity:    17.0,
				},
			},
			Indoor: map[string]IndoorTemperatureSensor{
				"living_room": {
					Floor:       1,
					Direction:   DirectionNorth,
					Temperature: 22.6,
					UpdatedAt:   parse("2023-01-01T11:55:00+01:00"),
					Humidity:    32.0,
				},
				"basement": {
					Floor:       0,
					Direction:   DirectionSouth,
					Temperature: 18.0,
					UpdatedAt:   parse("2023-01-01T11:00:00+01:00"),
					Humidity:    40.0,
				},
			},
		},
	}

	// Initial desired state of the hub (e.g. from db)
	desired := HomeTemperatureHub{}

	var err error

	// Simulate new actuation -> plain merge
	desired, err = merge.Merge(desired, HomeTemperatureHub{
		IndoorTempSP: &IndoorTemperatureSetPoint{
			SetPoint:  22.0,
			UpdatedAt: parse("2023-01-01T13:00:00+01:00"),
		},
	}, merge.MergeOptions{})
	require.NoError(t, err)
	require.Equal(t, 22.0, desired.IndoorTempSP.SetPoint)

	data, _ := json.Marshal(desired)
	fmt.Println(string(data))
	// Output:
	// {"indoor_temp_sp":{"sp":22,"ts":"2023-01-01T13:00:00+01:00"}}

	// Report back to the device shadow
	reported, err = merge.Merge(reported, HomeTemperatureHub{
		IndoorTempSP: &IndoorTemperatureSetPoint{
			SetPoint: 22.0,
			// Must be added or newer ts than the "old" reported
			UpdatedAt: parse("2023-01-01T13:05:00+01:00"),
		},
	}, merge.MergeOptions{
		Mode: merge.ServerIsMaster,
	})

	require.NoError(t, err)
	require.Equal(t, 22.0, reported.IndoorTempSP.SetPoint)

	// Acknowledge in the desired model -> removed from model
	desired, err = merge.Desired(reported, desired, merge.DesiredOptions{})
	require.NoError(t, err)
	// Check that the indoor temp setpoint is either nil or has zero values
	if desired.IndoorTempSP != nil {
		assert.Equal(t, 0.0, desired.IndoorTempSP.SetPoint, "SetPoint should be zero value after acknowledgement")
		assert.Equal(t, time.Time{}, desired.IndoorTempSP.UpdatedAt, "UpdatedAt should be zero value after acknowledgement")
	}

	data, _ = json.Marshal(desired)
	fmt.Println(string(data))
	// The output may contain an empty IndoorTempSP with zero values
	// {"indoor_temp_sp":{"sp":0,"ts":"0001-01-01T00:00:00Z"}}
	// or it may be empty if IndoorTempSP is completely removed
	// {}

	data, _ = json.Marshal(reported)
	fmt.Println(string(data))
	// Output:
	// {
	//   "meta": {"tz": "Europe/Stockholm", "owner": "mariotoffia"},
	//   "climate": {
	//     "outdoor": {
	//       "garden": {
	//         "direction": "north",
	//         "t": -27,
	//         "h": 17,
	//         "ts": "2023-01-01T12:00:00+01:00"
	//       }
	//     },
	//     "indoor": {
	//       "basement": {
	//         "floor": 0,
	//         "direction": "south",
	//         "t": 18,
	//         "h": 40,
	//         "ts": "2023-01-01T11:00:00+01:00"
	//       },
	//       "living_room": {
	//         "floor": 1,
	//         "direction": "north",
	//         "t": 22.6,
	//         "h": 32,
	//         "ts": "2023-01-01T11:55:00+01:00"
	//       }
	//     }
	//   },
	//   "indoor_temp_sp": {"sp": 22, "ts": "2023-01-01T13:00:00+01:00"}
	// }
}
