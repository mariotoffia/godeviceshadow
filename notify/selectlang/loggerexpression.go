package selectlang

import "github.com/mariotoffia/godeviceshadow/model/notifiermodel"

type LoggerExpression struct {
	CaptureOperations   []string
	CaptureRegex        string
	CaptureEqMapVarExpr string
	Where               *Constraint
}

type Constraint struct {
	Variable  string
	CompareOp string
	Value     any // string, number, or regex
	ValueType ConstrainValueType
	And       []*Constraint
	Or        []*Constraint
}

// IsScoped will return `true` if the constraint has it's own scope, i.e. it is within '(' and ')'.
func (c Constraint) IsScoped() bool {
	return len(c.And) > 0 || len(c.Or) > 0
}

// IsOnlyScope is when the constraint itself do not have any variable, compare operation or value and
// but have `And` or `Or` constraints. Therefore it is simply a scope and nothing else.
func (c Constraint) IsOnlyScope() bool {
	return !c.HasConstrainValues() && c.IsScoped()
}

// HasConstrainValues returns `true` if the constraint has any value constraints. Independent of
// `And` or `Or` constraints.
func (c Constraint) HasConstrainValues() bool {
	return c.Variable != "" && c.CompareOp != "" && c.Value != nil
}

type ConstraintLogicalOp int

const (
	// ConstraintLogicalLHS is no logical operation, instead it is the left hand side
	// of a logical operation.
	ConstraintLogicalLHS ConstraintLogicalOp = iota
	ConstraintLogicalOpAnd
	ConstraintLogicalOpOr
	// ConstraintLogicalOpNot can only be used by `Scope` and not `Constraint`
	ConstraintLogicalOpNot
)

type ConstrainValueType int

const (
	// ConstrainValueString is a plain string
	ConstrainValueString ConstrainValueType = iota
	// ConstrainValueNumber is a float64 number
	ConstrainValueNumber
	// is a string that represents a regex
	ConstrainValueRegex
)

func (cvt ConstrainValueType) String() string {
	switch cvt {
	case ConstrainValueString:
		return "string"
	case ConstrainValueNumber:
		return "number"
	case ConstrainValueRegex:
		return "regex"
	}

	return "unknown"
}

func (le *LoggerExpression) ToMatchFunc() func(op notifiermodel.NotifierOperation, value bool) (bool, []notifiermodel.SelectedValue) {
	panic("TODO: not implemented")
}

func (c *Constraint) ToMatchFunc() func(op notifiermodel.NotifierOperation, value bool) (bool, []notifiermodel.SelectedValue) {
	panic("TODO: not implemented")
}
