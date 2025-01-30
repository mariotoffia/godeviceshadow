package selectlang

import (
	"fmt"
	"io"
	"strings"
)

func (scope Scope) Dump(writer io.Writer, indent int) {
	prefix := strings.Repeat("  ", indent)

	fmt.Fprintln(writer, prefix, "{")
	defer fmt.Fprintln(writer, prefix, "}")

	doIndent := true

	if scope.Primary != nil {
		fmt.Fprint(writer, prefix, "  Primary", scope.Primary)
		doIndent = false
	} else if scope.Logger != nil {
		fmt.Fprint(writer, prefix, "  Logger")
		scope.Logger.Dump(writer, indent+1)

		doIndent = false
	}

	if doIndent {
		fmt.Fprint(writer, prefix)
	} else {
		fmt.Fprint(writer, " ")
	}

	if len(scope.And) > 0 {
		fmt.Fprint(writer, "AND ")

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

func (pe PrimaryExpression) String() string {
	s := ""

	if pe.ID != "" {
		s += fmt.Sprintf("ID: %s", pe.ID)
	}

	if pe.Name != "" {
		if s != "" {
			s += ", "
		}

		s += fmt.Sprintf("Name: %s", pe.Name)
	}

	if len(pe.Operation) > 0 {
		if s != "" {
			s += ", "
		}

		s += fmt.Sprintf("Operations: %v", pe.Operation)
	}

	return s
}

func (le LoggerExpression) String() string {
	s := ""

	if len(le.CaptureOperations) > 0 {
		s += fmt.Sprintf("Operations: %v", le.CaptureOperations)
	}

	if le.CaptureRegex != "" {
		if s != "" {
			s += ", "
		}

		s += fmt.Sprintf("RegEx: %s", le.CaptureRegex)
	}

	if le.CaptureEqMapVarExpr != "" {
		if s != "" {
			s += ", "
		}

		s += fmt.Sprintf("VarName: %s", le.CaptureEqMapVarExpr)
	}

	return s
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
