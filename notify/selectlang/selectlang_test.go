package selectlang_test

import (
	"fmt"
	"testing"

	"github.com/antlr4-go/antlr/v4"

	"github.com/mariotoffia/godeviceshadow/notify/selectlang"
)

type ExpressionListener struct {
	*selectlang.BaseselectlangListener

	insideLoggerExpression  bool
	insideLoggerConstraints bool
}

func NewTreeShapeListener() *ExpressionListener {
	return new(ExpressionListener)
}

func (s *ExpressionListener) VisitTerminal(node antlr.TerminalNode) {
	term := node.GetText()

	if term == "AND" {
		fmt.Println("AND")
	} else if term == "OR" {
		fmt.Println("OR")
	} else if term == "(" {
		fmt.Println("(")
	} else if term == ")" {
		fmt.Println(")")
	}
}

func (s *ExpressionListener) EnterIdExpr(ctx *selectlang.IdExprContext) {
	if child := selectlang.FirstChild[*selectlang.RegexOrStringContext](ctx); child != nil {
		fmt.Println("ID:", child.GetText())
	}
}

func (s *ExpressionListener) EnterNameExpr(ctx *selectlang.NameExprContext) {
	if child := selectlang.FirstChild[*selectlang.RegexOrStringContext](ctx); child != nil {
		fmt.Println("NAME:", child.GetText())
	}
}

func (s *ExpressionListener) EnterOperationExpr(ctx *selectlang.OperationExprContext) {
	operation := ctx.Operations()

	if operation == nil {
		return
	}

	fmt.Printf("    OPERATION:%#v\n", selectlang.ToStringList(operation, ","))
}

func (s *ExpressionListener) EnterLoggerExpr(ctx *selectlang.LoggerExprContext) {
	s.insideLoggerExpression = true

	fmt.Println("LOGGER START ")

	if child := selectlang.FirstChild[*selectlang.LoggerOpContext](ctx); child != nil {
		fmt.Printf("    LOGGER_OP:%#v\n", selectlang.ToStringList(child, ","))
	}

	if child := selectlang.FirstChild[*selectlang.RegexContext](ctx); child != nil {
		fmt.Println("    REGEX:", child.REGEX().GetText())
	}

	if child := selectlang.FirstChild[*selectlang.MapVarExprContext](ctx); child != nil {
		// child.EQ()
		fmt.Println("    MAP_VAR_EXPR:", child.STRING().GetText())
	}
}

func (s *ExpressionListener) EnterLoggerConstraints(ctx *selectlang.LoggerConstraintsContext) {
	s.insideLoggerConstraints = true

	fmt.Println("   ENTER_LOGGER_CONSTRAINTS")
}

func (s *ExpressionListener) ExitLoggerConstraints(ctx *selectlang.LoggerConstraintsContext) {
	s.insideLoggerConstraints = false

	fmt.Println("   EXIT_LOGGER_CONSTRAINTS")
}

func (s *ExpressionListener) EnterValueFactor(ctx *selectlang.ValueFactorContext) {
	if compareOp := ctx.CompareOp(); compareOp != nil {
		fmt.Println("    VARIABLE: var")

		switch {
		case compareOp.GT() != nil:
			fmt.Println("    COMPARE_OP: GT")
		case compareOp.GE() != nil:
			fmt.Println("    COMPARE_OP: GE")
		case compareOp.LT() != nil:
			fmt.Println("    COMPARE_OP: LT")
		case compareOp.LE() != nil:
			fmt.Println("    COMPARE_OP: LE")
		case compareOp.EQ() != nil:
			fmt.Println("    COMPARE_OP: EQ")
		case compareOp.NE() != nil:
			fmt.Println("    COMPARE_OP: NE")
		case compareOp.BEFORE() != nil:
			fmt.Println("    COMPARE_OP: BEFORE")
		case compareOp.AFTER() != nil:
			fmt.Println("    COMPARE_OP: AFTER")
		default:
			fmt.Println("    COMPARE_OP:", compareOp.GetText())
		}
	}

	if cr := ctx.ConstantOrRegex(); cr != nil {
		switch t := cr.(type) {
		case *selectlang.NumericLiteralContext:
			fmt.Println("    NUMBER:", t.NUMBER().GetText())
		case *selectlang.StringLiteralContext:
			fmt.Println("    STRING:", t.STRING().GetText())
		case *selectlang.RegexLiteralContext:
			fmt.Println("    REGEX:", t.Regex().REGEX().GetText())
		}
	}
}

func (s *ExpressionListener) ExitLoggerExpr(ctx *selectlang.LoggerExprContext) {
	s.insideLoggerExpression = false

	fmt.Println("LOGGER END")
}

func TestSimpleIdNameExpression(t *testing.T) {
	stmt := `
		(
			id: /myDevice-\d+/ AND 
			name: 'myShadow' AND 
			operation: report,desired
		) 
		AND 
		(add,update:/^Sensors-.*-indoor$/ == 'temp' WHERE
				(value > 20 OR (value < /re-\d+/ AND value != 'apa'))
		) 
		OR 
		(add,update)
	`

	input := antlr.NewInputStream(stmt)
	lexer := selectlang.NewselectlangLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := selectlang.NewselectlangParser(stream)
	p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
	tree := p.Filter()

	//fmt.Println(tree.ToStringTree(nil, p))

	antlr.ParseTreeWalkerDefault.Walk(NewTreeShapeListener(), tree)
}
