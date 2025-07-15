package merge_test

import (
	"context"
	"testing"

	"github.com/mariotoffia/godeviceshadow/merge"
	"github.com/stretchr/testify/assert"
)

func TestDesiredAnyPanicCases(t *testing.T) {
	t.Run("EdgeCase_NilReportedNonNilDesired", func(t *testing.T) {
		// This test specifically targets the case where reported is nil but desired is not
		var reported interface{} = nil
		desired := "not nil"

		// This should not panic and now returns nil with no error
		result, err := merge.DesiredAny(context.Background(), reported, desired, merge.DesiredOptions{})
		assert.NoError(t, err, "No error expected when reported is nil")
		assert.Nil(t, result, "Expected nil result")
	})

	t.Run("EdgeCase_NonNilReportedNilDesired", func(t *testing.T) {
		// This test specifically targets the case where reported is not nil but desired is
		reported := "not nil"
		var desired interface{} = nil

		// This should not panic
		result, err := merge.DesiredAny(context.Background(), reported, desired, merge.DesiredOptions{})
		assert.Error(t, err, "Expected an error when reported is not nil and desired is nil")
		assert.Nil(t, result, "Expected nil result")
	})

	t.Run("EdgeCase_ReflectInvalidValues", func(t *testing.T) {
		// Using a nil interface which will produce an invalid reflect.Value
		var reported interface{}
		var desired interface{}

		// This should not panic
		result, err := merge.DesiredAny(context.Background(), reported, desired, merge.DesiredOptions{})
		assert.NoError(t, err, "Nil interfaces should not cause an error")
		assert.Nil(t, result, "Expected nil result")
	})
}
