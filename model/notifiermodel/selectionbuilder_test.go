package notifiermodel_test

import (
	"fmt"
	"testing"

	"github.com/mariotoffia/godeviceshadow/model/notifiermodel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type S struct {
	R bool
	V []notifiermodel.SelectedValue
	N string
}

func (s *S) Select(operation notifiermodel.NotifierOperation, value bool) (bool, []notifiermodel.SelectedValue) {
	if dv, ok := operation.Custom["debug"]; ok {
		if d, ok := dv.(bool); ok && d {
			fmt.Printf("%s -> %t\n", s.N, s.R)
		}
	}
	return s.R, s.V
}

func TestSingleAndSuccess(t *testing.T) {
	a := S{R: true}
	b := S{R: true}

	model, _ := notifiermodel.NewSelectionBuilder(&a).
		And(&b).
		Build()

	res, _ := model.Select(notifiermodel.NotifierOperation{}, false)

	assert.True(t, res, "Expected AND of two true Selections to be true")
}

func TestSingleAndFail(t *testing.T) {
	a := S{R: true}
	b := S{R: false}

	model, _ := notifiermodel.NewSelectionBuilder(&a).
		And(&b).
		Build()

	res, _ := model.Select(notifiermodel.NotifierOperation{}, false)

	assert.False(t, res, "Expected AND of (true, false) to be false")
}

func TestSingleOrSuccess(t *testing.T) {
	a := S{R: true}
	b := S{R: false}

	model, _ := notifiermodel.NewSelectionBuilder(&a).
		Or(&b).
		Build()

	res, _ := model.Select(notifiermodel.NotifierOperation{}, false)

	assert.True(t, res, "Expected OR of (true, false) to be true")
}

func TestSingleOrFail(t *testing.T) {
	a := S{R: false}
	b := S{R: false}

	model, _ := notifiermodel.NewSelectionBuilder(&a).
		Or(&b).
		Build()

	res, _ := model.Select(notifiermodel.NotifierOperation{}, false)

	assert.False(t, res, "Expected OR of (false, false) to be false")
}

func TestSingleNotSuccess(t *testing.T) {
	// Create a Selection that returns false
	a := S{R: false}

	// Build the final selection by negating 'a'
	model, _ := notifiermodel.NewSelectionBuilder(&a).
		Not().
		Build()

	// Invoke Select and expect the negated result to be true
	res, _ := model.Select(notifiermodel.NotifierOperation{}, false)

	assert.True(t, res, "Expected NOT of false Selection to be true")
}

func TestSingleNotFail(t *testing.T) {
	// Create a Selection that returns true
	a := S{R: true}

	// Build the final selection by negating 'a'
	model, _ := notifiermodel.NewSelectionBuilder(&a).
		Not().
		Build()

	// Invoke Select and expect the negated result to be false
	res, _ := model.Select(notifiermodel.NotifierOperation{}, false)

	assert.False(t, res, "Expected NOT of true Selection to be false")
}

func TestNestedAndSuccess(t *testing.T) {
	// Create mock Selections that all return true
	a := S{R: true}
	b := S{R: true}
	c := S{R: true}

	// Build the final selection with nested AND using AndB
	model, _ := notifiermodel.NewSelectionBuilder(&a).
		And(notifiermodel.Scoped(&b, func(sb *notifiermodel.SelectionBuilder) {
			sb.And(&c)
		})).
		Build()

	// Invoke Select and expect the combined result to be true
	res, _ := model.Select(notifiermodel.NotifierOperation{}, false)

	assert.True(t, res, "Expected nested AND (a AND (b AND c)) to be true")
}

func TestNestedAndFail(t *testing.T) {
	// Create mock Selections
	a := S{R: true}  // Outer Selection returns true
	b := S{R: true}  // First nested Selection returns true
	c := S{R: false} // Second nested Selection returns false

	// Build the final selection with nested AND using AndB
	model, _ := notifiermodel.NewSelectionBuilder(&a).
		And(notifiermodel.Scoped(&b, func(sb *notifiermodel.SelectionBuilder) {
			sb.And(&c) // (b AND c) where c is false
		})).
		Build()

	// Invoke Select and expect the combined result to be false
	res, _ := model.Select(notifiermodel.NotifierOperation{}, false)

	assert.False(t, res, "Expected nested AND (a AND (b AND c)) to be false because c is false")
}

