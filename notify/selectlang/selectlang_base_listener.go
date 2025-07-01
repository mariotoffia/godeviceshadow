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

// EnterSelectStatement is called when production SelectStatement is entered.
func (s *BaseselectlangListener) EnterSelectStatement(ctx *SelectStatementContext) {}

// ExitSelectStatement is called when production SelectStatement is exited.
func (s *BaseselectlangListener) ExitSelectStatement(ctx *SelectStatementContext) {}

// EnterAllColumns is called when production AllColumns is entered.
func (s *BaseselectlangListener) EnterAllColumns(ctx *AllColumnsContext) {}

// ExitAllColumns is called when production AllColumns is exited.
func (s *BaseselectlangListener) ExitAllColumns(ctx *AllColumnsContext) {}

// EnterStreamName is called when production StreamName is entered.
func (s *BaseselectlangListener) EnterStreamName(ctx *StreamNameContext) {}

// ExitStreamName is called when production StreamName is exited.
func (s *BaseselectlangListener) ExitStreamName(ctx *StreamNameContext) {}

// EnterWhereClause is called when production WhereClause is entered.
func (s *BaseselectlangListener) EnterWhereClause(ctx *WhereClauseContext) {}

// ExitWhereClause is called when production WhereClause is exited.
func (s *BaseselectlangListener) ExitWhereClause(ctx *WhereClauseContext) {}

// EnterAndToExpression is called when production AndToExpression is entered.
func (s *BaseselectlangListener) EnterAndToExpression(ctx *AndToExpressionContext) {}

// ExitAndToExpression is called when production AndToExpression is exited.
func (s *BaseselectlangListener) ExitAndToExpression(ctx *AndToExpressionContext) {}

// EnterOrExpression is called when production OrExpression is entered.
func (s *BaseselectlangListener) EnterOrExpression(ctx *OrExpressionContext) {}

// ExitOrExpression is called when production OrExpression is exited.
func (s *BaseselectlangListener) ExitOrExpression(ctx *OrExpressionContext) {}

// EnterAndExpression is called when production AndExpression is entered.
func (s *BaseselectlangListener) EnterAndExpression(ctx *AndExpressionContext) {}

// ExitAndExpression is called when production AndExpression is exited.
func (s *BaseselectlangListener) ExitAndExpression(ctx *AndExpressionContext) {}

// EnterPrimaryExpression is called when production PrimaryExpression is entered.
func (s *BaseselectlangListener) EnterPrimaryExpression(ctx *PrimaryExpressionContext) {}

// ExitPrimaryExpression is called when production PrimaryExpression is exited.
func (s *BaseselectlangListener) ExitPrimaryExpression(ctx *PrimaryExpressionContext) {}

// EnterParenExpression is called when production ParenExpression is entered.
func (s *BaseselectlangListener) EnterParenExpression(ctx *ParenExpressionContext) {}

// ExitParenExpression is called when production ParenExpression is exited.
func (s *BaseselectlangListener) ExitParenExpression(ctx *ParenExpressionContext) {}

// EnterPredicateExpression is called when production PredicateExpression is entered.
func (s *BaseselectlangListener) EnterPredicateExpression(ctx *PredicateExpressionContext) {}

// ExitPredicateExpression is called when production PredicateExpression is exited.
func (s *BaseselectlangListener) ExitPredicateExpression(ctx *PredicateExpressionContext) {}

// EnterComparisonPredicate is called when production ComparisonPredicate is entered.
func (s *BaseselectlangListener) EnterComparisonPredicate(ctx *ComparisonPredicateContext) {}

// ExitComparisonPredicate is called when production ComparisonPredicate is exited.
func (s *BaseselectlangListener) ExitComparisonPredicate(ctx *ComparisonPredicateContext) {}

// EnterRegexPredicate is called when production RegexPredicate is entered.
func (s *BaseselectlangListener) EnterRegexPredicate(ctx *RegexPredicateContext) {}

// ExitRegexPredicate is called when production RegexPredicate is exited.
func (s *BaseselectlangListener) ExitRegexPredicate(ctx *RegexPredicateContext) {}

// EnterInPredicate is called when production InPredicate is entered.
func (s *BaseselectlangListener) EnterInPredicate(ctx *InPredicateContext) {}

// ExitInPredicate is called when production InPredicate is exited.
func (s *BaseselectlangListener) ExitInPredicate(ctx *InPredicateContext) {}

// EnterValueList is called when production ValueList is entered.
func (s *BaseselectlangListener) EnterValueList(ctx *ValueListContext) {}

// ExitValueList is called when production ValueList is exited.
func (s *BaseselectlangListener) ExitValueList(ctx *ValueListContext) {}

