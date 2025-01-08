package examples

import (
	"fmt"
	"testing"

	"github.com/mariotoffia/godeviceshadow/loggers/strlogger"
	"github.com/mariotoffia/godeviceshadow/merge"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMerge(t *testing.T) {
	hubZero := HomeTemperatureHub{ // <1>
		MetaInfo: &MetaInfo{
			TimeZone: "Europe/Stockholm",
			Owner:    "mariotoffia",
		},
		ClimateSensors: &ClimateSensors{
			Indoor: map[string]IndoorTemperatureSensor{
				"living_room": {
					Floor:       1,
					Direction:   DirectionNorth,
					Temperature: 22.6,
					UpdatedAt:   parse("2023-01-01T10:00:00+01:00"),
					Humidity:    32.0,
				},
			},
		},
	}

	hub := HomeTemperatureHub{ // <2>
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

	sl := strlogger.New() // <3>

	res, err := merge.Merge(hubZero, hub, merge.MergeOptions{
		Mode:    merge.ServerIsMaster,   // <4>
		Loggers: merge.MergeLoggers{sl}, // <5>
	})
	require.NoError(t, err)
	assert.Equal(t, hub, res) // <6>

	fmt.Println(sl.String())
	// Outputs:
	// Operation  Old Timestamp              New Timestamp              Path                           OldValue                                 NewValue
	// noop                                                             shadow.tz                      Europe/Stockholm                         Europe/Stockholm
	// noop                                                             shadow.owner                   mariotoffia                              mariotoffia
	// add        0001-01-01T00:00:00Z       2023-01-01T12:00:00+01:00  climate.outdoor.garden         nil                                      map[direction:north humidity:17 temperature:-27]
	// update     2023-01-01T10:00:00+01:00  2023-01-01T11:55:00+01:00  climate.indoor.living_room     map[direction:north floor:1 humidity:32 temperature:22.6] map[direction:north floor:1 humidity:32 temperature:22.6]
	// add        0001-01-01T00:00:00Z       2023-01-01T11:00:00+01:00  climate.indoor.basement        nil                                      map[direction:south floor:0 humidity:40 temperature:18]

}
