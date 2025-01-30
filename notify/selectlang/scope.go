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
	And       []Scope
	Or        []Scope
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

func (c Constraint) String() string {
	return fmt.Sprintf("%s %s %v (%s)", c.Variable, c.CompareOp, c.Value, c.ValueType)
}

type ConstraintLogicalOp int

const (
	// ConstraintLogicalLHS is no logical operation, instead it is the left hand side
	// of a logical operation.
	ConstraintLogicalLHS ConstraintLogicalOp = iota
	ConstraintLogicalOpAnd
	ConstraintLogicalOpOr
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

func (scope Scope) Children() []Scope {
	res := append(scope.And, scope.Or...)

	if scope.Not != nil {
		return append(res, *scope.Not)
	}

	return res
}

func (scope Scope) PrintScopeTree(writer io.Writer, indent int) {
	prefix := strings.Repeat("  ", indent)

	fmt.Fprintln(writer, prefix, "{")
	defer fmt.Fprintln(writer, prefix, "}")

	if scope.Primary != nil {
		fmt.Fprintln(writer, prefix, "  Primary:", scope.Primary)
	} else if scope.Logger != nil {
		fmt.Fprintln(writer, prefix, "  Logger:", scope.Logger)
	}

	if len(scope.And) > 0 {
		fmt.Fprint(writer, prefix, "AND ")

		for _, and := range scope.And {
			and.PrintScopeTree(writer, indent+1)
		}
	}

	if len(scope.Or) > 0 {
		fmt.Fprint(writer, prefix, "OR ")

		for _, or := range scope.Or {
			or.PrintScopeTree(writer, indent+1)
		}
	}

	if scope.Not != nil {
		fmt.Fprint(writer, prefix, "NOT ")

		scope.Not.PrintScopeTree(writer, indent+1)
	}
}
