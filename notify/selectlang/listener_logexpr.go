package selectlang

import (
	"fmt"
	"strconv"
)

type cop struct {
	c  *Constraint
	op ConstraintLogicalOp
}
type LogExprListener struct {
	debug bool

	expr              LoggerExpression
	currentOperation  ConstraintLogicalOp
	currentConstraint *Constraint
	stack             Stack[cop]
}

func (s *LogExprListener) Enter(ctx *LoggerExprContext) *LogExprListener {
	if s.debug {
		fmt.Println("LOGGER ENTER")
	}

	// Reset the expression
	s.expr = LoggerExpression{}
	s.currentConstraint = nil
	s.stack = Stack[cop]{}

	if operations := ctx.LoggerOp(); operations != nil {
		if s.debug {
			fmt.Printf("    LOGGER_OP:%#v\n", ToStringList(operations, ","))
		}

		s.expr.CaptureOperations = ToStringList(operations, ",")
	}

	if re := ctx.Regex(); re != nil {
		if s.debug {
			fmt.Println("    REGEX:", re.REGEX().GetText())
		}

		s.expr.CaptureRegex = re.REGEX().GetText()
	}

	if mv := ctx.MapVarExpr(); mv != nil {
		// child.EQ()
		if s.debug {
			fmt.Println("    MAP_VAR_EXPR:", mv.STRING().GetText())
		}

		s.expr.CaptureEqMapVarExpr = mv.STRING().GetText()
	}

	return s
}

func (s *LogExprListener) Exit(ctx *LoggerExprContext) *LoggerExpression {
	if s.debug {
		fmt.Println("LOGGER EXIT")
	}

	if len(s.expr.CaptureOperations) == 0 {
		return nil
	}

	s.expr.Where = s.currentConstraint

	return &s.expr
}

func (s *LogExprListener) EnterConstraints(ctx *LoggerConstraintsContext) {
	if s.debug {
		fmt.Println("   ENTER_LOGGER_CONSTRAINTS")
	}

	s.currentOperation = ConstraintLogicalLHS
}

func (s *LogExprListener) ExitConstraints(ctx *LoggerConstraintsContext) {
	if s.debug {
		fmt.Println("   EXIT_LOGGER_CONSTRAINTS")
	}
	// NOOP
}

func (s *LogExprListener) EnterValueCondition(ctx *ValueFactorContext) {
	if compareOp := ctx.CompareOp(); compareOp != nil {
		var constraint *Constraint

		if s.currentOperation == ConstraintLogicalLHS {
			constraint = s.currentConstraint
		} else {
			constraint = &Constraint{}
		}

		constraint.Variable = "value"
		constraint.CompareOp = compareOp.GetText()

		if cr := ctx.ConstantOrRegex(); cr != nil {
			switch t := cr.(type) {
			case *NumericLiteralContext:
				if n, err := strconv.ParseFloat(t.NUMBER().GetText(), 64); err == nil {
					constraint.Value = n
					constraint.ValueType = ConstrainValueNumber
				}
			case *StringLiteralContext:
				constraint.Value = t.STRING().GetText()
				constraint.ValueType = ConstrainValueString
			case *RegexLiteralContext:
				re := t.Regex().REGEX().GetText()
				// strip the pre and postfix '/' to get the regex
				constraint.Value = re[1 : len(re)-1]
				constraint.ValueType = ConstrainValueRegex
			}
		}

		if s.debug {
			fmt.Println("    VALUE_CONSTRAINT:", constraint)
		}

		switch s.currentOperation {
		case ConstraintLogicalOpAnd:
			s.currentConstraint.And = append(s.currentConstraint.And, constraint)
		case ConstraintLogicalOpOr:
			s.currentConstraint.Or = append(s.currentConstraint.Or, constraint)
		}
	}
}

func (s *LogExprListener) And() {
	s.currentOperation = ConstraintLogicalOpAnd
}

func (s *LogExprListener) Or() {
	s.currentOperation = ConstraintLogicalOpOr
}

func (s *LogExprListener) ScopeStart() {
	c := &Constraint{}

	switch s.currentOperation {
	case ConstraintLogicalOpAnd:
		s.currentConstraint.And = append(s.currentConstraint.And, c)
	case ConstraintLogicalOpOr:
		s.currentConstraint.Or = append(s.currentConstraint.Or, c)
	}

	if s.currentConstraint == nil {
		s.stack.Push(cop{c: c, op: s.currentOperation})
	} else {
		s.stack.Push(cop{c: s.currentConstraint, op: s.currentOperation})
	}

	s.currentConstraint = c

	s.currentOperation = ConstraintLogicalLHS // start with left hand side
}

func (s *LogExprListener) ScopeEnd() {
	// NOOP
	c := s.stack.Pop()

	s.currentConstraint = c.c
	s.currentOperation = c.op
}
