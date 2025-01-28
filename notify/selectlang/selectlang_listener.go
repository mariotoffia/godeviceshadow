// Code generated from selectlang.g4 by ANTLR 4.13.2. DO NOT EDIT.

package selectlang // selectlang
import "github.com/antlr4-go/antlr/v4"

// selectlangListener is a complete listener for a parse tree produced by selectlangParser.
type selectlangListener interface {
	antlr.ParseTreeListener

	// EnterFilter is called when entering the filter production.
	EnterFilter(c *FilterContext)

	// EnterExpression is called when entering the expression production.
	EnterExpression(c *ExpressionContext)

	// EnterPrimaryExpr is called when entering the primaryExpr production.
	EnterPrimaryExpr(c *PrimaryExprContext)

	// EnterIdExpr is called when entering the idExpr production.
	EnterIdExpr(c *IdExprContext)

	// EnterNameExpr is called when entering the nameExpr production.
	EnterNameExpr(c *NameExprContext)

	// EnterOperationExpr is called when entering the operationExpr production.
	EnterOperationExpr(c *OperationExprContext)

	// EnterOperations is called when entering the operations production.
	EnterOperations(c *OperationsContext)

	// EnterLoggerExpr is called when entering the loggerExpr production.
	EnterLoggerExpr(c *LoggerExprContext)

	// EnterMapVarExpr is called when entering the mapVarExpr production.
	EnterMapVarExpr(c *MapVarExprContext)

	// EnterLoggerOp is called when entering the loggerOp production.
	EnterLoggerOp(c *LoggerOpContext)

	// EnterLoggerConstraints is called when entering the loggerConstraints production.
	EnterLoggerConstraints(c *LoggerConstraintsContext)

	// EnterValueComparison is called when entering the valueComparison production.
	EnterValueComparison(c *ValueComparisonContext)

	// EnterValueCondition is called when entering the valueCondition production.
	EnterValueCondition(c *ValueConditionContext)

	// EnterValueFactor is called when entering the valueFactor production.
	EnterValueFactor(c *ValueFactorContext)

	// EnterCompareOp is called when entering the compareOp production.
	EnterCompareOp(c *CompareOpContext)

	// EnterNumericLiteral is called when entering the NumericLiteral production.
	EnterNumericLiteral(c *NumericLiteralContext)

	// EnterStringLiteral is called when entering the StringLiteral production.
	EnterStringLiteral(c *StringLiteralContext)

	// EnterTimeLiteral is called when entering the TimeLiteral production.
	EnterTimeLiteral(c *TimeLiteralContext)

	// EnterRegexLiteral is called when entering the RegexLiteral production.
	EnterRegexLiteral(c *RegexLiteralContext)

	// EnterRegexOrString is called when entering the regexOrString production.
	EnterRegexOrString(c *RegexOrStringContext)

	// EnterRegex is called when entering the regex production.
	EnterRegex(c *RegexContext)

	// ExitFilter is called when exiting the filter production.
	ExitFilter(c *FilterContext)

	// ExitExpression is called when exiting the expression production.
	ExitExpression(c *ExpressionContext)

	// ExitPrimaryExpr is called when exiting the primaryExpr production.
	ExitPrimaryExpr(c *PrimaryExprContext)

	// ExitIdExpr is called when exiting the idExpr production.
	ExitIdExpr(c *IdExprContext)

	// ExitNameExpr is called when exiting the nameExpr production.
	ExitNameExpr(c *NameExprContext)

	// ExitOperationExpr is called when exiting the operationExpr production.
	ExitOperationExpr(c *OperationExprContext)

	// ExitOperations is called when exiting the operations production.
	ExitOperations(c *OperationsContext)

	// ExitLoggerExpr is called when exiting the loggerExpr production.
	ExitLoggerExpr(c *LoggerExprContext)

	// ExitMapVarExpr is called when exiting the mapVarExpr production.
	ExitMapVarExpr(c *MapVarExprContext)

	// ExitLoggerOp is called when exiting the loggerOp production.
	ExitLoggerOp(c *LoggerOpContext)

	// ExitLoggerConstraints is called when exiting the loggerConstraints production.
	ExitLoggerConstraints(c *LoggerConstraintsContext)

	// ExitValueComparison is called when exiting the valueComparison production.
	ExitValueComparison(c *ValueComparisonContext)

	// ExitValueCondition is called when exiting the valueCondition production.
	ExitValueCondition(c *ValueConditionContext)

	// ExitValueFactor is called when exiting the valueFactor production.
	ExitValueFactor(c *ValueFactorContext)

	// ExitCompareOp is called when exiting the compareOp production.
	ExitCompareOp(c *CompareOpContext)

	// ExitNumericLiteral is called when exiting the NumericLiteral production.
	ExitNumericLiteral(c *NumericLiteralContext)

	// ExitStringLiteral is called when exiting the StringLiteral production.
	ExitStringLiteral(c *StringLiteralContext)

	// ExitTimeLiteral is called when exiting the TimeLiteral production.
	ExitTimeLiteral(c *TimeLiteralContext)

	// ExitRegexLiteral is called when exiting the RegexLiteral production.
	ExitRegexLiteral(c *RegexLiteralContext)

	// ExitRegexOrString is called when exiting the regexOrString production.
	ExitRegexOrString(c *RegexOrStringContext)

	// ExitRegex is called when exiting the regex production.
	ExitRegex(c *RegexContext)
}
