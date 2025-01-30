package selectlang

// Scope encapsulates a '(' and ')' pair with type info and
// collects the expression within the scope.
type Scope struct {
	// ScopeType is the type of scope. When the '(' is hit it will be
	// set to `ScopeLoggerUntyped` and then, depending on the content
	// it will be set to correct type.
	ScopeType ScopeType
	Primary   *PrimaryExpression
	Logger    *LoggerExpression
	And       []*Scope
	Or        []*Scope
	Not       *Scope
}

type ScopeType int

const (
	ScopeLoggerUntyped ScopeType = iota
	ScopeLoggerExpr
	ScopeTypePrimaryExpr
)

func (scope Scope) Children() []*Scope {
	res := append(scope.And, scope.Or...)

	if scope.Not != nil {
		return append(res, scope.Not)
	}

	return res
}
