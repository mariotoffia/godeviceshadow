package merge

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/mariotoffia/godeviceshadow/utils/vtsutils"
)

// DesiredOptions holds configuration for the desired state processing.
type DesiredOptions struct {
	// Loggers will be notified on add, updated, remove, not-changed operations while merging.
	// When a value in the reported state matches a value in the desired state, the logger
	// will be notified via NotifyAcknowledge with the path and value.
	Loggers DesiredLoggers

	// ContinueOnError determines whether to continue processing on error or stop immediately.
	// When true, errors are collected and returned at the end. When false (default),
	// processing stops on the first error.
	// This is useful for batch processing where you want to collect all errors rather than
	// stopping at the first one.
	ContinueOnError bool
}

type DesiredObject struct {
	DesiredOptions
	CurrentPath string
	Errors      DesiredErrors
}

// Desired is a special merge where a reported model is analyzed if it matches the desired model.
// All matched values are removed or set to `nil` in the desired model. The result is the desired
// model with the matched values removed.
//
// It works by comparing values that implement the `model.ValueAndTimestamp` interface:
//   - When values in both models are equal (via `vtsutils.Equals`), the value is removed from the desired model
//   - Non-matching values remain in the desired model
//   - Loggers are notified of acknowledged values via `NotifyAcknowledge`
//
// For complex data structures:
//   - Structs: Each field is processed recursively
//   - Maps: Each key-value pair is processed recursively
//   - Slices: Elements are matched by ID (if they implement IdValueAndTimestamp) or by position
//
// The intended usage pattern is:
//  1. Device sends a reported state
//  2. Compare with current desired state using this function
//  3. Result contains only values that haven't been acknowledged yet
//  4. Send this filtered desired state back to the device
//
// Example:
//
//	reported := MyDevice{
//	    Temperature: MyValueTS{Value: 72.5, Timestamp: time.Now()},
//	    Humidity: MyValueTS{Value: 40.0, Timestamp: time.Now()},
//	}
//
//	desired := MyDevice{
//	    Temperature: MyValueTS{Value: 72.5, Timestamp: time.Now()}, // Matches reported
//	    Humidity: MyValueTS{Value: 45.0, Timestamp: time.Now()},    // Doesn't match
//	    FanSpeed: MyValueTS{Value: "high", Timestamp: time.Now()},  // Not in reported
//	}
//
//	// Configure options
//	opts := merge.DesiredOptions{
//	    Loggers: merge.DesiredLoggers{myLogger},
//	    ContinueOnError: true,
//	}
//
//	// Process desired state
//	result, err := merge.Desired(reported, desired, opts)
//
//	// result will have:
//	// - Temperature: nil (removed because it matched)
//	// - Humidity: MyValueTS{Value: 45.0, ...} (kept because it didn't match)
//	// - FanSpeed: MyValueTS{Value: "high", ...} (kept because it wasn't in reported)
func Desired[T any](ctx context.Context, reportedModel, desiredModel T, opts DesiredOptions) (T, error) {
	//
	mergedVal, err := DesiredAny(ctx, reportedModel, desiredModel, opts)

	var zero T

	if err != nil {
		return zero, err
	}

	return mergedVal.(T), nil
}

func DesiredAny(ctx context.Context, reportedModel, desiredModel any, opts DesiredOptions) (any, error) {
	//
	// Handle nil cases
	if reportedModel == nil && desiredModel == nil {
		return nil, nil
	}

	// Handle nil reportedModel specifically
	if reportedModel == nil {
		return nil, nil
	}

	// Handle nil desiredModel
	if desiredModel == nil {
		return nil, fmt.Errorf("desired model cannot be nil when reported model is non-nil")
	}

	reportedVal := reflect.ValueOf(reportedModel)
	desiredVal := reflect.ValueOf(desiredModel)

	// Handle nil pointers and interfaces more specifically
	if (reportedVal.Kind() == reflect.Ptr || reportedVal.Kind() == reflect.Interface) && reportedVal.IsNil() {
		return nil, nil
	}

	// Ensure both values are valid
	if !reportedVal.IsValid() || !desiredVal.IsValid() {
		return nil, nil
	}

	// Check for kind mismatch
	if reportedVal.Kind() != desiredVal.Kind() {
		return nil, fmt.Errorf("reported and desired model must be of the same kind: %s != %s", reportedVal.Kind(), desiredVal.Kind())
	}

	// Create desired object with error tracking
	desiredObj := DesiredObject{
		DesiredOptions: opts,
		Errors:         make(DesiredErrors, 0),
	}

	// Process the desired model recursively
	result := desiredRecursive(ctx, reportedVal, desiredVal, desiredObj)

	// Check if we have any errors
	if len(desiredObj.Errors) > 0 && !opts.ContinueOnError {
		return nil, desiredObj.Errors
	}

	if !result.IsValid() {
		return nil, nil
	}

	// Make sure the result can be safely converted to an interface
	if result.CanInterface() {
		return result.Interface(), nil
	}

	return nil, nil
}

