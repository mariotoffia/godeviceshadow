package selectlang

import (
	"fmt"
	"io"
	"strings"
)

/*
(
ID: /myDevice-\d+/
AND
NAME: 'myShadow'
AND
    OPERATION:[]string{"report", "desired"}
)
AND
(
LOGGER START
    LOGGER_OP:[]string{"add", "update"}
    REGEX: /^Sensors-.*-indoor$/
    MAP_VAR_EXPR: 'temp'
   ENTER_LOGGER_CONSTRAINTS
(
    VARIABLE: var
    COMPARE_OP: GT
    NUMBER: 20
OR
(
    VARIABLE: var
    COMPARE_OP: LT
    REGEX: /re-\d+/
AND
    VARIABLE: var
    COMPARE_OP: NE
    STRING: 'apa'
)
)
   EXIT_LOGGER_CONSTRAINTS
LOGGER END
)
OR
(
LOGGER START
    LOGGER_OP:[]string{"add", "update"}
LOGGER END
)
*/

// Scope encapsulates a '(' and ')' pair with type info and
// collects the expression within the scope.
type Scope struct {
	// ScopeType is the type of scope. When the '(' is hit it will be
	// set to `ScopeLoggerUntyped` and then, depending on the content
	// it will be set to correct type.
	ScopeType ScopeType
	Primary   *PrimaryExpression
	Logger    []LoggerExpression
	Operator  OperatorType
	Children  []*Scope
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
	Constraints         Constraint
}

type Constraint struct {
	Variable  string
	CompareOp string
	Value     any          // string, number, or regex
	Operator  OperatorType // AND/OR operator for nested conditions
	Children  []Constraint // Nested constraints
}

type OperatorType int

const (
	OperatorNone OperatorType = iota
	OperatorAnd
	OperatorOr
)

type ScopeType int

const (
	ScopeLoggerUntyped ScopeType = iota
	ScopeLoggerExpr
	ScopeTypePrimaryExpr
)

func (scope Scope) PrintScopeTree(writer io.Writer, indent int) {
	prefix := strings.Repeat("  ", indent)

	if scope.Operator == OperatorAnd {
		fmt.Fprintln(writer, prefix, "AND {")
	} else if scope.Operator == OperatorOr {
		fmt.Fprintln(writer, prefix, "OR {")
	}

	if scope.Primary != nil {
		fmt.Fprintln(writer, prefix, "  Primary:", scope.Primary)
	}

	if len(scope.Logger) > 0 {
		fmt.Fprintln(writer, prefix, "  Logger:", scope.Logger)
	}

	for _, child := range scope.Children {
		child.PrintScopeTree(writer, indent+1)
	}

	if scope.Operator == OperatorAnd || scope.Operator == OperatorOr {
		fmt.Fprintln(writer, prefix, "}")
	}
}
