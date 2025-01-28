// Code generated from selectlang.g4 by ANTLR 4.13.2. DO NOT EDIT.

package selectlang // selectlang
import (
	"fmt"
	"strconv"
	"sync"

	"github.com/antlr4-go/antlr/v4"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}

type selectlangParser struct {
	*antlr.BaseParser
}

var SelectlangParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func selectlangParserInit() {
	staticData := &SelectlangParserStaticData
	staticData.LiteralNames = []string{
		"", "'report'", "'desired'", "'delete'", "'all'", "','", "':'", "'add'",
		"'remove'", "'update'", "'acknowledge'", "'no-change'", "'value'", "'WHERE'",
		"'('", "')'", "'=='", "'!='", "'>'", "'<'", "'>='", "'<='", "'before'",
		"'after'", "'id:'", "'AND'", "'OR'", "'NOT'", "'name:'", "'operation:'",
	}
	staticData.SymbolicNames = []string{
		"", "", "", "", "", "", "", "", "", "", "", "", "", "WHERE", "LPAREN",
		"RPAREN", "EQ", "NE", "GT", "LT", "GE", "LE", "BEFORE", "AFTER", "ID",
		"AND", "OR", "NOT", "NAME", "OPERATION", "NUMBER", "STRING", "TIME",
		"REGEX", "WS",
	}
	staticData.RuleNames = []string{
		"filter", "expression", "primaryExpr", "idExpr", "nameExpr", "operationExpr",
		"operations", "loggerExpr", "mapVarExpr", "loggerOp", "loggerConstraints",
		"valueComparison", "valueCondition", "valueFactor", "compareOp", "constantOrRegex",
		"regexOrString", "regex",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 34, 147, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7,
		10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15, 7, 15,
		2, 16, 7, 16, 2, 17, 7, 17, 1, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 3, 1, 47, 8, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 5, 1,
		55, 8, 1, 10, 1, 12, 1, 58, 9, 1, 1, 2, 1, 2, 1, 2, 3, 2, 63, 8, 2, 1,
		3, 1, 3, 1, 3, 1, 4, 1, 4, 1, 4, 1, 5, 1, 5, 1, 5, 1, 6, 1, 6, 1, 6, 5,
		6, 77, 8, 6, 10, 6, 12, 6, 80, 9, 6, 1, 7, 1, 7, 1, 7, 3, 7, 85, 8, 7,
		1, 7, 3, 7, 88, 8, 7, 1, 7, 3, 7, 91, 8, 7, 1, 8, 1, 8, 1, 8, 1, 9, 1,
		9, 1, 9, 5, 9, 99, 8, 9, 10, 9, 12, 9, 102, 9, 9, 1, 10, 1, 10, 1, 10,
		1, 10, 1, 10, 1, 11, 1, 11, 1, 11, 5, 11, 112, 8, 11, 10, 11, 12, 11, 115,
		9, 11, 1, 12, 1, 12, 1, 12, 5, 12, 120, 8, 12, 10, 12, 12, 12, 123, 9,
		12, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13, 3, 13, 133,
		8, 13, 1, 14, 1, 14, 1, 15, 1, 15, 1, 15, 1, 15, 3, 15, 141, 8, 15, 1,
		16, 1, 16, 1, 17, 1, 17, 1, 17, 0, 1, 2, 18, 0, 2, 4, 6, 8, 10, 12, 14,
		16, 18, 20, 22, 24, 26, 28, 30, 32, 34, 0, 5, 1, 0, 1, 4, 1, 0, 1, 3, 2,
		0, 4, 4, 7, 11, 1, 0, 16, 23, 2, 0, 31, 31, 33, 33, 145, 0, 36, 1, 0, 0,
		0, 2, 46, 1, 0, 0, 0, 4, 62, 1, 0, 0, 0, 6, 64, 1, 0, 0, 0, 8, 67, 1, 0,
		0, 0, 10, 70, 1, 0, 0, 0, 12, 73, 1, 0, 0, 0, 14, 81, 1, 0, 0, 0, 16, 92,
		1, 0, 0, 0, 18, 95, 1, 0, 0, 0, 20, 103, 1, 0, 0, 0, 22, 108, 1, 0, 0,
		0, 24, 116, 1, 0, 0, 0, 26, 132, 1, 0, 0, 0, 28, 134, 1, 0, 0, 0, 30, 140,
		1, 0, 0, 0, 32, 142, 1, 0, 0, 0, 34, 144, 1, 0, 0, 0, 36, 37, 3, 2, 1,
		0, 37, 38, 5, 0, 0, 1, 38, 1, 1, 0, 0, 0, 39, 40, 6, 1, -1, 0, 40, 47,
		3, 4, 2, 0, 41, 47, 3, 14, 7, 0, 42, 43, 5, 14, 0, 0, 43, 44, 3, 2, 1,
		0, 44, 45, 5, 15, 0, 0, 45, 47, 1, 0, 0, 0, 46, 39, 1, 0, 0, 0, 46, 41,
		1, 0, 0, 0, 46, 42, 1, 0, 0, 0, 47, 56, 1, 0, 0, 0, 48, 49, 10, 2, 0, 0,
		49, 50, 5, 25, 0, 0, 50, 55, 3, 2, 1, 3, 51, 52, 10, 1, 0, 0, 52, 53, 5,
		26, 0, 0, 53, 55, 3, 2, 1, 2, 54, 48, 1, 0, 0, 0, 54, 51, 1, 0, 0, 0, 55,
		58, 1, 0, 0, 0, 56, 54, 1, 0, 0, 0, 56, 57, 1, 0, 0, 0, 57, 3, 1, 0, 0,
		0, 58, 56, 1, 0, 0, 0, 59, 63, 3, 6, 3, 0, 60, 63, 3, 8, 4, 0, 61, 63,
		3, 10, 5, 0, 62, 59, 1, 0, 0, 0, 62, 60, 1, 0, 0, 0, 62, 61, 1, 0, 0, 0,
		63, 5, 1, 0, 0, 0, 64, 65, 5, 24, 0, 0, 65, 66, 3, 32, 16, 0, 66, 7, 1,
		0, 0, 0, 67, 68, 5, 28, 0, 0, 68, 69, 3, 32, 16, 0, 69, 9, 1, 0, 0, 0,
		70, 71, 5, 29, 0, 0, 71, 72, 3, 12, 6, 0, 72, 11, 1, 0, 0, 0, 73, 78, 7,
		0, 0, 0, 74, 75, 5, 5, 0, 0, 75, 77, 7, 1, 0, 0, 76, 74, 1, 0, 0, 0, 77,
		80, 1, 0, 0, 0, 78, 76, 1, 0, 0, 0, 78, 79, 1, 0, 0, 0, 79, 13, 1, 0, 0,
		0, 80, 78, 1, 0, 0, 0, 81, 84, 3, 18, 9, 0, 82, 83, 5, 6, 0, 0, 83, 85,
		3, 34, 17, 0, 84, 82, 1, 0, 0, 0, 84, 85, 1, 0, 0, 0, 85, 87, 1, 0, 0,
		0, 86, 88, 3, 16, 8, 0, 87, 86, 1, 0, 0, 0, 87, 88, 1, 0, 0, 0, 88, 90,
		1, 0, 0, 0, 89, 91, 3, 20, 10, 0, 90, 89, 1, 0, 0, 0, 90, 91, 1, 0, 0,
		0, 91, 15, 1, 0, 0, 0, 92, 93, 5, 16, 0, 0, 93, 94, 5, 31, 0, 0, 94, 17,
		1, 0, 0, 0, 95, 100, 7, 2, 0, 0, 96, 97, 5, 5, 0, 0, 97, 99, 7, 2, 0, 0,
		98, 96, 1, 0, 0, 0, 99, 102, 1, 0, 0, 0, 100, 98, 1, 0, 0, 0, 100, 101,
		1, 0, 0, 0, 101, 19, 1, 0, 0, 0, 102, 100, 1, 0, 0, 0, 103, 104, 5, 13,
		0, 0, 104, 105, 5, 14, 0, 0, 105, 106, 3, 22, 11, 0, 106, 107, 5, 15, 0,
		0, 107, 21, 1, 0, 0, 0, 108, 113, 3, 24, 12, 0, 109, 110, 5, 26, 0, 0,
		110, 112, 3, 24, 12, 0, 111, 109, 1, 0, 0, 0, 112, 115, 1, 0, 0, 0, 113,
		111, 1, 0, 0, 0, 113, 114, 1, 0, 0, 0, 114, 23, 1, 0, 0, 0, 115, 113, 1,
		0, 0, 0, 116, 121, 3, 26, 13, 0, 117, 118, 5, 25, 0, 0, 118, 120, 3, 26,
		13, 0, 119, 117, 1, 0, 0, 0, 120, 123, 1, 0, 0, 0, 121, 119, 1, 0, 0, 0,
		121, 122, 1, 0, 0, 0, 122, 25, 1, 0, 0, 0, 123, 121, 1, 0, 0, 0, 124, 125,
		5, 12, 0, 0, 125, 126, 3, 28, 14, 0, 126, 127, 3, 30, 15, 0, 127, 133,
		1, 0, 0, 0, 128, 129, 5, 14, 0, 0, 129, 130, 3, 22, 11, 0, 130, 131, 5,
		15, 0, 0, 131, 133, 1, 0, 0, 0, 132, 124, 1, 0, 0, 0, 132, 128, 1, 0, 0,
		0, 133, 27, 1, 0, 0, 0, 134, 135, 7, 3, 0, 0, 135, 29, 1, 0, 0, 0, 136,
		141, 5, 30, 0, 0, 137, 141, 5, 31, 0, 0, 138, 141, 5, 32, 0, 0, 139, 141,
		3, 34, 17, 0, 140, 136, 1, 0, 0, 0, 140, 137, 1, 0, 0, 0, 140, 138, 1,
		0, 0, 0, 140, 139, 1, 0, 0, 0, 141, 31, 1, 0, 0, 0, 142, 143, 7, 4, 0,
		0, 143, 33, 1, 0, 0, 0, 144, 145, 5, 33, 0, 0, 145, 35, 1, 0, 0, 0, 13,
		46, 54, 56, 62, 78, 84, 87, 90, 100, 113, 121, 132, 140,
	}
	deserializer := antlr.NewATNDeserializer(nil)
	staticData.atn = deserializer.Deserialize(staticData.serializedATN)
	atn := staticData.atn
	staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
	decisionToDFA := staticData.decisionToDFA
	for index, state := range atn.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(state, index)
	}
}

