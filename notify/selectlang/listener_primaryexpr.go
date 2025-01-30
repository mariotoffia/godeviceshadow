package selectlang

import "fmt"

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
