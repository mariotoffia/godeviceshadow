package selectlang

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"

	"github.com/mariotoffia/godeviceshadow/model/notifiermodel"
)

func ToSelection(expr string) (notifiermodel.Selection, error) {
	//
	var process func(scope *Scope) (notifiermodel.Selection, error)

	process = func(scope *Scope) (notifiermodel.Selection, error) {
		ands := make([]notifiermodel.Selection, 0, len(scope.And))
		andss := make([]notifiermodel.Selection, 0, len(scope.And))
		ors := make([]notifiermodel.Selection, 0, len(scope.Or))
		orss := make([]notifiermodel.Selection, 0, len(scope.Or))

		if len(scope.And) > 0 {
			for i := 0; i < len(scope.And); i++ {
				and := scope.And[i]
				sel, err := process(scope.And[i])

				if err != nil {
					return nil, err
				}

				if len(and.And) > 0 || len(and.Or) > 0 {
					andss = append(andss, sel)
				} else {
					ands = append(ands, sel)
				}
			}
		}

		if len(scope.Or) > 0 {
			for i := 0; i < len(scope.Or); i++ {
				or := scope.Or[i]
				sel, err := process(scope.Or[i])

				if err != nil {
					return nil, err
				}

				if len(or.And) > 0 || len(or.Or) > 0 {
					orss = append(orss, sel)
				} else {
					ors = append(ors, sel)
				}
			}
		}

		mf := scope.ToMatchFunc()

		if mf == nil {
			return nil, fmt.Errorf("must have either primary expression or log expression in a scope")
		}

		sb := notifiermodel.NewSelectionBuilder(notifiermodel.Func(mf))

		if len(ands) > 0 {
			sb.And(ands...)
		}

		for _, and := range andss {
			sb.And(and) // separate scope
		}

		if len(ors) > 0 {
			sb.Or(ors...)
		}

		for _, or := range orss {
			sb.Or(or) // separate scope
		}

		return sb.Build()
	}

	scope := ToScope(expr)

	return process(&scope)
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