func TestNestedOrSuccess(t *testing.T) {
	// Create mock Selections
	a := S{R: false} // Outer Selection returns false
	b := S{R: true}  // First nested Selection returns true
	c := S{R: true}  // Second nested Selection returns true

	// Build the final selection with nested OR using OrB
	model, _ := notifiermodel.NewSelectionBuilder(&a).
		Or(notifiermodel.Scoped(&b, func(sb *notifiermodel.SelectionBuilder) {
			sb.Or(&c) // (b OR c) where both are true
		})).
		Build()

	// Invoke Select and expect the combined result to be true
	res, _ := model.Select(notifiermodel.NotifierOperation{}, false)

	assert.True(t, res, "Expected nested OR (a OR (b OR c)) to be true because at least one nested Selection is true")
}

func TestNestedOrFail(t *testing.T) {
	// Create mock Selections
	a := S{R: false} // Outer Selection returns false
	b := S{R: false} // First nested Selection returns false
	c := S{R: false} // Second nested Selection returns false

	// Build the final selection with nested OR using OrB
	model, _ := notifiermodel.NewSelectionBuilder(&a).
		Or(notifiermodel.Scoped(&b, func(sb *notifiermodel.SelectionBuilder) {
			sb.Or(&c) // (b OR c) where both are false
		})).
		Build()

	// Invoke Select and expect the combined result to be false
	res, _ := model.Select(notifiermodel.NotifierOperation{}, false)

	assert.False(t, res, "Expected nested OR (a OR (b OR c)) to be false because all Selections are false")
}

func TestNestedAndOrCombination(t *testing.T) {
	// Create mock Selections
	a := S{R: true}  // a returns true
	b := S{R: false} // b returns false
	c := S{R: true}  // c returns true
	d := S{R: true}  // d returns true
	e := S{R: true}  // e returns true

	// Build the final selection with a combination of nested AND and OR
	// Expression: (a AND (b OR c)) OR (d AND e)
	model, _ := notifiermodel.NewSelectionBuilder(&a).
		And(notifiermodel.Scoped(&b, func(sb *notifiermodel.SelectionBuilder) {
			sb.Or(&c) // (b OR c)
		})).
		Or(notifiermodel.Scoped(&d, func(sb *notifiermodel.SelectionBuilder) {
			sb.And(&e) // (d AND e)
		})).
		Build()

	// Invoke Select and expect the combined result to be true
	res, _ := model.Select(notifiermodel.NotifierOperation{}, false)

	assert.True(t, res, "Expected ((a AND (b OR c)) OR (d AND e)) to be true")
}

func TestNestedNot(t *testing.T) {
	// Create mock Selections
	a := S{R: true}  // a returns true
	b := S{R: true}  // b returns true
	c := S{R: false} // c returns false

	// Build the sub-expression (b OR c)
	subOr, _ := notifiermodel.NewSelectionBuilder(&b).
		Or(&c).
		Build()

	// Build the main expression (a AND (b OR c))
	mainAnd, _ := notifiermodel.NewSelectionBuilder(&a).
		And(subOr).
		Build()

	// Apply NOT to the main expression: NOT (a AND (b OR c))
	finalSelection, _ := notifiermodel.NewSelectionBuilder(mainAnd).
		Not().
		Build()

	// Invoke Select and expect the negated result to be false because (a AND (b OR c)) is true
	res, _ := finalSelection.Select(notifiermodel.NotifierOperation{}, false)

	assert.False(t, res, "Expected NOT (a AND (b OR c)) to be false when (a AND (b OR c)) is true")
}

func TestMultipleNestedLevels(t *testing.T) {
	// Create mock Selections with varying return values
	a := S{N: "a", R: true}
	b := S{N: "b", R: false}
	c := S{N: "c", R: true}
	d := S{N: "d", R: true}
	e := S{N: "e", R: false}
	f := S{N: "f", R: true}
	g := S{N: "g", R: false}
	h := S{N: "h", R: true}

	/*
		Expected Expression: ((a AND (b OR (c AND d))) OR e) AND (f OR (g AND h))

		Step-by-step construction:
		1. Scoped(b, OR(c AND d))  --> (b OR (c AND d))
		2. a AND (b OR (c AND d))  --> (a AND (b OR (c AND d)))
		3. (a AND (b OR (c AND d))) OR e --> ((a AND (b OR (c AND d))) OR e)
		4. Scoped(g, AND(h)) --> (g AND h)
		5. f OR (g AND h) --> (f OR (g AND h))
		6. ((a AND (b OR (c AND d))) OR e) AND (f OR (g AND h))
	*/

	// Build the selection
	sb := notifiermodel.NewSelectionBuilder(&a).
		And(notifiermodel.Scoped(&b, func(sb *notifiermodel.SelectionBuilder) {
			sb.Or(notifiermodel.Scoped(&c, func(sb *notifiermodel.SelectionBuilder) {
				sb.And(&d)
			}))
		})).
		Or(&e).
		And(notifiermodel.Scoped(&f, func(sb *notifiermodel.SelectionBuilder) {
			sb.Or(notifiermodel.Scoped(&g, func(sb *notifiermodel.SelectionBuilder) {
				sb.And(&h)
			}))
		}))

	// Build the final selection
	finalSelection, err := sb.Build()
	require.NoError(t, err, "Expected selection to build without errors")

	// Invoke Select and expect the final result to be true
	res, _ := finalSelection.Select(notifiermodel.NotifierOperation{}, false)

	assert.True(t, res, "Expected the complex nested expression to evaluate to true")
}