func desiredRecursive(ctx context.Context, reportedVal, desiredVal reflect.Value, obj DesiredObject) reflect.Value {
	// Safety check for invalid values
	if !reportedVal.IsValid() || !desiredVal.IsValid() {
		return reflect.Value{}
	}

	// Handle nil pointers and interfaces
	if canBeNil(reportedVal) && reportedVal.IsNil() {
		return desiredVal // Keep desired value unchanged if reported is nil
	}

	if canBeNil(desiredVal) && desiredVal.IsNil() {
		return desiredVal // Keep nil desired value
	}

	// Special handling for pointer types
	if reportedVal.Kind() == reflect.Ptr && desiredVal.Kind() == reflect.Ptr {
		// Get the elements they point to
		reportedElem := reportedVal.Elem()
		desiredElem := desiredVal.Elem()

		// Process the elements
		if reportedElem.IsValid() && desiredElem.IsValid() {
			result := desiredRecursive(ctx, reportedElem, desiredElem, obj)

			if result.IsValid() {
				// Create a new pointer to the result
				newPtr := reflect.New(result.Type())
				newPtr.Elem().Set(result)
				return newPtr
			}
		}

		return desiredVal
	}

	// If we can't set the desired value, make it addressable
	if !desiredVal.CanSet() {
		desiredVal = makeAddressable(desiredVal)
	}

	// If both implement ValueAndTimestamp, check for equality
	if rvt, ok := unwrapValueAndTimestamp(reportedVal); ok {
		if dvt, ok := unwrapValueAndTimestamp(desiredVal); ok {
			if vtsutils.Equals(rvt, dvt) {
				// Safely notify about the acknowledgment
				if obj.Loggers != nil {
					obj.Loggers.NotifyAcknowledge(ctx, obj.CurrentPath, rvt)
				}

				// Remove from desired model
				return reflect.Zero(desiredVal.Type())
			} else {
				// Values don't match - could record this as diagnostic info
				_ = recordError(&obj, "Values don't match", reportedVal, desiredVal)
			}
		} else {
			// Desired doesn't implement ValueAndTimestamp
			_ = recordError(&obj, "Desired value doesn't implement ValueAndTimestamp", reportedVal, desiredVal)
		}

		return desiredVal // Keep desired value if not equal
	}

	// Unwrap values to work with concrete values
	reportedVal = unwrapReflectValue(reportedVal)
	desiredVal = unwrapReflectValue(desiredVal)

	// Safety check after unwrapping
	if !reportedVal.IsValid() || !desiredVal.IsValid() {
		return desiredVal
	}

	basePath := obj.CurrentPath

	switch reportedVal.Kind() {
	case reflect.Struct:
		for i := 0; i < reportedVal.NumField(); i++ {
			field := reportedVal.Type().Field(i)

			if field.PkgPath != "" {
				continue // Unexported field -> skip
			}

			tag := getJSONTag(field)

			if tag == "" {
				continue // No tag -> skip
			}

			obj.CurrentPath = concatPath(basePath, tag)

			if r := desiredRecursive(ctx, reportedVal.Field(i), desiredVal.Field(i), obj); r.IsValid() {
				desiredVal.Field(i).Set(r)
			}
		}
	case reflect.Map:
		// Safety check for nil maps
		if reportedVal.IsNil() || desiredVal.IsNil() {
			return desiredVal
		}

		for _, key := range reportedVal.MapKeys() {
			// Skip invalid keys
			if !key.IsValid() {
				continue
			}

			obj.CurrentPath = concatPath(basePath, formatKey(key))

			// Get map values safely
			reportedMapVal := reportedVal.MapIndex(key)
			desiredMapVal := desiredVal.MapIndex(key)

			// Skip if reported value is invalid
			if !reportedMapVal.IsValid() {
				continue
			}

			// If key exists in desired, process it
			if desiredMapVal.IsValid() {
				result := desiredRecursive(ctx, reportedMapVal, desiredMapVal, obj)

				// Update or remove map key based on result
				if result.IsValid() && !result.IsZero() {
					// Update key with the new value
					desiredVal.SetMapIndex(key, result)
				} else {
					// Remove key from the map
					desiredVal.SetMapIndex(key, reflect.Value{}) // This deletes the key
				}
			}
		}
	case reflect.Slice, reflect.Array:
		// Safety check for nil slices
		if reportedVal.Kind() == reflect.Slice && reportedVal.IsNil() {
			return desiredVal
		}
		if desiredVal.Kind() == reflect.Slice && desiredVal.IsNil() {
			return desiredVal
		}

		// Try ID-based matching first if elements implement IdValueAndTimestamp
		idBasedResult := desiredSliceById(ctx, reportedVal, desiredVal, obj)
		if idBasedResult.IsValid() {
			return idBasedResult
		}

		// Fall back to positional matching
		minLen := reportedVal.Len()

		if minLen > desiredVal.Len() {
			minLen = desiredVal.Len()
		}

		for i := 0; i < minLen; i++ {
			obj.CurrentPath = fmt.Sprintf("%s.%d", basePath, i)

			reportedItem := reportedVal.Index(i)
			desiredItem := desiredVal.Index(i)

			// Skip invalid items
			if !reportedItem.IsValid() || !desiredItem.IsValid() {
				continue
			}

			// Process items that can be set
			if desiredItem.CanSet() {
				if r := desiredRecursive(ctx, reportedItem, desiredItem, obj); r.IsValid() {
					desiredItem.Set(r)
				}
			}
		}
	}

	return desiredVal
}

