// Code generated from selectlang.g4 by ANTLR 4.13.2. DO NOT EDIT.

package selectlang // selectlang
import "github.com/antlr4-go/antlr/v4"

// BaseselectlangListener is a complete listener for a parse tree produced by selectlangParser.
type BaseselectlangListener struct{}

var _ selectlangListener = &BaseselectlangListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseselectlangListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseselectlangListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseselectlangListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseselectlangListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterFilter is called when production filter is entered.
func (s *BaseselectlangListener) EnterFilter(ctx *FilterContext) {}

// ExitFilter is called when production filter is exited.
func (s *BaseselectlangListener) ExitFilter(ctx *FilterContext) {}

// EnterExpression is called when production expression is entered.
func (s *BaseselectlangListener) EnterExpression(ctx *ExpressionContext) {}

// ExitExpression is called when production expression is exited.
func (s *BaseselectlangListener) ExitExpression(ctx *ExpressionContext) {}

// EnterPrimaryExpr is called when production primaryExpr is entered.
func (s *BaseselectlangListener) EnterPrimaryExpr(ctx *PrimaryExprContext) {}

// ExitPrimaryExpr is called when production primaryExpr is exited.
func (s *BaseselectlangListener) ExitPrimaryExpr(ctx *PrimaryExprContext) {}

// EnterIdExpr is called when production idExpr is entered.
func (s *BaseselectlangListener) EnterIdExpr(ctx *IdExprContext) {}

// ExitIdExpr is called when production idExpr is exited.
func (s *BaseselectlangListener) ExitIdExpr(ctx *IdExprContext) {}

// EnterNameExpr is called when production nameExpr is entered.
func (s *BaseselectlangListener) EnterNameExpr(ctx *NameExprContext) {}

// ExitNameExpr is called when production nameExpr is exited.
func (s *BaseselectlangListener) ExitNameExpr(ctx *NameExprContext) {}

// EnterOperationExpr is called when production operationExpr is entered.
func (s *BaseselectlangListener) EnterOperationExpr(ctx *OperationExprContext) {}

// ExitOperationExpr is called when production operationExpr is exited.
func (s *BaseselectlangListener) ExitOperationExpr(ctx *OperationExprContext) {}

// EnterOperations is called when production operations is entered.
func (s *BaseselectlangListener) EnterOperations(ctx *OperationsContext) {}

// ExitOperations is called when production operations is exited.
func (s *BaseselectlangListener) ExitOperations(ctx *OperationsContext) {}

// EnterLoggerExpr is called when production loggerExpr is entered.
func (s *BaseselectlangListener) EnterLoggerExpr(ctx *LoggerExprContext) {}

// ExitLoggerExpr is called when production loggerExpr is exited.
func (s *BaseselectlangListener) ExitLoggerExpr(ctx *LoggerExprContext) {}

// EnterMapVarExpr is called when production mapVarExpr is entered.
func (s *BaseselectlangListener) EnterMapVarExpr(ctx *MapVarExprContext) {}

// ExitMapVarExpr is called when production mapVarExpr is exited.
func (s *BaseselectlangListener) ExitMapVarExpr(ctx *MapVarExprContext) {}

// EnterLoggerOp is called when production loggerOp is entered.
func (s *BaseselectlangListener) EnterLoggerOp(ctx *LoggerOpContext) {}

// ExitLoggerOp is called when production loggerOp is exited.
func (s *BaseselectlangListener) ExitLoggerOp(ctx *LoggerOpContext) {}

// EnterValueComparison is called when production valueComparison is entered.
func (s *BaseselectlangListener) EnterValueComparison(ctx *ValueComparisonContext) {}

// ExitValueComparison is called when production valueComparison is exited.
func (s *BaseselectlangListener) ExitValueComparison(ctx *ValueComparisonContext) {}

// EnterValueCondition is called when production valueCondition is entered.
func (s *BaseselectlangListener) EnterValueCondition(ctx *ValueConditionContext) {}

// ExitValueCondition is called when production valueCondition is exited.
func (s *BaseselectlangListener) ExitValueCondition(ctx *ValueConditionContext) {}

// EnterValueFactor is called when production valueFactor is entered.
func (s *BaseselectlangListener) EnterValueFactor(ctx *ValueFactorContext) {}

// ExitValueFactor is called when production valueFactor is exited.
func (s *BaseselectlangListener) ExitValueFactor(ctx *ValueFactorContext) {}

// EnterCompareOp is called when production compareOp is entered.
func (s *BaseselectlangListener) EnterCompareOp(ctx *CompareOpContext) {}

// ExitCompareOp is called when production compareOp is exited.
func (s *BaseselectlangListener) ExitCompareOp(ctx *CompareOpContext) {}

// EnterNumericLiteral is called when production NumericLiteral is entered.
func (s *BaseselectlangListener) EnterNumericLiteral(ctx *NumericLiteralContext) {}

// ExitNumericLiteral is called when production NumericLiteral is exited.
func (s *BaseselectlangListener) ExitNumericLiteral(ctx *NumericLiteralContext) {}

// EnterStringLiteral is called when production StringLiteral is entered.
func (s *BaseselectlangListener) EnterStringLiteral(ctx *StringLiteralContext) {}

// ExitStringLiteral is called when production StringLiteral is exited.
func (s *BaseselectlangListener) ExitStringLiteral(ctx *StringLiteralContext) {}

// EnterTimeLiteral is called when production TimeLiteral is entered.
func (s *BaseselectlangListener) EnterTimeLiteral(ctx *TimeLiteralContext) {}

// ExitTimeLiteral is called when production TimeLiteral is exited.
func (s *BaseselectlangListener) ExitTimeLiteral(ctx *TimeLiteralContext) {}

// EnterRegexLiteral is called when production RegexLiteral is entered.
func (s *BaseselectlangListener) EnterRegexLiteral(ctx *RegexLiteralContext) {}

// ExitRegexLiteral is called when production RegexLiteral is exited.
func (s *BaseselectlangListener) ExitRegexLiteral(ctx *RegexLiteralContext) {}

// EnterRegexOrString is called when production regexOrString is entered.
func (s *BaseselectlangListener) EnterRegexOrString(ctx *RegexOrStringContext) {}

// ExitRegexOrString is called when production regexOrString is exited.
func (s *BaseselectlangListener) ExitRegexOrString(ctx *RegexOrStringContext) {}

// EnterRegex is called when production regex is entered.
func (s *BaseselectlangListener) EnterRegex(ctx *RegexContext) {}

// ExitRegex is called when production regex is exited.
func (s *BaseselectlangListener) ExitRegex(ctx *RegexContext) {}
