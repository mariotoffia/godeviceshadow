package selectlang

import (
	"fmt"
	"io"
	"strings"
)

// Scope encapsulates a '(' and ')' pair with type info and
// collects the expression within the scope.
type Scope struct {
	// ScopeType is the type of scope. When the '(' is hit it will be
	// set to `ScopeLoggerUntyped` and then, depending on the content
	// it will be set to correct type.
	ScopeType ScopeType
	Primary   *PrimaryExpression
	Logger    *LoggerExpression
	And       []*Scope
	Or        []*Scope
	Not       *Scope
}

type PrimaryExpression struct {
	ID        string
	Name      string
	Operation []string
}

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

func (le LoggerExpression) String() string {
	return fmt.Sprintf("op: %v, re: %s, eq: %s",
		le.CaptureOperations, le.CaptureRegex, le.CaptureEqMapVarExpr)
}

func (le LoggerExpression) Dump(writer io.Writer, indent int) {
	prefix := strings.Repeat("  ", indent)

	fmt.Fprintln(writer, prefix, le)

	if le.Where != nil {
		fmt.Fprintln(writer, prefix, "WHERE {")

		if le.Where.HasConstrainValues() {
			fmt.Fprint(writer, strings.Repeat("  ", indent+1), le.Where)
		}

		if le.Where.IsScoped() {
			le.Where.Dump(writer, indent+1)
		} else {
			fmt.Fprintln(writer, prefix)
		}
		fmt.Fprintln(writer, prefix, "}")
	}
}

func (c Constraint) String() string {
	return fmt.Sprintf("%s %s %v (%s)", c.Variable, c.CompareOp, c.Value, c.ValueType)
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

func (c Constraint) Dump(writer io.Writer, indent int) {
	prefix := strings.Repeat("  ", indent)

	for _, and := range c.And {
		if !and.IsScoped() {
			fmt.Fprint(writer, " AND ", and)
		}
	}

	for _, or := range c.Or {
		if !or.IsScoped() {
			fmt.Fprint(writer, " OR ", or)
		}
	}

	for _, and := range c.And {
		if and.IsScoped() {
			fmt.Fprintln(writer, " AND {")
			fmt.Fprint(writer, strings.Repeat("  ", indent+1), and)

			and.Dump(writer, indent+1)

			fmt.Fprintln(writer)
			fmt.Fprintln(writer, prefix, "}")
		}
	}

	for _, or := range c.Or {
		if or.IsScoped() {
			fmt.Fprintln(writer, " OR {")
			fmt.Fprint(writer, strings.Repeat("  ", indent+1), or)

			or.Dump(writer, indent+1)

			fmt.Fprintln(writer)
			fmt.Fprintln(writer, prefix, "}")
		}
	}
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

type ScopeType int

const (
	ScopeLoggerUntyped ScopeType = iota
	ScopeLoggerExpr
	ScopeTypePrimaryExpr
)

func (scope Scope) Children() []*Scope {
	res := append(scope.And, scope.Or...)

	if scope.Not != nil {
		return append(res, scope.Not)
	}

	return res
}

func (scope Scope) Dump(writer io.Writer, indent int) {
	prefix := strings.Repeat("  ", indent)

	fmt.Fprintln(writer, prefix, "{")
	defer fmt.Fprintln(writer, prefix, "}")

	if scope.Primary != nil {
		fmt.Fprintln(writer, prefix, "  Primary:", scope.Primary)
	} else if scope.Logger != nil {
		fmt.Fprint(writer, prefix, "  Logger:")
		scope.Logger.Dump(writer, indent+1)
	}

	if len(scope.And) > 0 {
		fmt.Fprint(writer, prefix, "AND ")

		for _, and := range scope.And {
			and.Dump(writer, indent+1)
		}
	}

	if len(scope.Or) > 0 {
		fmt.Fprint(writer, prefix, "OR ")

		for _, or := range scope.Or {
			or.Dump(writer, indent+1)
		}
	}

	if scope.Not != nil {
		fmt.Fprint(writer, prefix, "NOT ")

		scope.Not.Dump(writer, indent+1)
	}
}
