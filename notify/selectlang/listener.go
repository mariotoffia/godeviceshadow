package selectlang

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"
)

type ExpressionListener struct {
	*BaseselectlangListener

	debug bool

	scopes      Stack[Scope]
	constraints Stack[Constraint]
}

type ExpressionListenerOpts struct {
	Debug bool
}

func (s *ExpressionListener) RootScope() Scope {
	return s.scopes.Peek()
}

func NewExpressionListener(opts ...ExpressionListenerOpts) *ExpressionListener {
	var opt ExpressionListenerOpts

	if len(opts) > 0 {
		opt = opts[0]
	}

	return &ExpressionListener{
		debug: opt.Debug,
	}
}

func (s *ExpressionListener) VisitTerminal(node antlr.TerminalNode) {
	term := node.GetText()

	switch term {
	case "AND":
		if s.debug {
			fmt.Println("AND")
		}
		s.scopes.Update(func(scp Scope) Scope {
			scp.Operator = OperatorAnd
			return scp
		})

	case "OR":
		if s.debug {
			fmt.Println("OR")
		}
		s.scopes.Update(func(scp Scope) Scope {
			scp.Operator = OperatorOr
			return scp
		})

	case "(":
		if s.debug {
			fmt.Println("(")
		}
		s.scopes.Push(Scope{}) // Push new nested scope

	case ")":
		if s.debug {
			fmt.Println(")")
		}

		completedScope := s.scopes.Pop()

		// Attach the popped scope to the parent
		s.scopes.Update(func(scp Scope) Scope {
			scp.Children = append(scp.Children, &completedScope)
			return scp
		})
	}
}

func (s *ExpressionListener) EnterIdExpr(ctx *IdExprContext) {
	if child := FirstChild[*RegexOrStringContext](ctx); child != nil {
		if s.debug {
			fmt.Println("ID:", child.GetText())
		}

		s.scopes.Update(func(scp Scope) Scope {
			scp.ScopeType = ScopeTypePrimaryExpr

			if scp.Primary == nil {
				scp.Primary = &PrimaryExpression{ID: child.GetText()}
			} else {
				scp.Primary.ID = child.GetText()
			}

			return scp
		})
	}
}

func (s *ExpressionListener) EnterNameExpr(ctx *NameExprContext) {
	if child := FirstChild[*RegexOrStringContext](ctx); child != nil {
		if s.debug {
			fmt.Println("NAME:", child.GetText())
		}

		s.scopes.Update(func(scp Scope) Scope {
			scp.ScopeType = ScopeTypePrimaryExpr

			if scp.Primary == nil {
				scp.Primary = &PrimaryExpression{Name: child.GetText()}
			} else {
				scp.Primary.Name = child.GetText()
			}

			return scp
		})
	}
}

func (s *ExpressionListener) EnterOperationExpr(ctx *OperationExprContext) {
	operation := ctx.Operations()

	if operation == nil {
		return
	}

	s.scopes.Update(func(scp Scope) Scope {
		scp.ScopeType = ScopeTypePrimaryExpr

		if scp.Primary == nil {
			scp.Primary = &PrimaryExpression{Operation: ToStringList(operation, ",")}
		} else {
			scp.Primary.Operation = ToStringList(operation, ",")
		}

		return scp
	})

	if s.debug {
		fmt.Printf("    OPERATION:%#v\n", ToStringList(operation, ","))
	}
}

func (s *ExpressionListener) EnterLoggerExpr(ctx *LoggerExprContext) {

	if s.debug {
		fmt.Println("LOGGER START ")
	}

	s.scopes.Update(func(scp Scope) Scope {
		scp.ScopeType = ScopeLoggerExpr

		expr := LoggerExpression{}

		if child := FirstChild[*LoggerOpContext](ctx); child != nil {
			if s.debug {
				fmt.Printf("    LOGGER_OP:%#v\n", ToStringList(child, ","))
			}

			expr.CaptureOperations = ToStringList(child, ",")
		}

		if child := FirstChild[*RegexContext](ctx); child != nil {
			if s.debug {
				fmt.Println("    REGEX:", child.REGEX().GetText())
			}

			expr.CaptureRegex = child.REGEX().GetText()
		}

		if child := FirstChild[*MapVarExprContext](ctx); child != nil {
			// child.EQ()
			if s.debug {
				fmt.Println("    MAP_VAR_EXPR:", child.STRING().GetText())
			}

			expr.CaptureEqMapVarExpr = child.STRING().GetText()
		}

		scp.Logger = append(scp.Logger, expr)

		return scp
	})
}

func (s *ExpressionListener) EnterLoggerConstraints(ctx *LoggerConstraintsContext) {
	if s.debug {
		fmt.Println("   ENTER_LOGGER_CONSTRAINTS")
	}

	s.constraints = Stack[Constraint]{} // Reset constraints
}

func (s *ExpressionListener) ExitLoggerConstraints(ctx *LoggerConstraintsContext) {
	if s.debug {
		fmt.Println("   EXIT_LOGGER_CONSTRAINTS")
	}

	// Pop the completed WHERE condition tree
	whereCondition := s.constraints.Pop()

	// Attach it to the current LoggerExpression
	s.scopes.Update(func(scp Scope) Scope {
		if len(scp.Logger) > 0 {
			scp.Logger[len(scp.Logger)-1].Constraints = whereCondition
		}
		return scp
	})
}

func (s *ExpressionListener) EnterValueFactor(ctx *ValueFactorContext) {
	if compareOp := ctx.CompareOp(); compareOp != nil {
		if s.debug {
			fmt.Println("    VARIABLE: var")
		}

		switch {
		case compareOp.GT() != nil:
			if s.debug {
				fmt.Println("    COMPARE_OP: GT")
			}
		case compareOp.GE() != nil:
			if s.debug {
				fmt.Println("    COMPARE_OP: GE")
			}
		case compareOp.LT() != nil:
			if s.debug {
				fmt.Println("    COMPARE_OP: LT")
			}
		case compareOp.LE() != nil:
			if s.debug {
				fmt.Println("    COMPARE_OP: LE")
			}
		case compareOp.EQ() != nil:
			if s.debug {
				fmt.Println("    COMPARE_OP: EQ")
			}
		case compareOp.NE() != nil:
			if s.debug {
				fmt.Println("    COMPARE_OP: NE")
			}
		case compareOp.BEFORE() != nil:
			if s.debug {
				fmt.Println("    COMPARE_OP: BEFORE")
			}
		case compareOp.AFTER() != nil:
			if s.debug {
				fmt.Println("    COMPARE_OP: AFTER")
			}
		default:
			if s.debug {
				fmt.Println("    COMPARE_OP:", compareOp.GetText())
			}
		}
	}

	if cr := ctx.ConstantOrRegex(); cr != nil {
		switch t := cr.(type) {
		case *NumericLiteralContext:
			if s.debug {
				fmt.Println("    NUMBER:", t.NUMBER().GetText())
			}
		case *StringLiteralContext:
			if s.debug {
				fmt.Println("    STRING:", t.STRING().GetText())
			}
		case *RegexLiteralContext:
			if s.debug {
				fmt.Println("    REGEX:", t.Regex().REGEX().GetText())
			}
		}
	}
}

func (s *ExpressionListener) ExitLoggerExpr(ctx *LoggerExprContext) {
	if s.debug {
		fmt.Println("LOGGER END")
	}
}
