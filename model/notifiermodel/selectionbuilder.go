package notifiermodel

import "fmt"

// SelectionBuilder helps build a Selection with fluent method calls.
type SelectionBuilder struct {
	current Selection
	err     error
}

// NewSelectionBuilder initializes a new builder with a *single* initial Selection.
func NewSelectionBuilder(initial Selection) *SelectionBuilder {
	return &SelectionBuilder{current: initial}
}

func (b *SelectionBuilder) String() string {
	if b.current == nil {
		return "<empty selection>"
	}

	return fmt.Sprintf("Selection(%v)", b.current)
}

// Build returns the final composed Selection.
func (b *SelectionBuilder) Build() (Selection, error) {
	if b.err != nil {
		return nil, b.err
	}

	return b.current, b.err
}

// And will and one or more _selections_ with the current selection. This
// is equivalent to x AND y AND z.
func (b *SelectionBuilder) And(selections ...Selection) *SelectionBuilder {
	if b.err != nil {
		return b
	}

	if len(selections) == 0 {
		b.err = fmt.Errorf("cannot AND selections when no selections are provided")
		return b
	}

	if b.current == nil {
		b.err = fmt.Errorf("cannot AND selections when current is nil")
		return b
	}

	// Flatten all AND selections in one call
	allSelections := append([]Selection{b.current}, selections...)
	b.current = And(allSelections...)

	return b
}

// Or will or one or more _selections_ with the current selection. This
// is equivalent to x OR y OR z.
func (b *SelectionBuilder) Or(selections ...Selection) *SelectionBuilder {
	if b.err != nil {
		return b
	}

	if len(selections) == 0 {
		b.err = fmt.Errorf("cannot OR selections when no selections are provided")
		return b
	}

	if b.current == nil {
		b.err = fmt.Errorf("cannot OR selections when current is nil")
		return b
	}

	// Flatten all OR selections in one call
	allSelections := append([]Selection{b.current}, selections...)
	b.current = Or(allSelections...)

	return b
}

// Not wraps the current builder selection in a logical NOT.
func (b *SelectionBuilder) Not() *SelectionBuilder {
	if b.current == nil {
		b.err = fmt.Errorf("cannot NOT selections when current is nil")
		return b
	}

	b.current = Not(b.current)
	return b
}

// Scoped will push the _selection_ and the other selections in f inside a "parenthesis". This can be
// used in combination with `And` and `Or` to create more complex selections. For example:
// A AND (B OR C) can be created by:
//
// ```
// sb := NewSelectionBuilder(A)
//
//	sb.And(Scoped(B, func(sb *SelectionBuilder) {
//	  sb.Or(C)
//	}))
//
// ```
//
// The _selection_ is the one to start the builder with and the `f` is a function that will be called
// with a new `SelectionBuilder` instance. The `f` function can then call `And` and `Or` on the new
// `SelectionBuilder` instance to create the scoped selection.
//
// When it fails, it will return nil
func Scoped(selection Selection, f func(sb *SelectionBuilder)) Selection {
	sb := NewSelectionBuilder(selection)
	f(sb)

	if res, err := sb.Build(); err != nil {
		return nil
	} else {
		return res
	}
}
