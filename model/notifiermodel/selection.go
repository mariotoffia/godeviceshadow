package notifiermodel

import (
	"time"

	"github.com/mariotoffia/godeviceshadow/model"
)

type SelectedValue struct {
	Path         string
	OldValue     model.ValueAndTimestamp
	NewValue     model.ValueAndTimestamp
	OldTimeStamp time.Time
	NewTimeStamp time.Time
}

// Selection is a selection that will receive a `NotifierOperation` and if it is
// selected it will be returned. If the caller is interested in single values that where
// selected the `SelectedValue`(s) are returned.
//
// If no `SelectedValue` is wanted, the selection has higher performance, since it will return
// true on first match.
//
// If the id, name was matched and _value_ is set to `true` it will return the all values as selected.
//
// TIP: Use the `SelectionBuilder` or functions `And`, `Or` and `Not` to build a selection to create a complete
// expression of `Selection`.
type Selection interface {
	// Select will select the operation and return true if it is selected. If _value_ is set to `true`
	// it will return the selected values.
	//
	// If the selection is not matching in the operation it will return false.
	Select(operation NotifierOperation, value bool) (selected bool, values []SelectedValue)
}

// AndSelection implements the Selection interface
type AndSelection struct {
	Selections []Selection
}

// And returns a Selection that represents the logical AND of all the passed in Selections
func And(selections ...Selection) Selection {
	return &AndSelection{Selections: selections}
}

// Select executes each child Selection; if one fails (returns selected=false), we short-circuit.
func (a *AndSelection) Select(op NotifierOperation, value bool) (bool, []SelectedValue) {
	var allValues []SelectedValue

	for _, s := range a.Selections {
		selected, vals := s.Select(op, value)

		if !selected {
			return false, nil
		}

		// accumulate
		allValues = append(allValues, vals...)
	}
	return true, allValues
}

// OrSelection implements the Selection interface.
type OrSelection struct {
	Selections []Selection
}

// Or returns a Selection that represents the logical OR of all the passed-in Selections.
func Or(selections ...Selection) Selection {
	return &OrSelection{Selections: selections}
}

// Select executes each child Selection. If 'value' is true, it accumulates values
// from *all* matching Selections. If 'value' is false, it can short-circuit on the first match.
func (o *OrSelection) Select(op NotifierOperation, value bool) (bool, []SelectedValue) {
	var anyMatched bool
	var allValues []SelectedValue

	for _, s := range o.Selections {
		selected, vals := s.Select(op, value)

		if selected {
			anyMatched = true

			if value {
				allValues = append(allValues, vals...)
			} else {
				return true, nil // Short-circuits correctly
			}
		}
	}

	if !anyMatched {
		return false, nil
	}

	return true, allValues
}

// NotSelection implements the Selection interface
type NotSelection struct {
	Negated Selection
}

// Not returns a Selection that represents the logical NOT of a child Selection
func Not(s Selection) Selection {
	return &NotSelection{Negated: s}
}

func (n *NotSelection) Select(op NotifierOperation, value bool) (bool, []SelectedValue) {
	selected, values := n.Negated.Select(op, value)

	// negate since not
	selected = !selected

	if selected {
		return true, values
	}

	return false, nil
}

type FuncSelection struct {
	F func(op NotifierOperation, value bool) (bool, []SelectedValue)
}

func (s *FuncSelection) Select(operation NotifierOperation, value bool) (bool, []SelectedValue) {
	return s.F(operation, value)
}
