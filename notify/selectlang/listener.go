package selectlang

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"
)

type scop struct {
	scp *Scope
	op  ConstraintLogicalOp
}

type ExpressionListener struct {
	*BaseselectlangListener

	debug bool

	currentLogExpr *LogExprListener
	currentPrimary *PrimaryExprListener

	scopes           Stack[scop]
	currentLogicalOp ConstraintLogicalOp
}

type ExpressionListenerOpts struct {
	Debug bool
}

func (s *ExpressionListener) RootScope() Scope {
	if s.scopes.IsEmpty() {
		return Scope{}
	}

	if s.scopes.Size() == 1 {
		return *s.scopes.Peek().scp
	}

	child := s.scopes.Pop()

	var ptr *Scope

	for curr := s.scopes.Pop(); curr.scp != nil; curr = s.scopes.Pop() {
		ptr = curr.scp

		switch child.op {
		case ConstraintLogicalOpAnd:
			ptr.And = append(ptr.And, child.scp)
		case ConstraintLogicalOpOr:
			ptr.Or = append(ptr.Or, child.scp)
		case ConstraintLogicalOpNot:
			ptr.Not = child.scp
		}

		child = curr
	}

	s.scopes.Push(scop{scp: ptr, op: ConstraintLogicalLHS})
	return *ptr
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

func (s *ExpressionListener) EnterIdExpr(ctx *IdExprContext) {
	if s.currentPrimary == nil {
		s.currentLogExpr = nil
		s.currentPrimary = NewPrimaryExprListener(s.debug)
		s.currentPrimary.Enter()
	}

	s.currentPrimary.EnterIdExpr(ctx)
}

func (s *ExpressionListener) EnterNameExpr(ctx *NameExprContext) {
	if s.currentPrimary == nil {
		s.currentLogExpr = nil
		s.currentPrimary = NewPrimaryExprListener(s.debug)
		s.currentPrimary.Enter()
	}

	s.currentPrimary.EnterNameExpr(ctx)
}

func (s *ExpressionListener) EnterOperationExpr(ctx *OperationExprContext) {
	if s.currentPrimary == nil {
		s.currentLogExpr = nil
		s.currentPrimary = NewPrimaryExprListener(s.debug)
		s.currentPrimary.Enter()
	}

	s.currentPrimary.EnterOperationExpr(ctx)
}

func (s *ExpressionListener) VisitTerminal(node antlr.TerminalNode) {
	term := node.GetText()

	switch term {
	case "AND":
		s.and(node)
	case "OR":
		s.or(node)
	case "NOT":
		s.not(node)
	case "(":
		s.scopeStart(node)
	case ")":
		s.scopeEnd(node)
	}
}

func (s *ExpressionListener) EnterLoggerExpr(ctx *LoggerExprContext) {
	s.scopes.Update(func(scp scop) scop {
		scp.scp.ScopeType = ScopeLoggerExpr
		return scp
	})

	s.currentPrimary = nil
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
	s.scopes.Update(func(scp scop) scop {
		scp.scp.Logger = s.currentLogExpr.Exit(ctx)
		return scp
	})

	s.currentLogExpr = nil
}

func (s *ExpressionListener) scopeStart(_ antlr.TerminalNode) {
	if s.debug {
		fmt.Println("(")
	}

	if s.currentLogExpr != nil {
		s.currentLogExpr.ScopeStart()
	} else {
		scope := &Scope{}
		s.scopes.Push(scop{scp: scope, op: s.currentLogicalOp})

		s.currentLogicalOp = ConstraintLogicalLHS
	}
}

func (s *ExpressionListener) scopeEnd(_ antlr.TerminalNode) {
	if s.debug {
		fmt.Println(")")
	}

	if s.currentLogExpr != nil {
		s.currentLogExpr.ScopeEnd()
	} else {
		if s.currentPrimary != nil {
			s.scopes.Update(func(scp scop) scop {
				scp.scp.ScopeType = ScopeTypePrimaryExpr
				scp.scp.Primary = s.currentPrimary.Exit()

				return scp
			})

			s.currentPrimary = nil
		}

		if !s.scopes.IsEmpty() {
			current := s.scopes.Peek()
			s.currentLogicalOp = current.op
		}
	}
}

func (s *ExpressionListener) and(_ antlr.TerminalNode) {
	if s.debug {
		fmt.Println("AND")
	}

	if s.currentLogExpr != nil {
		s.currentLogExpr.And()
	} else {
		s.currentLogicalOp = ConstraintLogicalOpAnd
	}
}

func (s *ExpressionListener) or(_ antlr.TerminalNode) {
	if s.debug {
		fmt.Println("OR")
	}

	if s.currentLogExpr != nil {
		s.currentLogExpr.Or()
	} else {
		s.currentLogicalOp = ConstraintLogicalOpOr
	}
}

func (s *ExpressionListener) not(_ antlr.TerminalNode) {
	if s.debug {
		fmt.Println("NOT")
	}

	if s.currentLogExpr != nil {
		// NOT is not supported in the log expression
	} else {
		s.currentLogicalOp = ConstraintLogicalOpNot
	}
}