func TestComplexExpressionAllTrue(t *testing.T) {
	// Create mock Selections that all return true
	a := S{N: "a", R: true}
	b := S{N: "b", R: true}
	c := S{N: "c", R: true}
	d := S{N: "d", R: true}
	e := S{N: "e", R: true}
	f := S{N: "f", R: true}
	g := S{N: "g", R: true}
	h := S{N: "h", R: true}

	/*
		Expected Expression: ((a OR (b AND (c OR d))) AND e) OR (f AND (g OR h))

		Step-by-step construction:
		1. Scoped(c, OR(d))  --> (c OR d)
		2. Scoped(b, AND(c OR d))  --> (b AND (c OR d))
		3. a OR (b AND (c OR d))  --> (a OR (b AND (c OR d)))
		4. (a OR (b AND (c OR d))) AND e  --> ((a OR (b AND (c OR d))) AND e)
		5. Scoped(g, OR(h))  --> (g OR h)
		6. f AND (g OR h) --> (f AND (g OR h))
		7. ((a OR (b AND (c OR d))) AND e) OR (f AND (g OR h))
	*/

	// Build the selection
	sb := notifiermodel.NewSelectionBuilder(
		notifiermodel.Scoped(&a, func(sb *notifiermodel.SelectionBuilder) {
			sb.Or(notifiermodel.Scoped(&b, func(sb *notifiermodel.SelectionBuilder) {
				sb.And(notifiermodel.Scoped(&c, func(sb *notifiermodel.SelectionBuilder) {
					sb.Or(&d)
				}))
			}))
		})).
		And(&e).
		Or(notifiermodel.Scoped(&f, func(sb *notifiermodel.SelectionBuilder) {
			sb.And(notifiermodel.Scoped(&g, func(sb *notifiermodel.SelectionBuilder) {
				sb.Or(&h)
			}))
		}))

	// Build the final selection
	finalSelection, err := sb.Build()
	require.NoError(t, err, "Expected selection to build without errors")

	// Invoke Select and expect the final result to be true
	res, _ := finalSelection.Select(notifiermodel.NotifierOperation{}, false)

	assert.True(t, res, "Expected the complex nested expression to evaluate to true")
}

func TestComplexExpressionMixedValues(t *testing.T) {
	// Create mock Selections with mixed true/false values
	a := S{N: "a", R: false}
	b := S{N: "b", R: true}
	c := S{N: "c", R: false}
	d := S{N: "d", R: true}
	e := S{N: "e", R: true}
	f := S{N: "f", R: false}
	g := S{N: "g", R: true}
	h := S{N: "h", R: false}

	/*
		Expected Expression: ((a OR (b AND (c OR d))) AND e) OR (f AND (g OR h))

		Step-by-step construction:
		1. Scoped(c, OR(d))  --> (c OR d)  → (false OR true) → **true**
		2. Scoped(b, AND(c OR d))  --> (b AND (c OR d)) → (true AND true) → **true**
		3. a OR (b AND (c OR d))  --> (a OR true) → (false OR true) → **true**
		4. (a OR (b AND (c OR d))) AND e  --> (true AND e) → (true AND true) → **true**
		5. Scoped(g, OR(h))  --> (g OR h) → (true OR false) → **true**
		6. f AND (g OR h) --> (f AND true) → (false AND true) → **false**
		7. ((a OR (b AND (c OR d))) AND e) OR (f AND (g OR h)) → (true OR false) → **true**
	*/

	// Build the selection
	sb := notifiermodel.NewSelectionBuilder(
		notifiermodel.Scoped(&a, func(sb *notifiermodel.SelectionBuilder) {
			sb.Or(notifiermodel.Scoped(&b, func(sb *notifiermodel.SelectionBuilder) {
				sb.And(notifiermodel.Scoped(&c, func(sb *notifiermodel.SelectionBuilder) {
					sb.Or(&d)
				}))
			}))
		})).
		And(&e).
		Or(notifiermodel.Scoped(&f, func(sb *notifiermodel.SelectionBuilder) {
			sb.And(notifiermodel.Scoped(&g, func(sb *notifiermodel.SelectionBuilder) {
				sb.Or(&h)
			}))
		}))

	// Build the final selection
	finalSelection, err := sb.Build()
	require.NoError(t, err, "Expected selection to build without errors")

	// Invoke Select and expect the final result to be true
	res, _ := finalSelection.Select(notifiermodel.NotifierOperation{}, false)

	assert.True(t, res, "Expected the complex nested expression with mixed values to evaluate to true")
}