// selectlangParserInit initializes any static state used to implement selectlangParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewselectlangParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func SelectlangParserInit() {
	staticData := &SelectlangParserStaticData
	staticData.once.Do(selectlangParserInit)
}

// NewselectlangParser produces a new parser instance for the optional input antlr.TokenStream.
func NewselectlangParser(input antlr.TokenStream) *selectlangParser {
	SelectlangParserInit()
	this := new(selectlangParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &SelectlangParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	this.RuleNames = staticData.RuleNames
	this.LiteralNames = staticData.LiteralNames
	this.SymbolicNames = staticData.SymbolicNames
	this.GrammarFileName = "selectlang.g4"

	return this
}

// selectlangParser tokens.
const (
	selectlangParserEOF       = antlr.TokenEOF
	selectlangParserT__0      = 1
	selectlangParserT__1      = 2
	selectlangParserT__2      = 3
	selectlangParserT__3      = 4
	selectlangParserT__4      = 5
	selectlangParserT__5      = 6
	selectlangParserT__6      = 7
	selectlangParserT__7      = 8
	selectlangParserT__8      = 9
	selectlangParserT__9      = 10
	selectlangParserT__10     = 11
	selectlangParserT__11     = 12
	selectlangParserWHERE     = 13
	selectlangParserLPAREN    = 14
	selectlangParserRPAREN    = 15
	selectlangParserEQ        = 16
	selectlangParserNE        = 17
	selectlangParserGT        = 18
	selectlangParserLT        = 19
	selectlangParserGE        = 20
	selectlangParserLE        = 21
	selectlangParserBEFORE    = 22
	selectlangParserAFTER     = 23
	selectlangParserID        = 24
	selectlangParserAND       = 25
	selectlangParserOR        = 26
	selectlangParserNOT       = 27
	selectlangParserNAME      = 28
	selectlangParserOPERATION = 29
	selectlangParserNUMBER    = 30
	selectlangParserSTRING    = 31
	selectlangParserTIME      = 32
	selectlangParserREGEX     = 33
	selectlangParserWS        = 34
)

// selectlangParser rules.
const (
	selectlangParserRULE_filter            = 0
	selectlangParserRULE_expression        = 1
	selectlangParserRULE_primaryExpr       = 2
	selectlangParserRULE_idExpr            = 3
	selectlangParserRULE_nameExpr          = 4
	selectlangParserRULE_operationExpr     = 5
	selectlangParserRULE_operations        = 6
	selectlangParserRULE_loggerExpr        = 7
	selectlangParserRULE_mapVarExpr        = 8
	selectlangParserRULE_loggerOp          = 9
	selectlangParserRULE_loggerConstraints = 10
	selectlangParserRULE_valueComparison   = 11
	selectlangParserRULE_valueCondition    = 12
	selectlangParserRULE_valueFactor       = 13
	selectlangParserRULE_compareOp         = 14
	selectlangParserRULE_constantOrRegex   = 15
	selectlangParserRULE_regexOrString     = 16
	selectlangParserRULE_regex             = 17
)

// IFilterContext is an interface to support dynamic dispatch.
type IFilterContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Expression() IExpressionContext
	EOF() antlr.TerminalNode

	// IsFilterContext differentiates from other interfaces.
	IsFilterContext()
}

type FilterContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFilterContext() *FilterContext {
	var p = new(FilterContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_filter
	return p
}

func InitEmptyFilterContext(p *FilterContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_filter
}

func (*FilterContext) IsFilterContext() {}

func NewFilterContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FilterContext {
	var p = new(FilterContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = selectlangParserRULE_filter

	return p
}

func (s *FilterContext) GetParser() antlr.Parser { return s.parser }

func (s *FilterContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *FilterContext) EOF() antlr.TerminalNode {
	return s.GetToken(selectlangParserEOF, 0)
}

func (s *FilterContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FilterContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FilterContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.EnterFilter(s)
	}
}

func (s *FilterContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.ExitFilter(s)
	}
}

func (p *selectlangParser) Filter() (localctx IFilterContext) {
	localctx = NewFilterContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, selectlangParserRULE_filter)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(36)
		p.expression(0)
	}
	{
		p.SetState(37)
		p.Match(selectlangParserEOF)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IExpressionContext is an interface to support dynamic dispatch.
type IExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	PrimaryExpr() IPrimaryExprContext
	LoggerExpr() ILoggerExprContext
	LPAREN() antlr.TerminalNode
	AllExpression() []IExpressionContext
	Expression(i int) IExpressionContext
	RPAREN() antlr.TerminalNode
	AND() antlr.TerminalNode
	OR() antlr.TerminalNode

	// IsExpressionContext differentiates from other interfaces.
	IsExpressionContext()
}

type ExpressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExpressionContext() *ExpressionContext {
	var p = new(ExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_expression
	return p
}

func InitEmptyExpressionContext(p *ExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_expression
}

func (*ExpressionContext) IsExpressionContext() {}

func NewExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExpressionContext {
	var p = new(ExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = selectlangParserRULE_expression

	return p
}

func (s *ExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *ExpressionContext) PrimaryExpr() IPrimaryExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPrimaryExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPrimaryExprContext)
}

func (s *ExpressionContext) LoggerExpr() ILoggerExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILoggerExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILoggerExprContext)
}

func (s *ExpressionContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(selectlangParserLPAREN, 0)
}

func (s *ExpressionContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *ExpressionContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ExpressionContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(selectlangParserRPAREN, 0)
}

func (s *ExpressionContext) AND() antlr.TerminalNode {
	return s.GetToken(selectlangParserAND, 0)
}

func (s *ExpressionContext) OR() antlr.TerminalNode {
	return s.GetToken(selectlangParserOR, 0)
}

func (s *ExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.EnterExpression(s)
	}
}

func (s *ExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.ExitExpression(s)
	}
}

func (p *selectlangParser) Expression() (localctx IExpressionContext) {
	return p.expression(0)
}

func (p *selectlangParser) expression(_p int) (localctx IExpressionContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()

	_parentState := p.GetState()
	localctx = NewExpressionContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IExpressionContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 2
	p.EnterRecursionRule(localctx, 2, selectlangParserRULE_expression, _p)
	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(46)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case selectlangParserID, selectlangParserNAME, selectlangParserOPERATION:
		{
			p.SetState(40)
			p.PrimaryExpr()
		}

	case selectlangParserT__3, selectlangParserT__6, selectlangParserT__7, selectlangParserT__8, selectlangParserT__9, selectlangParserT__10:
		{
			p.SetState(41)
			p.LoggerExpr()
		}

	case selectlangParserLPAREN:
		{
			p.SetState(42)
			p.Match(selectlangParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(43)
			p.expression(0)
		}
		{
			p.SetState(44)
			p.Match(selectlangParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(56)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 2, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(54)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}

			switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 1, p.GetParserRuleContext()) {
			case 1:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, selectlangParserRULE_expression)
				p.SetState(48)

				if !(p.Precpred(p.GetParserRuleContext(), 2)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
					goto errorExit
				}
				{
					p.SetState(49)
					p.Match(selectlangParserAND)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(50)
					p.expression(3)
				}

			case 2:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, selectlangParserRULE_expression)
				p.SetState(51)

				if !(p.Precpred(p.GetParserRuleContext(), 1)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 1)", ""))
					goto errorExit
				}
				{
					p.SetState(52)
					p.Match(selectlangParserOR)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(53)
					p.expression(2)
				}

			case antlr.ATNInvalidAltNumber:
				goto errorExit
			}

		}
		p.SetState(58)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 2, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.UnrollRecursionContexts(_parentctx)
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IPrimaryExprContext is an interface to support dynamic dispatch.
type IPrimaryExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IdExpr() IIdExprContext
	NameExpr() INameExprContext
	OperationExpr() IOperationExprContext

	// IsPrimaryExprContext differentiates from other interfaces.
	IsPrimaryExprContext()
}

type PrimaryExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPrimaryExprContext() *PrimaryExprContext {
	var p = new(PrimaryExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_primaryExpr
	return p
}

func InitEmptyPrimaryExprContext(p *PrimaryExprContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_primaryExpr
}

func (*PrimaryExprContext) IsPrimaryExprContext() {}

func NewPrimaryExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PrimaryExprContext {
	var p = new(PrimaryExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = selectlangParserRULE_primaryExpr

	return p
}

func (s *PrimaryExprContext) GetParser() antlr.Parser { return s.parser }

func (s *PrimaryExprContext) IdExpr() IIdExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdExprContext)
}

func (s *PrimaryExprContext) NameExpr() INameExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(INameExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(INameExprContext)
}

func (s *PrimaryExprContext) OperationExpr() IOperationExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IOperationExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IOperationExprContext)
}

func (s *PrimaryExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PrimaryExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PrimaryExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.EnterPrimaryExpr(s)
	}
}

func (s *PrimaryExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.ExitPrimaryExpr(s)
	}
}

func (p *selectlangParser) PrimaryExpr() (localctx IPrimaryExprContext) {
	localctx = NewPrimaryExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, selectlangParserRULE_primaryExpr)
	p.SetState(62)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case selectlangParserID:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(59)
			p.IdExpr()
		}

	case selectlangParserNAME:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(60)
			p.NameExpr()
		}

	case selectlangParserOPERATION:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(61)
			p.OperationExpr()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IIdExprContext is an interface to support dynamic dispatch.
type IIdExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ID() antlr.TerminalNode
	RegexOrString() IRegexOrStringContext

	// IsIdExprContext differentiates from other interfaces.
	IsIdExprContext()
}

type IdExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIdExprContext() *IdExprContext {
	var p = new(IdExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_idExpr
	return p
}

func InitEmptyIdExprContext(p *IdExprContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_idExpr
}

func (*IdExprContext) IsIdExprContext() {}

func NewIdExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IdExprContext {
	var p = new(IdExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = selectlangParserRULE_idExpr

	return p
}

func (s *IdExprContext) GetParser() antlr.Parser { return s.parser }

func (s *IdExprContext) ID() antlr.TerminalNode {
	return s.GetToken(selectlangParserID, 0)
}

func (s *IdExprContext) RegexOrString() IRegexOrStringContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRegexOrStringContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRegexOrStringContext)
}

func (s *IdExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IdExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *IdExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.EnterIdExpr(s)
	}
}

func (s *IdExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.ExitIdExpr(s)
	}
}

func (p *selectlangParser) IdExpr() (localctx IIdExprContext) {
	localctx = NewIdExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, selectlangParserRULE_idExpr)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(64)
		p.Match(selectlangParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(65)
		p.RegexOrString()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// INameExprContext is an interface to support dynamic dispatch.
type INameExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	NAME() antlr.TerminalNode
	RegexOrString() IRegexOrStringContext

	// IsNameExprContext differentiates from other interfaces.
	IsNameExprContext()
}

type NameExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNameExprContext() *NameExprContext {
	var p = new(NameExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_nameExpr
	return p
}

func InitEmptyNameExprContext(p *NameExprContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_nameExpr
}

func (*NameExprContext) IsNameExprContext() {}

func NewNameExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NameExprContext {
	var p = new(NameExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = selectlangParserRULE_nameExpr

	return p
}

func (s *NameExprContext) GetParser() antlr.Parser { return s.parser }

func (s *NameExprContext) NAME() antlr.TerminalNode {
	return s.GetToken(selectlangParserNAME, 0)
}

func (s *NameExprContext) RegexOrString() IRegexOrStringContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRegexOrStringContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRegexOrStringContext)
}

func (s *NameExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NameExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *NameExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.EnterNameExpr(s)
	}
}

func (s *NameExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.ExitNameExpr(s)
	}
}

func (p *selectlangParser) NameExpr() (localctx INameExprContext) {
	localctx = NewNameExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, selectlangParserRULE_nameExpr)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(67)
		p.Match(selectlangParserNAME)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(68)
		p.RegexOrString()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IOperationExprContext is an interface to support dynamic dispatch.
type IOperationExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	OPERATION() antlr.TerminalNode
	Operations() IOperationsContext

	// IsOperationExprContext differentiates from other interfaces.
	IsOperationExprContext()
}

type OperationExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyOperationExprContext() *OperationExprContext {
	var p = new(OperationExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_operationExpr
	return p
}

func InitEmptyOperationExprContext(p *OperationExprContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_operationExpr
}

func (*OperationExprContext) IsOperationExprContext() {}

func NewOperationExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *OperationExprContext {
	var p = new(OperationExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = selectlangParserRULE_operationExpr

	return p
}

func (s *OperationExprContext) GetParser() antlr.Parser { return s.parser }

func (s *OperationExprContext) OPERATION() antlr.TerminalNode {
	return s.GetToken(selectlangParserOPERATION, 0)
}

func (s *OperationExprContext) Operations() IOperationsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IOperationsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IOperationsContext)
}

func (s *OperationExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *OperationExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *OperationExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.EnterOperationExpr(s)
	}
}

func (s *OperationExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.ExitOperationExpr(s)
	}
}

func (p *selectlangParser) OperationExpr() (localctx IOperationExprContext) {
	localctx = NewOperationExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, selectlangParserRULE_operationExpr)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(70)
		p.Match(selectlangParserOPERATION)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(71)
		p.Operations()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IOperationsContext is an interface to support dynamic dispatch.
type IOperationsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsOperationsContext differentiates from other interfaces.
	IsOperationsContext()
}

type OperationsContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyOperationsContext() *OperationsContext {
	var p = new(OperationsContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_operations
	return p
}

func InitEmptyOperationsContext(p *OperationsContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_operations
}

func (*OperationsContext) IsOperationsContext() {}

func NewOperationsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *OperationsContext {
	var p = new(OperationsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = selectlangParserRULE_operations

	return p
}

func (s *OperationsContext) GetParser() antlr.Parser { return s.parser }
func (s *OperationsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *OperationsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *OperationsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.EnterOperations(s)
	}
}

func (s *OperationsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.ExitOperations(s)
	}
}

func (p *selectlangParser) Operations() (localctx IOperationsContext) {
	localctx = NewOperationsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, selectlangParserRULE_operations)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(73)
		_la = p.GetTokenStream().LA(1)

		if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&30) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}
	p.SetState(78)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 4, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(74)
				p.Match(selectlangParserT__4)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(75)
				_la = p.GetTokenStream().LA(1)

				if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&14) != 0) {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}

		}
		p.SetState(80)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 4, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ILoggerExprContext is an interface to support dynamic dispatch.
type ILoggerExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LoggerOp() ILoggerOpContext
	Regex() IRegexContext
	MapVarExpr() IMapVarExprContext
	LoggerConstraints() ILoggerConstraintsContext

	// IsLoggerExprContext differentiates from other interfaces.
	IsLoggerExprContext()
}

type LoggerExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLoggerExprContext() *LoggerExprContext {
	var p = new(LoggerExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_loggerExpr
	return p
}

func InitEmptyLoggerExprContext(p *LoggerExprContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_loggerExpr
}

func (*LoggerExprContext) IsLoggerExprContext() {}

func NewLoggerExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LoggerExprContext {
	var p = new(LoggerExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = selectlangParserRULE_loggerExpr

	return p
}

func (s *LoggerExprContext) GetParser() antlr.Parser { return s.parser }

func (s *LoggerExprContext) LoggerOp() ILoggerOpContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILoggerOpContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILoggerOpContext)
}

func (s *LoggerExprContext) Regex() IRegexContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRegexContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRegexContext)
}

func (s *LoggerExprContext) MapVarExpr() IMapVarExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMapVarExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMapVarExprContext)
}

func (s *LoggerExprContext) LoggerConstraints() ILoggerConstraintsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILoggerConstraintsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILoggerConstraintsContext)
}

func (s *LoggerExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LoggerExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LoggerExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.EnterLoggerExpr(s)
	}
}

func (s *LoggerExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.ExitLoggerExpr(s)
	}
}

func (p *selectlangParser) LoggerExpr() (localctx ILoggerExprContext) {
	localctx = NewLoggerExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, selectlangParserRULE_loggerExpr)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(81)
		p.LoggerOp()
	}
	p.SetState(84)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 5, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(82)
			p.Match(selectlangParserT__5)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(83)
			p.Regex()
		}

	} else if p.HasError() { // JIM
		goto errorExit
	}
	p.SetState(87)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 6, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(86)
			p.MapVarExpr()
		}

	} else if p.HasError() { // JIM
		goto errorExit
	}
	p.SetState(90)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 7, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(89)
			p.LoggerConstraints()
		}

	} else if p.HasError() { // JIM
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IMapVarExprContext is an interface to support dynamic dispatch.
type IMapVarExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	EQ() antlr.TerminalNode
	STRING() antlr.TerminalNode

	// IsMapVarExprContext differentiates from other interfaces.
	IsMapVarExprContext()
}

type MapVarExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMapVarExprContext() *MapVarExprContext {
	var p = new(MapVarExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_mapVarExpr
	return p
}

func InitEmptyMapVarExprContext(p *MapVarExprContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_mapVarExpr
}

func (*MapVarExprContext) IsMapVarExprContext() {}

func NewMapVarExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MapVarExprContext {
	var p = new(MapVarExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = selectlangParserRULE_mapVarExpr

	return p
}

func (s *MapVarExprContext) GetParser() antlr.Parser { return s.parser }

func (s *MapVarExprContext) EQ() antlr.TerminalNode {
	return s.GetToken(selectlangParserEQ, 0)
}

func (s *MapVarExprContext) STRING() antlr.TerminalNode {
	return s.GetToken(selectlangParserSTRING, 0)
}

func (s *MapVarExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MapVarExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *MapVarExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.EnterMapVarExpr(s)
	}
}

func (s *MapVarExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.ExitMapVarExpr(s)
	}
}

func (p *selectlangParser) MapVarExpr() (localctx IMapVarExprContext) {
	localctx = NewMapVarExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, selectlangParserRULE_mapVarExpr)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(92)
		p.Match(selectlangParserEQ)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(93)
		p.Match(selectlangParserSTRING)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ILoggerOpContext is an interface to support dynamic dispatch.
type ILoggerOpContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsLoggerOpContext differentiates from other interfaces.
	IsLoggerOpContext()
}

type LoggerOpContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLoggerOpContext() *LoggerOpContext {
	var p = new(LoggerOpContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_loggerOp
	return p
}

func InitEmptyLoggerOpContext(p *LoggerOpContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_loggerOp
}

func (*LoggerOpContext) IsLoggerOpContext() {}

func NewLoggerOpContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LoggerOpContext {
	var p = new(LoggerOpContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = selectlangParserRULE_loggerOp

	return p
}

func (s *LoggerOpContext) GetParser() antlr.Parser { return s.parser }
func (s *LoggerOpContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LoggerOpContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LoggerOpContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.EnterLoggerOp(s)
	}
}

func (s *LoggerOpContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.ExitLoggerOp(s)
	}
}

func (p *selectlangParser) LoggerOp() (localctx ILoggerOpContext) {
	localctx = NewLoggerOpContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, selectlangParserRULE_loggerOp)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(95)
		_la = p.GetTokenStream().LA(1)

		if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&3984) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}
	p.SetState(100)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 8, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(96)
				p.Match(selectlangParserT__4)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(97)
				_la = p.GetTokenStream().LA(1)

				if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&3984) != 0) {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}

		}
		p.SetState(102)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 8, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ILoggerConstraintsContext is an interface to support dynamic dispatch.
type ILoggerConstraintsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	WHERE() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	ValueComparison() IValueComparisonContext
	RPAREN() antlr.TerminalNode

	// IsLoggerConstraintsContext differentiates from other interfaces.
	IsLoggerConstraintsContext()
}

type LoggerConstraintsContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLoggerConstraintsContext() *LoggerConstraintsContext {
	var p = new(LoggerConstraintsContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_loggerConstraints
	return p
}

func InitEmptyLoggerConstraintsContext(p *LoggerConstraintsContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_loggerConstraints
}

func (*LoggerConstraintsContext) IsLoggerConstraintsContext() {}

func NewLoggerConstraintsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LoggerConstraintsContext {
	var p = new(LoggerConstraintsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = selectlangParserRULE_loggerConstraints

	return p
}

func (s *LoggerConstraintsContext) GetParser() antlr.Parser { return s.parser }

func (s *LoggerConstraintsContext) WHERE() antlr.TerminalNode {
	return s.GetToken(selectlangParserWHERE, 0)
}

func (s *LoggerConstraintsContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(selectlangParserLPAREN, 0)
}

func (s *LoggerConstraintsContext) ValueComparison() IValueComparisonContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IValueComparisonContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IValueComparisonContext)
}

func (s *LoggerConstraintsContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(selectlangParserRPAREN, 0)
}

func (s *LoggerConstraintsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LoggerConstraintsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LoggerConstraintsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.EnterLoggerConstraints(s)
	}
}

func (s *LoggerConstraintsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.ExitLoggerConstraints(s)
	}
}

func (p *selectlangParser) LoggerConstraints() (localctx ILoggerConstraintsContext) {
	localctx = NewLoggerConstraintsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, selectlangParserRULE_loggerConstraints)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(103)
		p.Match(selectlangParserWHERE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(104)
		p.Match(selectlangParserLPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(105)
		p.ValueComparison()
	}
	{
		p.SetState(106)
		p.Match(selectlangParserRPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IValueComparisonContext is an interface to support dynamic dispatch.
type IValueComparisonContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllValueCondition() []IValueConditionContext
	ValueCondition(i int) IValueConditionContext
	AllOR() []antlr.TerminalNode
	OR(i int) antlr.TerminalNode

	// IsValueComparisonContext differentiates from other interfaces.
	IsValueComparisonContext()
}

type ValueComparisonContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyValueComparisonContext() *ValueComparisonContext {
	var p = new(ValueComparisonContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_valueComparison
	return p
}

func InitEmptyValueComparisonContext(p *ValueComparisonContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_valueComparison
}

func (*ValueComparisonContext) IsValueComparisonContext() {}

func NewValueComparisonContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ValueComparisonContext {
	var p = new(ValueComparisonContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = selectlangParserRULE_valueComparison

	return p
}

func (s *ValueComparisonContext) GetParser() antlr.Parser { return s.parser }

func (s *ValueComparisonContext) AllValueCondition() []IValueConditionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IValueConditionContext); ok {
			len++
		}
	}

	tst := make([]IValueConditionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IValueConditionContext); ok {
			tst[i] = t.(IValueConditionContext)
			i++
		}
	}

	return tst
}

func (s *ValueComparisonContext) ValueCondition(i int) IValueConditionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IValueConditionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IValueConditionContext)
}

func (s *ValueComparisonContext) AllOR() []antlr.TerminalNode {
	return s.GetTokens(selectlangParserOR)
}

func (s *ValueComparisonContext) OR(i int) antlr.TerminalNode {
	return s.GetToken(selectlangParserOR, i)
}

func (s *ValueComparisonContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ValueComparisonContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ValueComparisonContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.EnterValueComparison(s)
	}
}

func (s *ValueComparisonContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.ExitValueComparison(s)
	}
}

func (p *selectlangParser) ValueComparison() (localctx IValueComparisonContext) {
	localctx = NewValueComparisonContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, selectlangParserRULE_valueComparison)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(108)
		p.ValueCondition()
	}
	p.SetState(113)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == selectlangParserOR {
		{
			p.SetState(109)
			p.Match(selectlangParserOR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(110)
			p.ValueCondition()
		}

		p.SetState(115)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IValueConditionContext is an interface to support dynamic dispatch.
type IValueConditionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllValueFactor() []IValueFactorContext
	ValueFactor(i int) IValueFactorContext
	AllAND() []antlr.TerminalNode
	AND(i int) antlr.TerminalNode

	// IsValueConditionContext differentiates from other interfaces.
	IsValueConditionContext()
}

type ValueConditionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyValueConditionContext() *ValueConditionContext {
	var p = new(ValueConditionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_valueCondition
	return p
}

func InitEmptyValueConditionContext(p *ValueConditionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_valueCondition
}

func (*ValueConditionContext) IsValueConditionContext() {}

func NewValueConditionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ValueConditionContext {
	var p = new(ValueConditionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = selectlangParserRULE_valueCondition

	return p
}

func (s *ValueConditionContext) GetParser() antlr.Parser { return s.parser }

func (s *ValueConditionContext) AllValueFactor() []IValueFactorContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IValueFactorContext); ok {
			len++
		}
	}

	tst := make([]IValueFactorContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IValueFactorContext); ok {
			tst[i] = t.(IValueFactorContext)
			i++
		}
	}

	return tst
}

func (s *ValueConditionContext) ValueFactor(i int) IValueFactorContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IValueFactorContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IValueFactorContext)
}

func (s *ValueConditionContext) AllAND() []antlr.TerminalNode {
	return s.GetTokens(selectlangParserAND)
}

func (s *ValueConditionContext) AND(i int) antlr.TerminalNode {
	return s.GetToken(selectlangParserAND, i)
}

func (s *ValueConditionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ValueConditionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ValueConditionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.EnterValueCondition(s)
	}
}

func (s *ValueConditionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.ExitValueCondition(s)
	}
}

func (p *selectlangParser) ValueCondition() (localctx IValueConditionContext) {
	localctx = NewValueConditionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, selectlangParserRULE_valueCondition)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(116)
		p.ValueFactor()
	}
	p.SetState(121)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == selectlangParserAND {
		{
			p.SetState(117)
			p.Match(selectlangParserAND)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(118)
			p.ValueFactor()
		}

		p.SetState(123)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IValueFactorContext is an interface to support dynamic dispatch.
type IValueFactorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CompareOp() ICompareOpContext
	ConstantOrRegex() IConstantOrRegexContext
	LPAREN() antlr.TerminalNode
	ValueComparison() IValueComparisonContext
	RPAREN() antlr.TerminalNode

	// IsValueFactorContext differentiates from other interfaces.
	IsValueFactorContext()
}

type ValueFactorContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyValueFactorContext() *ValueFactorContext {
	var p = new(ValueFactorContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_valueFactor
	return p
}

func InitEmptyValueFactorContext(p *ValueFactorContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_valueFactor
}

func (*ValueFactorContext) IsValueFactorContext() {}

func NewValueFactorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ValueFactorContext {
	var p = new(ValueFactorContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = selectlangParserRULE_valueFactor

	return p
}

func (s *ValueFactorContext) GetParser() antlr.Parser { return s.parser }

func (s *ValueFactorContext) CompareOp() ICompareOpContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICompareOpContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICompareOpContext)
}

func (s *ValueFactorContext) ConstantOrRegex() IConstantOrRegexContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IConstantOrRegexContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IConstantOrRegexContext)
}

func (s *ValueFactorContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(selectlangParserLPAREN, 0)
}

func (s *ValueFactorContext) ValueComparison() IValueComparisonContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IValueComparisonContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IValueComparisonContext)
}

func (s *ValueFactorContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(selectlangParserRPAREN, 0)
}

func (s *ValueFactorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ValueFactorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ValueFactorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.EnterValueFactor(s)
	}
}

func (s *ValueFactorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.ExitValueFactor(s)
	}
}

func (p *selectlangParser) ValueFactor() (localctx IValueFactorContext) {
	localctx = NewValueFactorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, selectlangParserRULE_valueFactor)
	p.SetState(132)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case selectlangParserT__11:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(124)
			p.Match(selectlangParserT__11)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(125)
			p.CompareOp()
		}
		{
			p.SetState(126)
			p.ConstantOrRegex()
		}

	case selectlangParserLPAREN:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(128)
			p.Match(selectlangParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(129)
			p.ValueComparison()
		}
		{
			p.SetState(130)
			p.Match(selectlangParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ICompareOpContext is an interface to support dynamic dispatch.
type ICompareOpContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	EQ() antlr.TerminalNode
	NE() antlr.TerminalNode
	GT() antlr.TerminalNode
	LT() antlr.TerminalNode
	GE() antlr.TerminalNode
	LE() antlr.TerminalNode
	BEFORE() antlr.TerminalNode
	AFTER() antlr.TerminalNode

	// IsCompareOpContext differentiates from other interfaces.
	IsCompareOpContext()
}

type CompareOpContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCompareOpContext() *CompareOpContext {
	var p = new(CompareOpContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_compareOp
	return p
}

func InitEmptyCompareOpContext(p *CompareOpContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_compareOp
}

func (*CompareOpContext) IsCompareOpContext() {}

func NewCompareOpContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CompareOpContext {
	var p = new(CompareOpContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = selectlangParserRULE_compareOp

	return p
}

func (s *CompareOpContext) GetParser() antlr.Parser { return s.parser }

func (s *CompareOpContext) EQ() antlr.TerminalNode {
	return s.GetToken(selectlangParserEQ, 0)
}

func (s *CompareOpContext) NE() antlr.TerminalNode {
	return s.GetToken(selectlangParserNE, 0)
}

func (s *CompareOpContext) GT() antlr.TerminalNode {
	return s.GetToken(selectlangParserGT, 0)
}

func (s *CompareOpContext) LT() antlr.TerminalNode {
	return s.GetToken(selectlangParserLT, 0)
}

func (s *CompareOpContext) GE() antlr.TerminalNode {
	return s.GetToken(selectlangParserGE, 0)
}

func (s *CompareOpContext) LE() antlr.TerminalNode {
	return s.GetToken(selectlangParserLE, 0)
}

func (s *CompareOpContext) BEFORE() antlr.TerminalNode {
	return s.GetToken(selectlangParserBEFORE, 0)
}

func (s *CompareOpContext) AFTER() antlr.TerminalNode {
	return s.GetToken(selectlangParserAFTER, 0)
}

func (s *CompareOpContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CompareOpContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CompareOpContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.EnterCompareOp(s)
	}
}

func (s *CompareOpContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.ExitCompareOp(s)
	}
}

func (p *selectlangParser) CompareOp() (localctx ICompareOpContext) {
	localctx = NewCompareOpContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, selectlangParserRULE_compareOp)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(134)
		_la = p.GetTokenStream().LA(1)

		if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&16711680) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IConstantOrRegexContext is an interface to support dynamic dispatch.
type IConstantOrRegexContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsConstantOrRegexContext differentiates from other interfaces.
	IsConstantOrRegexContext()
}

type ConstantOrRegexContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyConstantOrRegexContext() *ConstantOrRegexContext {
	var p = new(ConstantOrRegexContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_constantOrRegex
	return p
}

func InitEmptyConstantOrRegexContext(p *ConstantOrRegexContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_constantOrRegex
}

func (*ConstantOrRegexContext) IsConstantOrRegexContext() {}

func NewConstantOrRegexContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ConstantOrRegexContext {
	var p = new(ConstantOrRegexContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = selectlangParserRULE_constantOrRegex

	return p
}

func (s *ConstantOrRegexContext) GetParser() antlr.Parser { return s.parser }

func (s *ConstantOrRegexContext) CopyAll(ctx *ConstantOrRegexContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *ConstantOrRegexContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ConstantOrRegexContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type NumericLiteralContext struct {
	ConstantOrRegexContext
}

func NewNumericLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *NumericLiteralContext {
	var p = new(NumericLiteralContext)

	InitEmptyConstantOrRegexContext(&p.ConstantOrRegexContext)
	p.parser = parser
	p.CopyAll(ctx.(*ConstantOrRegexContext))

	return p
}

func (s *NumericLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NumericLiteralContext) NUMBER() antlr.TerminalNode {
	return s.GetToken(selectlangParserNUMBER, 0)
}

func (s *NumericLiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.EnterNumericLiteral(s)
	}
}

func (s *NumericLiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.ExitNumericLiteral(s)
	}
}

type StringLiteralContext struct {
	ConstantOrRegexContext
}

func NewStringLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *StringLiteralContext {
	var p = new(StringLiteralContext)

	InitEmptyConstantOrRegexContext(&p.ConstantOrRegexContext)
	p.parser = parser
	p.CopyAll(ctx.(*ConstantOrRegexContext))

	return p
}

func (s *StringLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StringLiteralContext) STRING() antlr.TerminalNode {
	return s.GetToken(selectlangParserSTRING, 0)
}

func (s *StringLiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.EnterStringLiteral(s)
	}
}

