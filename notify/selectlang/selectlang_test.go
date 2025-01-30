package selectlang_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/antlr4-go/antlr/v4"

	"github.com/mariotoffia/godeviceshadow/notify/selectlang"
)

func TestSimpleIdNameExpression(t *testing.T) {
	stmt := `
		(
			id: /myDevice-\d+/ AND 
			name: 'myShadow' AND 
			operation: report,desired
		)
		AND
		(add,update:/^Sensors-.*-indoor$/ == 'temp'  
		WHERE (
			value > 20 OR (value < /re-\d+/ AND value != 'apa'))
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

	listener := selectlang.NewExpressionListener(selectlang.ExpressionListenerOpts{
		Debug: false,
	})

	antlr.ParseTreeWalkerDefault.Walk(listener, tree)

	var buff bytes.Buffer
	listener.RootScope().Dump(&buff, 0 /*indent*/)
	fmt.Println()
	fmt.Println("-----------------------------")
	fmt.Println(buff.String())
}