// EnterObjField is called when production ObjField is entered.
func (s *BaseselectlangListener) EnterObjField(ctx *ObjFieldContext) {}

// ExitObjField is called when production ObjField is exited.
func (s *BaseselectlangListener) ExitObjField(ctx *ObjFieldContext) {}

// EnterLogField is called when production LogField is entered.
func (s *BaseselectlangListener) EnterLogField(ctx *LogFieldContext) {}

// ExitLogField is called when production LogField is exited.
func (s *BaseselectlangListener) ExitLogField(ctx *LogFieldContext) {}

// EnterObjFieldAccess is called when production ObjFieldAccess is entered.
func (s *BaseselectlangListener) EnterObjFieldAccess(ctx *ObjFieldAccessContext) {}

// ExitObjFieldAccess is called when production ObjFieldAccess is exited.
func (s *BaseselectlangListener) ExitObjFieldAccess(ctx *ObjFieldAccessContext) {}

// EnterLogFieldAccess is called when production LogFieldAccess is entered.
func (s *BaseselectlangListener) EnterLogFieldAccess(ctx *LogFieldAccessContext) {}

// ExitLogFieldAccess is called when production LogFieldAccess is exited.
func (s *BaseselectlangListener) ExitLogFieldAccess(ctx *LogFieldAccessContext) {}

// EnterNumberValue is called when production NumberValue is entered.
func (s *BaseselectlangListener) EnterNumberValue(ctx *NumberValueContext) {}

// ExitNumberValue is called when production NumberValue is exited.
func (s *BaseselectlangListener) ExitNumberValue(ctx *NumberValueContext) {}

// EnterStringValue is called when production StringValue is entered.
func (s *BaseselectlangListener) EnterStringValue(ctx *StringValueContext) {}

// ExitStringValue is called when production StringValue is exited.
func (s *BaseselectlangListener) ExitStringValue(ctx *StringValueContext) {}

// EnterEqualsOp is called when production EqualsOp is entered.
func (s *BaseselectlangListener) EnterEqualsOp(ctx *EqualsOpContext) {}

// ExitEqualsOp is called when production EqualsOp is exited.
func (s *BaseselectlangListener) ExitEqualsOp(ctx *EqualsOpContext) {}

// EnterNotEqualsOp is called when production NotEqualsOp is entered.
func (s *BaseselectlangListener) EnterNotEqualsOp(ctx *NotEqualsOpContext) {}

// ExitNotEqualsOp is called when production NotEqualsOp is exited.
func (s *BaseselectlangListener) ExitNotEqualsOp(ctx *NotEqualsOpContext) {}

// EnterGreaterThanOp is called when production GreaterThanOp is entered.
func (s *BaseselectlangListener) EnterGreaterThanOp(ctx *GreaterThanOpContext) {}

// ExitGreaterThanOp is called when production GreaterThanOp is exited.
func (s *BaseselectlangListener) ExitGreaterThanOp(ctx *GreaterThanOpContext) {}

// EnterLessThanOp is called when production LessThanOp is entered.
func (s *BaseselectlangListener) EnterLessThanOp(ctx *LessThanOpContext) {}

// ExitLessThanOp is called when production LessThanOp is exited.
func (s *BaseselectlangListener) ExitLessThanOp(ctx *LessThanOpContext) {}

// EnterGreaterOrEqualOp is called when production GreaterOrEqualOp is entered.
func (s *BaseselectlangListener) EnterGreaterOrEqualOp(ctx *GreaterOrEqualOpContext) {}

// ExitGreaterOrEqualOp is called when production GreaterOrEqualOp is exited.
func (s *BaseselectlangListener) ExitGreaterOrEqualOp(ctx *GreaterOrEqualOpContext) {}

// EnterLessOrEqualOp is called when production LessOrEqualOp is entered.
func (s *BaseselectlangListener) EnterLessOrEqualOp(ctx *LessOrEqualOpContext) {}

// ExitLessOrEqualOp is called when production LessOrEqualOp is exited.
func (s *BaseselectlangListener) ExitLessOrEqualOp(ctx *LessOrEqualOpContext) {}

// EnterRegexOp is called when production RegexOp is entered.
func (s *BaseselectlangListener) EnterRegexOp(ctx *RegexOpContext) {}

// ExitRegexOp is called when production RegexOp is exited.
func (s *BaseselectlangListener) ExitRegexOp(ctx *RegexOpContext) {}

// EnterRegexValue is called when production RegexValue is entered.
func (s *BaseselectlangListener) EnterRegexValue(ctx *RegexValueContext) {}

// ExitRegexValue is called when production RegexValue is exited.
func (s *BaseselectlangListener) ExitRegexValue(ctx *RegexValueContext) {}