func (s *StringLiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.ExitStringLiteral(s)
	}
}

type RegexLiteralContext struct {
	ConstantOrRegexContext
}

func NewRegexLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *RegexLiteralContext {
	var p = new(RegexLiteralContext)

	InitEmptyConstantOrRegexContext(&p.ConstantOrRegexContext)
	p.parser = parser
	p.CopyAll(ctx.(*ConstantOrRegexContext))

	return p
}

func (s *RegexLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RegexLiteralContext) Regex() IRegexContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRegexContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRegexContext)
}

func (s *RegexLiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.EnterRegexLiteral(s)
	}
}

func (s *RegexLiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.ExitRegexLiteral(s)
	}
}

type TimeLiteralContext struct {
	ConstantOrRegexContext
}

func NewTimeLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *TimeLiteralContext {
	var p = new(TimeLiteralContext)

	InitEmptyConstantOrRegexContext(&p.ConstantOrRegexContext)
	p.parser = parser
	p.CopyAll(ctx.(*ConstantOrRegexContext))

	return p
}

func (s *TimeLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TimeLiteralContext) TIME() antlr.TerminalNode {
	return s.GetToken(selectlangParserTIME, 0)
}

func (s *TimeLiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.EnterTimeLiteral(s)
	}
}

func (s *TimeLiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.ExitTimeLiteral(s)
	}
}

func (p *selectlangParser) ConstantOrRegex() (localctx IConstantOrRegexContext) {
	localctx = NewConstantOrRegexContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, selectlangParserRULE_constantOrRegex)
	p.SetState(140)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case selectlangParserNUMBER:
		localctx = NewNumericLiteralContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(136)
			p.Match(selectlangParserNUMBER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case selectlangParserSTRING:
		localctx = NewStringLiteralContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(137)
			p.Match(selectlangParserSTRING)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case selectlangParserTIME:
		localctx = NewTimeLiteralContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(138)
			p.Match(selectlangParserTIME)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case selectlangParserREGEX:
		localctx = NewRegexLiteralContext(p, localctx)
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(139)
			p.Regex()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IRegexOrStringContext is an interface to support dynamic dispatch.
type IRegexOrStringContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	REGEX() antlr.TerminalNode
	STRING() antlr.TerminalNode

	// IsRegexOrStringContext differentiates from other interfaces.
	IsRegexOrStringContext()
}

type RegexOrStringContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRegexOrStringContext() *RegexOrStringContext {
	var p = new(RegexOrStringContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_regexOrString
	return p
}

func InitEmptyRegexOrStringContext(p *RegexOrStringContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_regexOrString
}

func (*RegexOrStringContext) IsRegexOrStringContext() {}

func NewRegexOrStringContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RegexOrStringContext {
	var p = new(RegexOrStringContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = selectlangParserRULE_regexOrString

	return p
}

func (s *RegexOrStringContext) GetParser() antlr.Parser { return s.parser }

func (s *RegexOrStringContext) REGEX() antlr.TerminalNode {
	return s.GetToken(selectlangParserREGEX, 0)
}

func (s *RegexOrStringContext) STRING() antlr.TerminalNode {
	return s.GetToken(selectlangParserSTRING, 0)
}

func (s *RegexOrStringContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RegexOrStringContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RegexOrStringContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.EnterRegexOrString(s)
	}
}

func (s *RegexOrStringContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.ExitRegexOrString(s)
	}
}

func (p *selectlangParser) RegexOrString() (localctx IRegexOrStringContext) {
	localctx = NewRegexOrStringContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 32, selectlangParserRULE_regexOrString)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(142)
		_la = p.GetTokenStream().LA(1)

		if !(_la == selectlangParserSTRING || _la == selectlangParserREGEX) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IRegexContext is an interface to support dynamic dispatch.
type IRegexContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	REGEX() antlr.TerminalNode

	// IsRegexContext differentiates from other interfaces.
	IsRegexContext()
}

type RegexContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRegexContext() *RegexContext {
	var p = new(RegexContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_regex
	return p
}

func InitEmptyRegexContext(p *RegexContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_regex
}

func (*RegexContext) IsRegexContext() {}

func NewRegexContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RegexContext {
	var p = new(RegexContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = selectlangParserRULE_regex

	return p
}

func (s *RegexContext) GetParser() antlr.Parser { return s.parser }

func (s *RegexContext) REGEX() antlr.TerminalNode {
	return s.GetToken(selectlangParserREGEX, 0)
}

func (s *RegexContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RegexContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RegexContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.EnterRegex(s)
	}
}

func (s *RegexContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.ExitRegex(s)
	}
}

func (p *selectlangParser) Regex() (localctx IRegexContext) {
	localctx = NewRegexContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 34, selectlangParserRULE_regex)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(144)
		p.Match(selectlangParserREGEX)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

func (p *selectlangParser) Sempred(localctx antlr.RuleContext, ruleIndex, predIndex int) bool {
	switch ruleIndex {
	case 1:
		var t *ExpressionContext = nil
		if localctx != nil {
			t = localctx.(*ExpressionContext)
		}
		return p.Expression_Sempred(t, predIndex)

	default:
		panic("No predicate with index: " + fmt.Sprint(ruleIndex))
	}
}

func (p *selectlangParser) Expression_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 0:
		return p.Precpred(p.GetParserRuleContext(), 2)

	case 1:
		return p.Precpred(p.GetParserRuleContext(), 1)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}