func TestTruthTableExpression(t *testing.T) {
	// tc represents a single row in the truth table
	type tc struct {
		A, B, C, D, E, F, G, H, I, J, K bool
		Expected                        bool
	}

	// evaluateExpression manually computes the expected result for given inputs
	evaluateExpression := func(a, b, c, d, e, f, g, h, i, j, k bool) bool { /* NOSONAR */
		leftInner := d || e          // (d OR e)
		leftMiddle := c && leftInner // (c AND (d OR e))
		leftOuter := b || leftMiddle // (b OR (c AND (d OR e)))
		leftSide := a && leftOuter   // (a AND (b OR (c AND (d OR e))))
		leftFinal := leftSide || f   // ((a AND (b OR (c AND (d OR e)))) OR f)

		rightInner := i || j           // (i OR j)
		rightMiddle := h && rightInner // (h AND (i OR j))
		rightSide := g || rightMiddle  // (g OR (h AND (i OR j)))

		result := (leftFinal && rightSide) || k // (((leftFinal) AND (rightSide)) OR k)
		return result
	}

	// Generate all possible combinations of 11 boolean variables
	var testCases []tc
	for a := 0; a < 2; a++ {
		for b := 0; b < 2; b++ {
			for c := 0; c < 2; c++ {
				for d := 0; d < 2; d++ {
					for e := 0; e < 2; e++ {
						for f := 0; f < 2; f++ {
							for g := 0; g < 2; g++ {
								for h := 0; h < 2; h++ {
									for i := 0; i < 2; i++ {
										for j := 0; j < 2; j++ {
											for k := 0; k < 2; k++ {
												tc := tc{
													A:        a == 1,
													B:        b == 1,
													C:        c == 1,
													D:        d == 1,
													E:        e == 1,
													F:        f == 1,
													G:        g == 1,
													H:        h == 1,
													I:        i == 1,
													J:        j == 1,
													K:        k == 1,
													Expected: evaluateExpression(a == 1, b == 1, c == 1, d == 1, e == 1, f == 1, g == 1, h == 1, i == 1, j == 1, k == 1),
												}
												testCases = append(testCases, tc)
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	// Run tests for all truth table cases
	for _, tc := range testCases {
		// Create mock selections
		a := S{N: "a", R: tc.A}
		b := S{N: "b", R: tc.B}
		c := S{N: "c", R: tc.C}
		d := S{N: "d", R: tc.D}
		e := S{N: "e", R: tc.E}
		f := S{N: "f", R: tc.F}
		g := S{N: "g", R: tc.G}
		h := S{N: "h", R: tc.H}
		i := S{N: "i", R: tc.I}
		j := S{N: "j", R: tc.J}
		k := S{N: "k", R: tc.K}

		// Build selection using Scoped to match expression (((a AND (b OR (c AND (d OR e)))) OR f) AND (g OR (h AND (i OR j)))) OR k)
		sb := notifiermodel.NewSelectionBuilder(
			notifiermodel.Scoped(&a, func(sb *notifiermodel.SelectionBuilder) {
				sb.And(notifiermodel.Scoped(&b, func(sb *notifiermodel.SelectionBuilder) {
					sb.Or(notifiermodel.Scoped(&c, func(sb *notifiermodel.SelectionBuilder) {
						sb.And(notifiermodel.Scoped(&d, func(sb *notifiermodel.SelectionBuilder) {
							sb.Or(&e)
						}))
					}))
				}))
			})).
			Or(&f).
			And(notifiermodel.Scoped(&g, func(sb *notifiermodel.SelectionBuilder) {
				sb.Or(notifiermodel.Scoped(&h, func(sb *notifiermodel.SelectionBuilder) {
					sb.And(notifiermodel.Scoped(&i, func(sb *notifiermodel.SelectionBuilder) {
						sb.Or(&j)
					}))
				}))
			})).
			Or(&k)

		// Build final selection
		finalSelection, err := sb.Build()
		require.NoError(t, err, "Expected selection to build without errors")

		// Invoke Select and check against expected result
		res, _ := finalSelection.Select(notifiermodel.NotifierOperation{}, false)

		assert.Equal(t, tc.Expected, res, "Mismatch for combination: %+v", tc)
	}
}
