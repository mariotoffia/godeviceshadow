package merge_test

import (
	"fmt"
	"time"

	"github.com/mariotoffia/godeviceshadow/model"
)

// Common test types used across multiple test files

type Sensor struct {
	ID        int
	TimeStamp time.Time
	Value     any
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
func (s *Sensor) GetValue() any {
	return s.Value
}

// For map testing
type TimestampedMapVal struct {
	Value     string
	UpdatedAt time.Time
}

func (tmv *TimestampedMapVal) GetTimestamp() time.Time {
	return tmv.UpdatedAt
}
func (tmv *TimestampedMapVal) GetValue() any {
	return tmv.Value
}

// IdSensor implements model.IdValueAndTimestamp for ID-based merging tests
type IdSensor struct {
	ID        string
	TimeStamp time.Time
	Value     float64
}

func (s *IdSensor) GetTimestamp() time.Time {
	return s.TimeStamp
}

func (s *IdSensor) GetValue() any {
	return s.Value
}

func (s *IdSensor) GetID() string {
	return s.ID
}

// CustomMergeable implements the model.Merger interface
type CustomMergeable struct {
	Name      string
	Value     int
	Timestamp time.Time
}

func (c *CustomMergeable) Merge(other any, mode model.MergeMode) (any, error) {
	otherMergeable, ok := other.(*CustomMergeable)
	if !ok {
		return nil, fmt.Errorf("expected *CustomMergeable, got %T", other)
	}

	result := &CustomMergeable{
		Name:      c.Name + "+" + otherMergeable.Name,
		Timestamp: time.Now().UTC(),
	}

	// Merge based on mode
	if mode == model.ClientIsMaster {
		result.Value = otherMergeable.Value
	} else {
		// If timestamp of other is newer, use its value
		if otherMergeable.Timestamp.After(c.Timestamp) {
			result.Value = otherMergeable.Value
		} else {
			result.Value = c.Value
		}
	}

	return result, nil
}

// ErrorMergeable implements Merger but returns an error when merging
type ErrorMergeable struct {
	ShouldError bool
}

func (e *ErrorMergeable) Merge(other any, mode model.MergeMode) (any, error) {
	if e.ShouldError {
		return nil, fmt.Errorf("simulated error from custom merger")
	}
	return e, nil
}

// MockIdValueType is a struct that requires an explicit pointer
// to implement IdValueAndTimestamp
type MockIdValueType struct {
	ID        string
	TimeStamp time.Time
	Value     string
}

func (m *MockIdValueType) GetTimestamp() time.Time {
	return m.TimeStamp
}

func (m *MockIdValueType) GetValue() any {
	return m.Value
}

func (m *MockIdValueType) GetID() string {
	return m.ID
}

// SensorContainer is a struct that contains ID-based sensors
type SensorContainer struct {
	Name    string
	Sensors []IdSensor
}
