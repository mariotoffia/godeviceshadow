package vtsutils

import (
	"reflect"

	"github.com/mariotoffia/godeviceshadow/model"
)

// ValueAndTimeStampEqual compares _a_, _b_ using `ValueAndTimestamp.GetValue()` function.
//
// It will check if they are same type, if not, it will return `false`. It it is a scalar
// value it will just perform == check. If it is a map[string]any it will iterate over
// all keys and compare the values. If all values exists in both _a_ and _b_ and are equal
// it will return `true`, otherwise `false`. If any of the _a_ or _b_ map do miss a key,
// it will return `false`.
//
// All other, it will do a reflect.DeepEqual(a, b) and return the result.
func ValueAndTimeStampEqual(a, b model.ValueAndTimestamp) bool {
	// Check for nil values
	if a == nil || b == nil {
		return a == b
	}

	// Retrieve values
	valueA := a.GetValue()
	valueB := b.GetValue()

	// Check if types are the same
	if reflect.TypeOf(valueA) != reflect.TypeOf(valueB) {
		return false
	}

	// Handle scalar values
	if reflect.TypeOf(valueA).Kind() != reflect.Map && reflect.TypeOf(valueA).Kind() != reflect.Slice {
		return valueA == valueB
	}

	// Handle map comparison
	if mapA, ok := valueA.(map[string]any); ok {
		mapB, ok := valueB.(map[string]any)
		if !ok {
			return false
		}
		// Check keys and values
		if len(mapA) != len(mapB) {
			return false
		}
		for key, valA := range mapA {
			valB, exists := mapB[key]
			if !exists {
				return false
			}
			if !compareValues(valA, valB) {
				return false
			}
		}
		return true
	}

	return reflect.DeepEqual(valueA, valueB)
}

// compareValues is a helper function to handle comparison of scalar values,
// maps, slices, and recursive ValueAndTimestamp.
func compareValues(valA, valB any) bool {
	if reflect.TypeOf(valA) != reflect.TypeOf(valB) {
		return false
	}

	if scalarVal, ok := valA.(model.ValueAndTimestamp); ok {
		return ValueAndTimeStampEqual(scalarVal, valB.(model.ValueAndTimestamp))
	}

	return valA == valB
}
