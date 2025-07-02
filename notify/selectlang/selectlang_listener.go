// Code generated from selectlang.g4 by ANTLR 4.13.2. DO NOT EDIT.

package selectlang // selectlang
import "github.com/antlr4-go/antlr/v4"

// selectlangListener is a complete listener for a parse tree produced by selectlangParser.
type selectlangListener interface {
	antlr.ParseTreeListener

	// EnterSelectStatement is called when entering the SelectStatement production.
	EnterSelectStatement(c *SelectStatementContext)

	// EnterAllColumns is called when entering the AllColumns production.
	EnterAllColumns(c *AllColumnsContext)

	// EnterStreamName is called when entering the StreamName production.
	EnterStreamName(c *StreamNameContext)

	// EnterWhereClause is called when entering the WhereClause production.
	EnterWhereClause(c *WhereClauseContext)

	// EnterAndToExpression is called when entering the AndToExpression production.
	EnterAndToExpression(c *AndToExpressionContext)

	// EnterOrExpression is called when entering the OrExpression production.
	EnterOrExpression(c *OrExpressionContext)

	// EnterAndExpression is called when entering the AndExpression production.
	EnterAndExpression(c *AndExpressionContext)

	// EnterPrimaryExpression is called when entering the PrimaryExpression production.
	EnterPrimaryExpression(c *PrimaryExpressionContext)

	// EnterParenExpression is called when entering the ParenExpression production.
	EnterParenExpression(c *ParenExpressionContext)

	// EnterPredicateExpression is called when entering the PredicateExpression production.
	EnterPredicateExpression(c *PredicateExpressionContext)

	// EnterComparisonPredicate is called when entering the ComparisonPredicate production.
	EnterComparisonPredicate(c *ComparisonPredicateContext)

	// EnterRegexPredicate is called when entering the RegexPredicate production.
	EnterRegexPredicate(c *RegexPredicateContext)

	// EnterInPredicate is called when entering the InPredicate production.
	EnterInPredicate(c *InPredicateContext)

	// EnterHasPredicate is called when entering the HasPredicate production.
	EnterHasPredicate(c *HasPredicateContext)

	// EnterValueList is called when entering the ValueList production.
	EnterValueList(c *ValueListContext)

	// EnterObjField is called when entering the ObjField production.
	EnterObjField(c *ObjFieldContext)

	// EnterLogField is called when entering the LogField production.
	EnterLogField(c *LogFieldContext)

	// EnterObjFieldAccess is called when entering the ObjFieldAccess production.
	EnterObjFieldAccess(c *ObjFieldAccessContext)

	// EnterLogFieldAccess is called when entering the LogFieldAccess production.
	EnterLogFieldAccess(c *LogFieldAccessContext)

	// EnterNumberValue is called when entering the NumberValue production.
	EnterNumberValue(c *NumberValueContext)

	// EnterStringValue is called when entering the StringValue production.
	EnterStringValue(c *StringValueContext)

	// EnterEqualsOp is called when entering the EqualsOp production.
	EnterEqualsOp(c *EqualsOpContext)

	// EnterNotEqualsOp is called when entering the NotEqualsOp production.
	EnterNotEqualsOp(c *NotEqualsOpContext)

	// EnterGreaterThanOp is called when entering the GreaterThanOp production.
	EnterGreaterThanOp(c *GreaterThanOpContext)

	// EnterLessThanOp is called when entering the LessThanOp production.
	EnterLessThanOp(c *LessThanOpContext)

	// EnterGreaterOrEqualOp is called when entering the GreaterOrEqualOp production.
	EnterGreaterOrEqualOp(c *GreaterOrEqualOpContext)

	// EnterLessOrEqualOp is called when entering the LessOrEqualOp production.
	EnterLessOrEqualOp(c *LessOrEqualOpContext)

	// EnterRegexOp is called when entering the RegexOp production.
	EnterRegexOp(c *RegexOpContext)

	// EnterRegexNotOp is called when entering the RegexNotOp production.
	EnterRegexNotOp(c *RegexNotOpContext)

	// EnterRegexValue is called when entering the RegexValue production.
	EnterRegexValue(c *RegexValueContext)

	// ExitSelectStatement is called when exiting the SelectStatement production.
	ExitSelectStatement(c *SelectStatementContext)

	// ExitAllColumns is called when exiting the AllColumns production.
	ExitAllColumns(c *AllColumnsContext)

	// ExitStreamName is called when exiting the StreamName production.
	ExitStreamName(c *StreamNameContext)

	// ExitWhereClause is called when exiting the WhereClause production.
	ExitWhereClause(c *WhereClauseContext)

	// ExitAndToExpression is called when exiting the AndToExpression production.
	ExitAndToExpression(c *AndToExpressionContext)

	// ExitOrExpression is called when exiting the OrExpression production.
	ExitOrExpression(c *OrExpressionContext)

	// ExitAndExpression is called when exiting the AndExpression production.
	ExitAndExpression(c *AndExpressionContext)

	// ExitPrimaryExpression is called when exiting the PrimaryExpression production.
	ExitPrimaryExpression(c *PrimaryExpressionContext)

	// ExitParenExpression is called when exiting the ParenExpression production.
	ExitParenExpression(c *ParenExpressionContext)

	// ExitPredicateExpression is called when exiting the PredicateExpression production.
	ExitPredicateExpression(c *PredicateExpressionContext)

	// ExitComparisonPredicate is called when exiting the ComparisonPredicate production.
	ExitComparisonPredicate(c *ComparisonPredicateContext)

	// ExitRegexPredicate is called when exiting the RegexPredicate production.
	ExitRegexPredicate(c *RegexPredicateContext)

	// ExitInPredicate is called when exiting the InPredicate production.
	ExitInPredicate(c *InPredicateContext)

	// ExitHasPredicate is called when exiting the HasPredicate production.
	ExitHasPredicate(c *HasPredicateContext)

	// ExitValueList is called when exiting the ValueList production.
	ExitValueList(c *ValueListContext)

	// ExitObjField is called when exiting the ObjField production.
	ExitObjField(c *ObjFieldContext)

	// ExitLogField is called when exiting the LogField production.
	ExitLogField(c *LogFieldContext)

	// ExitObjFieldAccess is called when exiting the ObjFieldAccess production.
	ExitObjFieldAccess(c *ObjFieldAccessContext)

	// ExitLogFieldAccess is called when exiting the LogFieldAccess production.
	ExitLogFieldAccess(c *LogFieldAccessContext)

	// ExitNumberValue is called when exiting the NumberValue production.
	ExitNumberValue(c *NumberValueContext)

	// ExitStringValue is called when exiting the StringValue production.
	ExitStringValue(c *StringValueContext)

	// ExitEqualsOp is called when exiting the EqualsOp production.
	ExitEqualsOp(c *EqualsOpContext)

	// ExitNotEqualsOp is called when exiting the NotEqualsOp production.
	ExitNotEqualsOp(c *NotEqualsOpContext)

	// ExitGreaterThanOp is called when exiting the GreaterThanOp production.
	ExitGreaterThanOp(c *GreaterThanOpContext)

	// ExitLessThanOp is called when exiting the LessThanOp production.
	ExitLessThanOp(c *LessThanOpContext)

	// ExitGreaterOrEqualOp is called when exiting the GreaterOrEqualOp production.
	ExitGreaterOrEqualOp(c *GreaterOrEqualOpContext)

	// ExitLessOrEqualOp is called when exiting the LessOrEqualOp production.
	ExitLessOrEqualOp(c *LessOrEqualOpContext)

	// ExitRegexOp is called when exiting the RegexOp production.
	ExitRegexOp(c *RegexOpContext)

	// ExitRegexNotOp is called when exiting the RegexNotOp production.
	ExitRegexNotOp(c *RegexNotOpContext)

	// ExitRegexValue is called when exiting the RegexValue production.
	ExitRegexValue(c *RegexValueContext)
}
