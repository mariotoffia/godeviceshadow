package selectlang

import (
	"fmt"
)

// PrimaryExprListener handles parsing of primary expressions (ID, Name, Operations)
type PrimaryExprListener struct {
	debug   bool
	primary *PrimaryExpression
}

func NewPrimaryExprListener(debug bool) *PrimaryExprListener {
	return &PrimaryExprListener{
		debug:   debug,
		primary: &PrimaryExpression{},
	}
}

func (p *PrimaryExprListener) Enter() {
	if p.debug {
		fmt.Println("PRIMARY ENTER")
	}
}

func (p *PrimaryExprListener) Exit() *PrimaryExpression {
	if p.debug {
		fmt.Println("PRIMARY EXIT")
	}

	return p.primary
}

func (p *PrimaryExprListener) EnterIdExpr(ctx *IdExprContext) {
	if child := FirstChild[*RegexOrStringContext](ctx); child != nil {
		if p.debug {
			fmt.Println("    ID:", child.GetText())
		}

		p.primary.ID = child.GetText()
	}
}

func (p *PrimaryExprListener) EnterNameExpr(ctx *NameExprContext) {
	if child := FirstChild[*RegexOrStringContext](ctx); child != nil {
		if p.debug {
			fmt.Println("    NAME:", child.GetText())
		}

		p.primary.Name = child.GetText()
	}
}

func (p *PrimaryExprListener) EnterOperationExpr(ctx *OperationExprContext) {
	operation := ctx.Operations()

	if operation == nil {
		return
	}

	p.primary.Operation = ToStringList(operation, ",")

	if p.debug {
		fmt.Printf("    OPERATION:%#v\n", p.primary.Operation)
	}
}
