package selectlang_test

import (
	"fmt"
	"testing"

	"github.com/antlr4-go/antlr/v4"
	"github.com/mariotoffia/godeviceshadow/notify/selectlang"
)

type ExpressionListener struct {
	*selectlang.BaseselectlangListener
}

func NewTreeShapeListener() *ExpressionListener {
	return new(ExpressionListener)
}

func (this *ExpressionListener) EnterEveryRule(ctx antlr.ParserRuleContext) {
	fmt.Println(ctx.GetText())
}

func TestSimpleIdNameExpression(t *testing.T) {
	input := antlr.NewInputStream(`(id: /myDevice-\d+/ AND name: 'myShadow')`)
	lexer := selectlang.NewselectlangLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := selectlang.NewselectlangParser(stream)
	p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
	tree := p.Filter()
	antlr.ParseTreeWalkerDefault.Walk(NewTreeShapeListener(), tree)

}
