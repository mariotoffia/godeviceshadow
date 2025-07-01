package selectlang_test

import (
	"testing"

	antlr "github.com/antlr4-go/antlr/v4"

	"github.com/mariotoffia/godeviceshadow/notify/selectlang"
)

func TestComplexExpression(t *testing.T) {
	stmt := `
        SELECT * FROM Notification WHERE
        (
            obj.ID ~= 'myDevice-\\d+' AND
            obj.Name == 'myShadow' AND
            obj.Operation IN ('report','desired')
        )
        AND
        (
            log.Operation IN ('add','update') AND
            log.Path ~= '^Sensors-.*-indoor$' AND
            log.Value == 'temp' AND
            (
                log.Value > 20 OR (log.Value ~= 're-\\d+' AND log.Value != 'apa' OR (log.Value > 99 AND log.Value != 'bubben-\\d+'))
            )
        )
        OR
        (log.Operation IN ('add','update'))
    `

	input := antlr.NewInputStream(stmt)
	lexer := selectlang.NewselectlangLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := selectlang.NewselectlangParser(stream)
	p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
	//tree := p.Select_stmt()

	//fmt.Println(tree.ToStringTree(nil, p))

	/*
	   	listener := selectlang.NewExpressionListener(selectlang.ExpressionListenerOpts{
	   		Debug: false,
	   	})

	   antlr.ParseTreeWalkerDefault.Walk(listener, tree)

	   var buff bytes.Buffer
	   listener.RootScope().Dump(&buff, 0)
	   fmt.Println()
	   fmt.Println("-----------------------------")
	   fmt.Println(buff.String())
	*/
}
