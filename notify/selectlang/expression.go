package selectlang

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/mariotoffia/godeviceshadow/model/notifiermodel"
)

func ToSelectionBuilder(expr string) *notifiermodel.SelectionBuilder {
	/*
	   Example:
	   	sb := notifiermodel.NewSelectionBuilder(
	   		notifiermodel.Scoped(&a, func(sb *notifiermodel.SelectionBuilder) {
	   			sb.Or(notifiermodel.Scoped(&b, func(sb *notifiermodel.SelectionBuilder) {
	   				sb.And(notifiermodel.Scoped(&c, func(sb *notifiermodel.SelectionBuilder) {
	   					sb.Or(&d)
	   				}))
	   			}))
	   		})).
	   		And(&e).
	   		Or(notifiermodel.Scoped(&f, func(sb *notifiermodel.SelectionBuilder) {
	   			sb.And(notifiermodel.Scoped(&g, func(sb *notifiermodel.SelectionBuilder) {
	   				sb.Or(&h)
	   			}))
	   		}))
	*/

	//scope := ToScope(expr)

	return nil
}

func ToScope(expr string) Scope {
	input := antlr.NewInputStream(expr)
	lexer := NewselectlangLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := NewselectlangParser(stream)

	p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))

	listener := NewExpressionListener(ExpressionListenerOpts{
		Debug: false,
	})

	antlr.ParseTreeWalkerDefault.Walk(listener, p.Filter())

	return listener.RootScope()
}
