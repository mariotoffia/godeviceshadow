package selectlang

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"
)

type ExpressionListener struct {
	*BaseselectlangListener

	debug bool

	scopes         Stack[Scope]
	currentLogExpr *LogExprListener
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

		if s.currentLogExpr != nil {
			s.currentLogExpr.And()
		} else {
			s.scopes.Update(func(scp Scope) Scope {
				return scp
			})
		}
	case "OR":
		if s.debug {
			fmt.Println("OR")
		}

		if s.currentLogExpr != nil {
			s.currentLogExpr.Or()
		} else {
			s.scopes.Update(func(scp Scope) Scope {
				return scp
			})
		}
	case "NOT":
		if s.debug {
			fmt.Println("NOT")
		}

		if s.currentLogExpr != nil {
			// This is not supported in the log expression
		} else {
			s.scopes.Update(func(scp Scope) Scope {
				return scp
			})
		}
	case "(":
		if s.debug {
			fmt.Println("(")
		}
		if s.currentLogExpr != nil {
			s.currentLogExpr.ScopeStart()
		} else {
			// Push new nested scope
			s.scopes.Push(Scope{})
		}
	case ")":
		if s.debug {
			fmt.Println(")")
		}
		if s.currentLogExpr != nil {
			s.currentLogExpr.ScopeEnd()
		} else {
			s.scopes.Update(func(scp Scope) Scope {
				return scp
			})
		}
	}
}

func (s *ExpressionListener) EnterLoggerExpr(ctx *LoggerExprContext) {
	s.scopes.Update(func(scp Scope) Scope {
		scp.ScopeType = ScopeLoggerExpr
		return scp
	})

	s.currentLogExpr = (&LogExprListener{debug: s.debug}).Enter(ctx)
}

func (s *ExpressionListener) EnterLoggerConstraints(ctx *LoggerConstraintsContext) {
	s.currentLogExpr.EnterConstraints(ctx)
}

func (s *ExpressionListener) ExitLoggerConstraints(ctx *LoggerConstraintsContext) {
	s.currentLogExpr.ExitConstraints(ctx)
}

func (s *ExpressionListener) EnterValueFactor(ctx *ValueFactorContext) {
	s.currentLogExpr.EnterValueCondition(ctx)
}

func (s *ExpressionListener) ExitLoggerExpr(ctx *LoggerExprContext) {
	s.scopes.Update(func(scp Scope) Scope {
		scp.Logger = s.currentLogExpr.Exit(ctx)
		return scp
	})

	s.currentLogExpr = nil
}