// desiredSliceById attempts to match slice elements by ID if they implement IdValueAndTimestamp
// This function enhances the device shadow pattern by allowing slices to be matched by ID rather than position.
// For example, if you have a list of sensors in both reported and desired state, they can be matched by their unique IDs.
//
// Returns the processed desired slice value if successful, or an invalid reflect.Value if ID-based matching couldn't be applied.
func desiredSliceById(ctx context.Context, reportedVal, desiredVal reflect.Value, obj DesiredObject) reflect.Value {
	basePath := obj.CurrentPath

	// Create maps to track elements by ID
	reportedMap := make(map[string]int) // ID -> index in reportedVal
	processed := make(map[string]bool)  // IDs that have been processed

	// Extract IDs from reported elements
	for i := 0; i < reportedVal.Len(); i++ {
		elem := reportedVal.Index(i)
		if idvt, ok := unwrapIdValueAndTimestamp(elem); ok {
			id := idvt.GetID()
			reportedMap[id] = i
		} else {
			// If any element doesn't implement IdValueAndTimestamp, fall back to positional matching
			return reflect.Value{}
		}
	}

	// Process each element in the desired slice
	for i := 0; i < desiredVal.Len(); i++ {
		desiredElem := desiredVal.Index(i)

		// Skip invalid elements
		if !desiredElem.IsValid() || !desiredElem.CanSet() {
			continue
		}

		// Extract ID from desired element
		desiredIdVt, ok := unwrapIdValueAndTimestamp(desiredElem)
		if !ok {
			// If any element doesn't implement IdValueAndTimestamp, fall back to positional matching
			return reflect.Value{}
		}

		desiredId := desiredIdVt.GetID()
		processed[desiredId] = true

		// Find matching element in reported slice
		if reportedIdx, exists := reportedMap[desiredId]; exists {
			reportedElem := reportedVal.Index(reportedIdx)
			obj.CurrentPath = fmt.Sprintf("%s.%s", basePath, desiredId)

			// Process the elements recursively
			if r := desiredRecursive(ctx, reportedElem, desiredElem, obj); r.IsValid() {
				desiredElem.Set(r)
			}
		}
	}

	return desiredVal
}

func makeAddressable(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Ptr {
		return v
	}
	ptr := reflect.New(v.Type())
	ptr.Elem().Set(v)
	return ptr.Elem()
}

// canBeNil checks if a reflect.Value can be nil.
// Some reflect.Value kinds (like basic types) can never be nil, while others (like pointers, interfaces,
// maps, slices, etc.) can be nil. This function helps determine if the IsNil check is valid for a given value.
func canBeNil(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return true
	}
	return false
}

// DesiredError represents an error that occurred during the desired operation.
// This provides detailed context about where and why a desired operation failed,
// which can be helpful for debugging complex device shadow synchronization issues.
type DesiredError struct {
	// Path is the JSON path where the error occurred
	Path string
	// Message is the error message
	Message string
	// ReportedValue is the value from the reported model
	ReportedValue any
	// DesiredValue is the value from the desired model
	DesiredValue any
}

// Error implements the error interface
func (e DesiredError) Error() string {
	return fmt.Sprintf("Error at path '%s': %s", e.Path, e.Message)
}

// DesiredErrors is a collection of errors that occurred during processing
type DesiredErrors []DesiredError

// Error implements the error interface for DesiredErrors
func (e DesiredErrors) Error() string {
	if len(e) == 0 {
		return "No errors"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d errors occurred during desired operation:\n", len(e)))
	for i, err := range e {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, err.Error()))
	}
	return sb.String()
}

// recordError adds an error to the DesiredObject's error list
// Returns false if processing should stop, true if it should continue
func recordError(obj *DesiredObject, message string, reportedVal, desiredVal reflect.Value) bool {
	if obj == nil {
		return false
	}

	var reportedValue, desiredValue any

	if reportedVal.IsValid() && reportedVal.CanInterface() {
		reportedValue = reportedVal.Interface()
	}

	if desiredVal.IsValid() && desiredVal.CanInterface() {
		desiredValue = desiredVal.Interface()
	}

	err := DesiredError{
		Path:          obj.CurrentPath,
		Message:       message,
		ReportedValue: reportedValue,
		DesiredValue:  desiredValue,
	}

	obj.Errors = append(obj.Errors, err)

	return obj.ContinueOnError
}
