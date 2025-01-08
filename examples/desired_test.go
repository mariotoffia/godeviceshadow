package examples

import "testing"

func HandleDesiredReported(t *testing.T) {
	hub := HomeTemperatureHub{ // <1>
		DeviceShadow: DeviceShadow{
			TimeZone: "Europe/Stockholm",
			Owner:    "mariotoffia",
		},
		ClimateSensors: ClimateSensors{
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

	_ = hub

}
