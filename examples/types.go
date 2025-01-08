package examples

import (
	"time"
)

type HomeTemperatureHub struct {
	*MetaInfo      `json:"meta,omitempty"`
	ClimateSensors *ClimateSensors            `json:"climate,omitempty"`
	IndoorTempSP   *IndoorTemperatureSetPoint `json:"indoor_temp_sp,omitempty"` // Important omitempty when used in desired
}

type MetaInfo struct {
	TimeZone string `json:"tz,omitempty"`
	Owner    string `json:"owner,omitempty"`
}

type Direction string

const (
	DirectionNorth Direction = "north"
	DirectionSouth Direction = "south"
	DirectionEast  Direction = "east"
	DirectionWest  Direction = "west"
)

type IndoorTemperatureSensor struct {
	Floor       int       `json:"floor"`
	Direction   Direction `json:"direction"`
	Temperature float64   `json:"t"`
	Humidity    float64   `json:"h"`
	UpdatedAt   time.Time `json:"ts"`
}

func (idt *IndoorTemperatureSensor) GetTimestamp() time.Time {
	return idt.UpdatedAt
}

func (idt *IndoorTemperatureSensor) GetValue() any {
	return map[string]any{
		"floor":       idt.Floor,
		"direction":   idt.Direction,
		"temperature": idt.Temperature,
		"humidity":    idt.Humidity,
	}
}

type OutdoorTemperatureSensor struct {
	Direction   Direction `json:"direction"`
	Temperature float64   `json:"t"`
	Humidity    float64   `json:"h"`
	UpdatedAt   time.Time `json:"ts"`
}

func (ots *OutdoorTemperatureSensor) GetTimestamp() time.Time {
	return ots.UpdatedAt
}

func (ots *OutdoorTemperatureSensor) GetValue() any {
	return map[string]any{
		"direction":   ots.Direction,
		"temperature": ots.Temperature,
		"humidity":    ots.Humidity,
	}
}

type IndoorTemperatureSetPoint struct {
	SetPoint  float64   `json:"sp"`
	UpdatedAt time.Time `json:"ts"`
}

func (sp *IndoorTemperatureSetPoint) GetTimestamp() time.Time {
	return sp.UpdatedAt
}

func (sp *IndoorTemperatureSetPoint) GetValue() any {
	return map[string]any{
		"sp": sp.SetPoint,
	}
}

type ClimateSensors struct {
	Outdoor map[string]OutdoorTemperatureSensor `json:"outdoor,omitempty"`
	Indoor  map[string]IndoorTemperatureSensor  `json:"indoor,omitempty"`
}
