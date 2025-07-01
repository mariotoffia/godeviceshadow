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
		"", "'SELECT'", "'FROM'", "'WHERE'", "'AND'", "'OR'", "'IN'", "'*'",
		"'.'", "'=='", "'!='", "'>'", "'<'", "'>='", "'<='", "'~='", "'('",
		"')'", "','", "'obj'", "'log'", "'ID'", "'Name'", "'Operation'", "'Path'",
		"'Value'",
	}
	staticData.SymbolicNames = []string{
		"", "SELECT", "FROM", "WHERE", "AND", "OR", "IN", "STAR", "DOT", "EQ",
		"NE", "GT", "LT", "GE", "LE", "REGEX_OP", "LPAREN", "RPAREN", "COMMA",
		"OBJ", "LOG", "ID_FIELD", "NAME_FIELD", "OP_FIELD", "PATH_FIELD", "VAL_FIELD",
		"IDENTIFIER", "NUMBER", "STRING", "WS",
	}
	staticData.RuleNames = []string{
		"select_stmt", "columns", "stream", "where_clause", "expression", "and_expr",
		"primary_expr", "predicate", "value_list", "field", "obj_field", "log_field",
		"value", "comp_operator", "regex_operator", "regex_value",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 29, 128, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7,
		10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15, 7, 15,
		1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 3, 0, 38, 8, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1,
		2, 1, 2, 1, 3, 1, 3, 1, 3, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 5, 4, 55,
		8, 4, 10, 4, 12, 4, 58, 9, 4, 1, 5, 1, 5, 1, 5, 1, 5, 1, 5, 1, 5, 5, 5,
		66, 8, 5, 10, 5, 12, 5, 69, 9, 5, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 3, 6, 76,
		8, 6, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7,
		1, 7, 3, 7, 90, 8, 7, 1, 8, 1, 8, 1, 8, 5, 8, 95, 8, 8, 10, 8, 12, 8, 98,
		9, 8, 1, 9, 1, 9, 3, 9, 102, 8, 9, 1, 10, 1, 10, 1, 10, 1, 10, 1, 11, 1,
		11, 1, 11, 1, 11, 1, 12, 1, 12, 3, 12, 114, 8, 12, 1, 13, 1, 13, 1, 13,
		1, 13, 1, 13, 1, 13, 3, 13, 122, 8, 13, 1, 14, 1, 14, 1, 15, 1, 15, 1,
		15, 0, 2, 8, 10, 16, 0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26,
		28, 30, 0, 2, 1, 0, 21, 23, 1, 0, 22, 25, 125, 0, 32, 1, 0, 0, 0, 2, 41,
		1, 0, 0, 0, 4, 43, 1, 0, 0, 0, 6, 45, 1, 0, 0, 0, 8, 48, 1, 0, 0, 0, 10,
		59, 1, 0, 0, 0, 12, 75, 1, 0, 0, 0, 14, 89, 1, 0, 0, 0, 16, 91, 1, 0, 0,
		0, 18, 101, 1, 0, 0, 0, 20, 103, 1, 0, 0, 0, 22, 107, 1, 0, 0, 0, 24, 113,
		1, 0, 0, 0, 26, 121, 1, 0, 0, 0, 28, 123, 1, 0, 0, 0, 30, 125, 1, 0, 0,
		0, 32, 33, 5, 1, 0, 0, 33, 34, 3, 2, 1, 0, 34, 35, 5, 2, 0, 0, 35, 37,
		3, 4, 2, 0, 36, 38, 3, 6, 3, 0, 37, 36, 1, 0, 0, 0, 37, 38, 1, 0, 0, 0,
		38, 39, 1, 0, 0, 0, 39, 40, 5, 0, 0, 1, 40, 1, 1, 0, 0, 0, 41, 42, 5, 7,
		0, 0, 42, 3, 1, 0, 0, 0, 43, 44, 5, 26, 0, 0, 44, 5, 1, 0, 0, 0, 45, 46,
		5, 3, 0, 0, 46, 47, 3, 8, 4, 0, 47, 7, 1, 0, 0, 0, 48, 49, 6, 4, -1, 0,
		49, 50, 3, 10, 5, 0, 50, 56, 1, 0, 0, 0, 51, 52, 10, 2, 0, 0, 52, 53, 5,
		5, 0, 0, 53, 55, 3, 10, 5, 0, 54, 51, 1, 0, 0, 0, 55, 58, 1, 0, 0, 0, 56,
		54, 1, 0, 0, 0, 56, 57, 1, 0, 0, 0, 57, 9, 1, 0, 0, 0, 58, 56, 1, 0, 0,
		0, 59, 60, 6, 5, -1, 0, 60, 61, 3, 12, 6, 0, 61, 67, 1, 0, 0, 0, 62, 63,
		10, 2, 0, 0, 63, 64, 5, 4, 0, 0, 64, 66, 3, 12, 6, 0, 65, 62, 1, 0, 0,
		0, 66, 69, 1, 0, 0, 0, 67, 65, 1, 0, 0, 0, 67, 68, 1, 0, 0, 0, 68, 11,
		1, 0, 0, 0, 69, 67, 1, 0, 0, 0, 70, 71, 5, 16, 0, 0, 71, 72, 3, 8, 4, 0,
		72, 73, 5, 17, 0, 0, 73, 76, 1, 0, 0, 0, 74, 76, 3, 14, 7, 0, 75, 70, 1,
		0, 0, 0, 75, 74, 1, 0, 0, 0, 76, 13, 1, 0, 0, 0, 77, 78, 3, 18, 9, 0, 78,
		79, 3, 26, 13, 0, 79, 80, 3, 24, 12, 0, 80, 90, 1, 0, 0, 0, 81, 82, 3,
		18, 9, 0, 82, 83, 3, 28, 14, 0, 83, 84, 3, 30, 15, 0, 84, 90, 1, 0, 0,
		0, 85, 86, 3, 18, 9, 0, 86, 87, 5, 6, 0, 0, 87, 88, 3, 16, 8, 0, 88, 90,
		1, 0, 0, 0, 89, 77, 1, 0, 0, 0, 89, 81, 1, 0, 0, 0, 89, 85, 1, 0, 0, 0,
		90, 15, 1, 0, 0, 0, 91, 96, 3, 24, 12, 0, 92, 93, 5, 18, 0, 0, 93, 95,
		3, 24, 12, 0, 94, 92, 1, 0, 0, 0, 95, 98, 1, 0, 0, 0, 96, 94, 1, 0, 0,
		0, 96, 97, 1, 0, 0, 0, 97, 17, 1, 0, 0, 0, 98, 96, 1, 0, 0, 0, 99, 102,
		3, 20, 10, 0, 100, 102, 3, 22, 11, 0, 101, 99, 1, 0, 0, 0, 101, 100, 1,
		0, 0, 0, 102, 19, 1, 0, 0, 0, 103, 104, 5, 19, 0, 0, 104, 105, 5, 8, 0,
		0, 105, 106, 7, 0, 0, 0, 106, 21, 1, 0, 0, 0, 107, 108, 5, 20, 0, 0, 108,
		109, 5, 8, 0, 0, 109, 110, 7, 1, 0, 0, 110, 23, 1, 0, 0, 0, 111, 114, 5,
		27, 0, 0, 112, 114, 5, 28, 0, 0, 113, 111, 1, 0, 0, 0, 113, 112, 1, 0,
		0, 0, 114, 25, 1, 0, 0, 0, 115, 122, 5, 9, 0, 0, 116, 122, 5, 10, 0, 0,
		117, 122, 5, 11, 0, 0, 118, 122, 5, 12, 0, 0, 119, 122, 5, 13, 0, 0, 120,
		122, 5, 14, 0, 0, 121, 115, 1, 0, 0, 0, 121, 116, 1, 0, 0, 0, 121, 117,
		1, 0, 0, 0, 121, 118, 1, 0, 0, 0, 121, 119, 1, 0, 0, 0, 121, 120, 1, 0,
		0, 0, 122, 27, 1, 0, 0, 0, 123, 124, 5, 15, 0, 0, 124, 29, 1, 0, 0, 0,
		125, 126, 5, 28, 0, 0, 126, 31, 1, 0, 0, 0, 9, 37, 56, 67, 75, 89, 96,
		101, 113, 121,
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
	selectlangParserEOF        = antlr.TokenEOF
	selectlangParserSELECT     = 1
	selectlangParserFROM       = 2
	selectlangParserWHERE      = 3
	selectlangParserAND        = 4
	selectlangParserOR         = 5
	selectlangParserIN         = 6
	selectlangParserSTAR       = 7
	selectlangParserDOT        = 8
	selectlangParserEQ         = 9
	selectlangParserNE         = 10
	selectlangParserGT         = 11
	selectlangParserLT         = 12
	selectlangParserGE         = 13
	selectlangParserLE         = 14
	selectlangParserREGEX_OP   = 15
	selectlangParserLPAREN     = 16
	selectlangParserRPAREN     = 17
	selectlangParserCOMMA      = 18
	selectlangParserOBJ        = 19
	selectlangParserLOG        = 20
	selectlangParserID_FIELD   = 21
	selectlangParserNAME_FIELD = 22
	selectlangParserOP_FIELD   = 23
	selectlangParserPATH_FIELD = 24
	selectlangParserVAL_FIELD  = 25
	selectlangParserIDENTIFIER = 26
	selectlangParserNUMBER     = 27
	selectlangParserSTRING     = 28
	selectlangParserWS         = 29
)

// selectlangParser rules.
const (
	selectlangParserRULE_select_stmt    = 0
	selectlangParserRULE_columns        = 1
	selectlangParserRULE_stream         = 2
	selectlangParserRULE_where_clause   = 3
	selectlangParserRULE_expression     = 4
	selectlangParserRULE_and_expr       = 5
	selectlangParserRULE_primary_expr   = 6
	selectlangParserRULE_predicate      = 7
	selectlangParserRULE_value_list     = 8
	selectlangParserRULE_field          = 9
	selectlangParserRULE_obj_field      = 10
	selectlangParserRULE_log_field      = 11
	selectlangParserRULE_value          = 12
	selectlangParserRULE_comp_operator  = 13
	selectlangParserRULE_regex_operator = 14
	selectlangParserRULE_regex_value    = 15
)

// ISelect_stmtContext is an interface to support dynamic dispatch.
type ISelect_stmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsSelect_stmtContext differentiates from other interfaces.
	IsSelect_stmtContext()
}

type Select_stmtContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySelect_stmtContext() *Select_stmtContext {
	var p = new(Select_stmtContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_select_stmt
	return p
}

func InitEmptySelect_stmtContext(p *Select_stmtContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_select_stmt
}

func (*Select_stmtContext) IsSelect_stmtContext() {}

func NewSelect_stmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Select_stmtContext {
	var p = new(Select_stmtContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = selectlangParserRULE_select_stmt

	return p
}

func (s *Select_stmtContext) GetParser() antlr.Parser { return s.parser }

func (s *Select_stmtContext) CopyAll(ctx *Select_stmtContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *Select_stmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Select_stmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type SelectStatementContext struct {
	Select_stmtContext
}

func NewSelectStatementContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SelectStatementContext {
	var p = new(SelectStatementContext)

	InitEmptySelect_stmtContext(&p.Select_stmtContext)
	p.parser = parser
	p.CopyAll(ctx.(*Select_stmtContext))

	return p
}

func (s *SelectStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SelectStatementContext) SELECT() antlr.TerminalNode {
	return s.GetToken(selectlangParserSELECT, 0)
}

func (s *SelectStatementContext) Columns() IColumnsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IColumnsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IColumnsContext)
}

func (s *SelectStatementContext) FROM() antlr.TerminalNode {
	return s.GetToken(selectlangParserFROM, 0)
}

func (s *SelectStatementContext) Stream() IStreamContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStreamContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStreamContext)
}

func (s *SelectStatementContext) EOF() antlr.TerminalNode {
	return s.GetToken(selectlangParserEOF, 0)
}

func (s *SelectStatementContext) Where_clause() IWhere_clauseContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IWhere_clauseContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IWhere_clauseContext)
}

func (s *SelectStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.EnterSelectStatement(s)
	}
}

func (s *SelectStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.ExitSelectStatement(s)
	}
}

func (p *selectlangParser) Select_stmt() (localctx ISelect_stmtContext) {
	localctx = NewSelect_stmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, selectlangParserRULE_select_stmt)
	var _la int

	localctx = NewSelectStatementContext(p, localctx)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(32)
		p.Match(selectlangParserSELECT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(33)
		p.Columns()
	}
	{
		p.SetState(34)
		p.Match(selectlangParserFROM)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(35)
		p.Stream()
	}
	p.SetState(37)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == selectlangParserWHERE {
		{
			p.SetState(36)
			p.Where_clause()
		}

	}
	{
		p.SetState(39)
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

// IColumnsContext is an interface to support dynamic dispatch.
type IColumnsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsColumnsContext differentiates from other interfaces.
	IsColumnsContext()
}

type ColumnsContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyColumnsContext() *ColumnsContext {
	var p = new(ColumnsContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_columns
	return p
}

func InitEmptyColumnsContext(p *ColumnsContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_columns
}

func (*ColumnsContext) IsColumnsContext() {}

func NewColumnsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ColumnsContext {
	var p = new(ColumnsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = selectlangParserRULE_columns

	return p
}

func (s *ColumnsContext) GetParser() antlr.Parser { return s.parser }

func (s *ColumnsContext) CopyAll(ctx *ColumnsContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *ColumnsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ColumnsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type AllColumnsContext struct {
	ColumnsContext
}

func NewAllColumnsContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *AllColumnsContext {
	var p = new(AllColumnsContext)

	InitEmptyColumnsContext(&p.ColumnsContext)
	p.parser = parser
	p.CopyAll(ctx.(*ColumnsContext))

	return p
}

func (s *AllColumnsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AllColumnsContext) STAR() antlr.TerminalNode {
	return s.GetToken(selectlangParserSTAR, 0)
}

func (s *AllColumnsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.EnterAllColumns(s)
	}
}

func (s *AllColumnsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.ExitAllColumns(s)
	}
}

func (p *selectlangParser) Columns() (localctx IColumnsContext) {
	localctx = NewColumnsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, selectlangParserRULE_columns)
	localctx = NewAllColumnsContext(p, localctx)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(41)
		p.Match(selectlangParserSTAR)
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

// IStreamContext is an interface to support dynamic dispatch.
type IStreamContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsStreamContext differentiates from other interfaces.
	IsStreamContext()
}

type StreamContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStreamContext() *StreamContext {
	var p = new(StreamContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_stream
	return p
}

func InitEmptyStreamContext(p *StreamContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_stream
}

func (*StreamContext) IsStreamContext() {}

func NewStreamContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StreamContext {
	var p = new(StreamContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = selectlangParserRULE_stream

	return p
}

func (s *StreamContext) GetParser() antlr.Parser { return s.parser }

func (s *StreamContext) CopyAll(ctx *StreamContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *StreamContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StreamContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type StreamNameContext struct {
	StreamContext
}

func NewStreamNameContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *StreamNameContext {
	var p = new(StreamNameContext)

	InitEmptyStreamContext(&p.StreamContext)
	p.parser = parser
	p.CopyAll(ctx.(*StreamContext))

	return p
}

func (s *StreamNameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StreamNameContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(selectlangParserIDENTIFIER, 0)
}

func (s *StreamNameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.EnterStreamName(s)
	}
}

func (s *StreamNameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.ExitStreamName(s)
	}
}

func (p *selectlangParser) Stream() (localctx IStreamContext) {
	localctx = NewStreamContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, selectlangParserRULE_stream)
	localctx = NewStreamNameContext(p, localctx)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(43)
		p.Match(selectlangParserIDENTIFIER)
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

// IWhere_clauseContext is an interface to support dynamic dispatch.
type IWhere_clauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsWhere_clauseContext differentiates from other interfaces.
	IsWhere_clauseContext()
}

type Where_clauseContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWhere_clauseContext() *Where_clauseContext {
	var p = new(Where_clauseContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_where_clause
	return p
}

func InitEmptyWhere_clauseContext(p *Where_clauseContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_where_clause
}

func (*Where_clauseContext) IsWhere_clauseContext() {}

func NewWhere_clauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Where_clauseContext {
	var p = new(Where_clauseContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = selectlangParserRULE_where_clause

	return p
}

func (s *Where_clauseContext) GetParser() antlr.Parser { return s.parser }

func (s *Where_clauseContext) CopyAll(ctx *Where_clauseContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *Where_clauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Where_clauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type WhereClauseContext struct {
	Where_clauseContext
}

func NewWhereClauseContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *WhereClauseContext {
	var p = new(WhereClauseContext)

	InitEmptyWhere_clauseContext(&p.Where_clauseContext)
	p.parser = parser
	p.CopyAll(ctx.(*Where_clauseContext))

	return p
}

func (s *WhereClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WhereClauseContext) WHERE() antlr.TerminalNode {
	return s.GetToken(selectlangParserWHERE, 0)
}

func (s *WhereClauseContext) Expression() IExpressionContext {
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

func (s *WhereClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.EnterWhereClause(s)
	}
}

func (s *WhereClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.ExitWhereClause(s)
	}
}

func (p *selectlangParser) Where_clause() (localctx IWhere_clauseContext) {
	localctx = NewWhere_clauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, selectlangParserRULE_where_clause)
	localctx = NewWhereClauseContext(p, localctx)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(45)
		p.Match(selectlangParserWHERE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(46)
		p.expression(0)
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

func (s *ExpressionContext) CopyAll(ctx *ExpressionContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *ExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type AndToExpressionContext struct {
	ExpressionContext
}

func NewAndToExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *AndToExpressionContext {
	var p = new(AndToExpressionContext)

	InitEmptyExpressionContext(&p.ExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExpressionContext))

	return p
}

func (s *AndToExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AndToExpressionContext) And_expr() IAnd_exprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAnd_exprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAnd_exprContext)
}

func (s *AndToExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.EnterAndToExpression(s)
	}
}

func (s *AndToExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.ExitAndToExpression(s)
	}
}

type OrExpressionContext struct {
	ExpressionContext
}

func NewOrExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *OrExpressionContext {
	var p = new(OrExpressionContext)

	InitEmptyExpressionContext(&p.ExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExpressionContext))

	return p
}

func (s *OrExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *OrExpressionContext) Expression() IExpressionContext {
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

func (s *OrExpressionContext) OR() antlr.TerminalNode {
	return s.GetToken(selectlangParserOR, 0)
}

func (s *OrExpressionContext) And_expr() IAnd_exprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAnd_exprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAnd_exprContext)
}

func (s *OrExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.EnterOrExpression(s)
	}
}

func (s *OrExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.ExitOrExpression(s)
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
	_startState := 8
	p.EnterRecursionRule(localctx, 8, selectlangParserRULE_expression, _p)
	var _alt int

	p.EnterOuterAlt(localctx, 1)
	localctx = NewAndToExpressionContext(p, localctx)
	p.SetParserRuleContext(localctx)
	_prevctx = localctx

	{
		p.SetState(49)
		p.and_expr(0)
	}

	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(56)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 1, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			localctx = NewOrExpressionContext(p, NewExpressionContext(p, _parentctx, _parentState))
			p.PushNewRecursionContext(localctx, _startState, selectlangParserRULE_expression)
			p.SetState(51)

			if !(p.Precpred(p.GetParserRuleContext(), 2)) {
				p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
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
				p.and_expr(0)
			}

		}
		p.SetState(58)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 1, p.GetParserRuleContext())
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

// IAnd_exprContext is an interface to support dynamic dispatch.
type IAnd_exprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsAnd_exprContext differentiates from other interfaces.
	IsAnd_exprContext()
}

type And_exprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAnd_exprContext() *And_exprContext {
	var p = new(And_exprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_and_expr
	return p
}

func InitEmptyAnd_exprContext(p *And_exprContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_and_expr
}

func (*And_exprContext) IsAnd_exprContext() {}

func NewAnd_exprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *And_exprContext {
	var p = new(And_exprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = selectlangParserRULE_and_expr

	return p
}

func (s *And_exprContext) GetParser() antlr.Parser { return s.parser }

func (s *And_exprContext) CopyAll(ctx *And_exprContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *And_exprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *And_exprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type AndExpressionContext struct {
	And_exprContext
}

func NewAndExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *AndExpressionContext {
	var p = new(AndExpressionContext)

	InitEmptyAnd_exprContext(&p.And_exprContext)
	p.parser = parser
	p.CopyAll(ctx.(*And_exprContext))

	return p
}

func (s *AndExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AndExpressionContext) And_expr() IAnd_exprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAnd_exprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAnd_exprContext)
}

func (s *AndExpressionContext) AND() antlr.TerminalNode {
	return s.GetToken(selectlangParserAND, 0)
}

func (s *AndExpressionContext) Primary_expr() IPrimary_exprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPrimary_exprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPrimary_exprContext)
}

func (s *AndExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.EnterAndExpression(s)
	}
}

func (s *AndExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.ExitAndExpression(s)
	}
}

type PrimaryExpressionContext struct {
	And_exprContext
}

func NewPrimaryExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *PrimaryExpressionContext {
	var p = new(PrimaryExpressionContext)

	InitEmptyAnd_exprContext(&p.And_exprContext)
	p.parser = parser
	p.CopyAll(ctx.(*And_exprContext))

	return p
}

func (s *PrimaryExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PrimaryExpressionContext) Primary_expr() IPrimary_exprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPrimary_exprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPrimary_exprContext)
}

func (s *PrimaryExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.EnterPrimaryExpression(s)
	}
}

func (s *PrimaryExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.ExitPrimaryExpression(s)
	}
}

func (p *selectlangParser) And_expr() (localctx IAnd_exprContext) {
	return p.and_expr(0)
}

func (p *selectlangParser) and_expr(_p int) (localctx IAnd_exprContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()

	_parentState := p.GetState()
	localctx = NewAnd_exprContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IAnd_exprContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 10
	p.EnterRecursionRule(localctx, 10, selectlangParserRULE_and_expr, _p)
	var _alt int

	p.EnterOuterAlt(localctx, 1)
	localctx = NewPrimaryExpressionContext(p, localctx)
	p.SetParserRuleContext(localctx)
	_prevctx = localctx

	{
		p.SetState(60)
		p.Primary_expr()
	}

	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(67)
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
			localctx = NewAndExpressionContext(p, NewAnd_exprContext(p, _parentctx, _parentState))
			p.PushNewRecursionContext(localctx, _startState, selectlangParserRULE_and_expr)
			p.SetState(62)

			if !(p.Precpred(p.GetParserRuleContext(), 2)) {
				p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
				goto errorExit
			}
			{
				p.SetState(63)
				p.Match(selectlangParserAND)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(64)
				p.Primary_expr()
			}

		}
		p.SetState(69)
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

// IPrimary_exprContext is an interface to support dynamic dispatch.
type IPrimary_exprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsPrimary_exprContext differentiates from other interfaces.
	IsPrimary_exprContext()
}

type Primary_exprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPrimary_exprContext() *Primary_exprContext {
	var p = new(Primary_exprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_primary_expr
	return p
}

func InitEmptyPrimary_exprContext(p *Primary_exprContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_primary_expr
}

func (*Primary_exprContext) IsPrimary_exprContext() {}

func NewPrimary_exprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Primary_exprContext {
	var p = new(Primary_exprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = selectlangParserRULE_primary_expr

	return p
}

func (s *Primary_exprContext) GetParser() antlr.Parser { return s.parser }

func (s *Primary_exprContext) CopyAll(ctx *Primary_exprContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *Primary_exprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Primary_exprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type ParenExpressionContext struct {
	Primary_exprContext
}

func NewParenExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ParenExpressionContext {
	var p = new(ParenExpressionContext)

	InitEmptyPrimary_exprContext(&p.Primary_exprContext)
	p.parser = parser
	p.CopyAll(ctx.(*Primary_exprContext))

	return p
}

func (s *ParenExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ParenExpressionContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(selectlangParserLPAREN, 0)
}

func (s *ParenExpressionContext) Expression() IExpressionContext {
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

func (s *ParenExpressionContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(selectlangParserRPAREN, 0)
}

func (s *ParenExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.EnterParenExpression(s)
	}
}

func (s *ParenExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.ExitParenExpression(s)
	}
}

type PredicateExpressionContext struct {
	Primary_exprContext
}

func NewPredicateExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *PredicateExpressionContext {
	var p = new(PredicateExpressionContext)

	InitEmptyPrimary_exprContext(&p.Primary_exprContext)
	p.parser = parser
	p.CopyAll(ctx.(*Primary_exprContext))

	return p
}

func (s *PredicateExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PredicateExpressionContext) Predicate() IPredicateContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPredicateContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPredicateContext)
}

func (s *PredicateExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.EnterPredicateExpression(s)
	}
}

func (s *PredicateExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.ExitPredicateExpression(s)
	}
}

func (p *selectlangParser) Primary_expr() (localctx IPrimary_exprContext) {
	localctx = NewPrimary_exprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, selectlangParserRULE_primary_expr)
	p.SetState(75)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case selectlangParserLPAREN:
		localctx = NewParenExpressionContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(70)
			p.Match(selectlangParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(71)
			p.expression(0)
		}
		{
			p.SetState(72)
			p.Match(selectlangParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case selectlangParserOBJ, selectlangParserLOG:
		localctx = NewPredicateExpressionContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(74)
			p.Predicate()
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

// IPredicateContext is an interface to support dynamic dispatch.
type IPredicateContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsPredicateContext differentiates from other interfaces.
	IsPredicateContext()
}

type PredicateContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPredicateContext() *PredicateContext {
	var p = new(PredicateContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_predicate
	return p
}

func InitEmptyPredicateContext(p *PredicateContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_predicate
}

func (*PredicateContext) IsPredicateContext() {}

func NewPredicateContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PredicateContext {
	var p = new(PredicateContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = selectlangParserRULE_predicate

	return p
}

func (s *PredicateContext) GetParser() antlr.Parser { return s.parser }

func (s *PredicateContext) CopyAll(ctx *PredicateContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *PredicateContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PredicateContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type RegexPredicateContext struct {
	PredicateContext
}

func NewRegexPredicateContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *RegexPredicateContext {
	var p = new(RegexPredicateContext)

	InitEmptyPredicateContext(&p.PredicateContext)
	p.parser = parser
	p.CopyAll(ctx.(*PredicateContext))

	return p
}

func (s *RegexPredicateContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RegexPredicateContext) Field() IFieldContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFieldContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFieldContext)
}

func (s *RegexPredicateContext) Regex_operator() IRegex_operatorContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRegex_operatorContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRegex_operatorContext)
}

func (s *RegexPredicateContext) Regex_value() IRegex_valueContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRegex_valueContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRegex_valueContext)
}

func (s *RegexPredicateContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.EnterRegexPredicate(s)
	}
}

func (s *RegexPredicateContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.ExitRegexPredicate(s)
	}
}

type InPredicateContext struct {
	PredicateContext
}

func NewInPredicateContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *InPredicateContext {
	var p = new(InPredicateContext)

	InitEmptyPredicateContext(&p.PredicateContext)
	p.parser = parser
	p.CopyAll(ctx.(*PredicateContext))

	return p
}

func (s *InPredicateContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *InPredicateContext) Field() IFieldContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFieldContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFieldContext)
}

func (s *InPredicateContext) IN() antlr.TerminalNode {
	return s.GetToken(selectlangParserIN, 0)
}

func (s *InPredicateContext) Value_list() IValue_listContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IValue_listContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IValue_listContext)
}

func (s *InPredicateContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.EnterInPredicate(s)
	}
}

func (s *InPredicateContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.ExitInPredicate(s)
	}
}

type ComparisonPredicateContext struct {
	PredicateContext
}

func NewComparisonPredicateContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ComparisonPredicateContext {
	var p = new(ComparisonPredicateContext)

	InitEmptyPredicateContext(&p.PredicateContext)
	p.parser = parser
	p.CopyAll(ctx.(*PredicateContext))

	return p
}

func (s *ComparisonPredicateContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ComparisonPredicateContext) Field() IFieldContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFieldContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFieldContext)
}

func (s *ComparisonPredicateContext) Comp_operator() IComp_operatorContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IComp_operatorContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IComp_operatorContext)
}

func (s *ComparisonPredicateContext) Value() IValueContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IValueContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IValueContext)
}

func (s *ComparisonPredicateContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.EnterComparisonPredicate(s)
	}
}

func (s *ComparisonPredicateContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.ExitComparisonPredicate(s)
	}
}

func (p *selectlangParser) Predicate() (localctx IPredicateContext) {
	localctx = NewPredicateContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, selectlangParserRULE_predicate)
	p.SetState(89)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 4, p.GetParserRuleContext()) {
	case 1:
		localctx = NewComparisonPredicateContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(77)
			p.Field()
		}
		{
			p.SetState(78)
			p.Comp_operator()
		}
		{
			p.SetState(79)
			p.Value()
		}

	case 2:
		localctx = NewRegexPredicateContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(81)
			p.Field()
		}
		{
			p.SetState(82)
			p.Regex_operator()
		}
		{
			p.SetState(83)
			p.Regex_value()
		}

	case 3:
		localctx = NewInPredicateContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(85)
			p.Field()
		}
		{
			p.SetState(86)
			p.Match(selectlangParserIN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(87)
			p.Value_list()
		}

	case antlr.ATNInvalidAltNumber:
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

// IValue_listContext is an interface to support dynamic dispatch.
type IValue_listContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsValue_listContext differentiates from other interfaces.
	IsValue_listContext()
}

type Value_listContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyValue_listContext() *Value_listContext {
	var p = new(Value_listContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_value_list
	return p
}

func InitEmptyValue_listContext(p *Value_listContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_value_list
}

func (*Value_listContext) IsValue_listContext() {}

func NewValue_listContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Value_listContext {
	var p = new(Value_listContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = selectlangParserRULE_value_list

	return p
}

func (s *Value_listContext) GetParser() antlr.Parser { return s.parser }

func (s *Value_listContext) CopyAll(ctx *Value_listContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *Value_listContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Value_listContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type ValueListContext struct {
	Value_listContext
}

func NewValueListContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ValueListContext {
	var p = new(ValueListContext)

	InitEmptyValue_listContext(&p.Value_listContext)
	p.parser = parser
	p.CopyAll(ctx.(*Value_listContext))

	return p
}

func (s *ValueListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ValueListContext) AllValue() []IValueContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IValueContext); ok {
			len++
		}
	}

	tst := make([]IValueContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IValueContext); ok {
			tst[i] = t.(IValueContext)
			i++
		}
	}

	return tst
}

func (s *ValueListContext) Value(i int) IValueContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IValueContext); ok {
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

	return t.(IValueContext)
}

func (s *ValueListContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(selectlangParserCOMMA)
}

func (s *ValueListContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(selectlangParserCOMMA, i)
}

func (s *ValueListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.EnterValueList(s)
	}
}

func (s *ValueListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.ExitValueList(s)
	}
}

func (p *selectlangParser) Value_list() (localctx IValue_listContext) {
	localctx = NewValue_listContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, selectlangParserRULE_value_list)
	var _alt int

	localctx = NewValueListContext(p, localctx)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(91)
		p.Value()
	}
	p.SetState(96)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 5, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(92)
				p.Match(selectlangParserCOMMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(93)
				p.Value()
			}

		}
		p.SetState(98)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 5, p.GetParserRuleContext())
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

// IFieldContext is an interface to support dynamic dispatch.
type IFieldContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsFieldContext differentiates from other interfaces.
	IsFieldContext()
}

type FieldContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFieldContext() *FieldContext {
	var p = new(FieldContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_field
	return p
}

func InitEmptyFieldContext(p *FieldContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_field
}

func (*FieldContext) IsFieldContext() {}

func NewFieldContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FieldContext {
	var p = new(FieldContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = selectlangParserRULE_field

	return p
}

func (s *FieldContext) GetParser() antlr.Parser { return s.parser }

func (s *FieldContext) CopyAll(ctx *FieldContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *FieldContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FieldContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type LogFieldContext struct {
	FieldContext
}

func NewLogFieldContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *LogFieldContext {
	var p = new(LogFieldContext)

	InitEmptyFieldContext(&p.FieldContext)
	p.parser = parser
	p.CopyAll(ctx.(*FieldContext))

	return p
}

func (s *LogFieldContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LogFieldContext) Log_field() ILog_fieldContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILog_fieldContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILog_fieldContext)
}

func (s *LogFieldContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.EnterLogField(s)
	}
}

func (s *LogFieldContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.ExitLogField(s)
	}
}

type ObjFieldContext struct {
	FieldContext
}

func NewObjFieldContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ObjFieldContext {
	var p = new(ObjFieldContext)

	InitEmptyFieldContext(&p.FieldContext)
	p.parser = parser
	p.CopyAll(ctx.(*FieldContext))

	return p
}

func (s *ObjFieldContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ObjFieldContext) Obj_field() IObj_fieldContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IObj_fieldContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IObj_fieldContext)
}

func (s *ObjFieldContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.EnterObjField(s)
	}
}

func (s *ObjFieldContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.ExitObjField(s)
	}
}

func (p *selectlangParser) Field() (localctx IFieldContext) {
	localctx = NewFieldContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, selectlangParserRULE_field)
	p.SetState(101)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case selectlangParserOBJ:
		localctx = NewObjFieldContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(99)
			p.Obj_field()
		}

	case selectlangParserLOG:
		localctx = NewLogFieldContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(100)
			p.Log_field()
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

// IObj_fieldContext is an interface to support dynamic dispatch.
type IObj_fieldContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsObj_fieldContext differentiates from other interfaces.
	IsObj_fieldContext()
}

type Obj_fieldContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyObj_fieldContext() *Obj_fieldContext {
	var p = new(Obj_fieldContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_obj_field
	return p
}

func InitEmptyObj_fieldContext(p *Obj_fieldContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_obj_field
}

func (*Obj_fieldContext) IsObj_fieldContext() {}

func NewObj_fieldContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Obj_fieldContext {
	var p = new(Obj_fieldContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = selectlangParserRULE_obj_field

	return p
}

func (s *Obj_fieldContext) GetParser() antlr.Parser { return s.parser }

func (s *Obj_fieldContext) CopyAll(ctx *Obj_fieldContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *Obj_fieldContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Obj_fieldContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type ObjFieldAccessContext struct {
	Obj_fieldContext
}

func NewObjFieldAccessContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ObjFieldAccessContext {
	var p = new(ObjFieldAccessContext)

	InitEmptyObj_fieldContext(&p.Obj_fieldContext)
	p.parser = parser
	p.CopyAll(ctx.(*Obj_fieldContext))

	return p
}

func (s *ObjFieldAccessContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ObjFieldAccessContext) OBJ() antlr.TerminalNode {
	return s.GetToken(selectlangParserOBJ, 0)
}

func (s *ObjFieldAccessContext) DOT() antlr.TerminalNode {
	return s.GetToken(selectlangParserDOT, 0)
}

func (s *ObjFieldAccessContext) ID_FIELD() antlr.TerminalNode {
	return s.GetToken(selectlangParserID_FIELD, 0)
}

func (s *ObjFieldAccessContext) NAME_FIELD() antlr.TerminalNode {
	return s.GetToken(selectlangParserNAME_FIELD, 0)
}

func (s *ObjFieldAccessContext) OP_FIELD() antlr.TerminalNode {
	return s.GetToken(selectlangParserOP_FIELD, 0)
}

func (s *ObjFieldAccessContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.EnterObjFieldAccess(s)
	}
}

func (s *ObjFieldAccessContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.ExitObjFieldAccess(s)
	}
}

func (p *selectlangParser) Obj_field() (localctx IObj_fieldContext) {
	localctx = NewObj_fieldContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, selectlangParserRULE_obj_field)
	var _la int

	localctx = NewObjFieldAccessContext(p, localctx)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(103)
		p.Match(selectlangParserOBJ)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(104)
		p.Match(selectlangParserDOT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(105)
		_la = p.GetTokenStream().LA(1)

		if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&14680064) != 0) {
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

// ILog_fieldContext is an interface to support dynamic dispatch.
type ILog_fieldContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsLog_fieldContext differentiates from other interfaces.
	IsLog_fieldContext()
}

type Log_fieldContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLog_fieldContext() *Log_fieldContext {
	var p = new(Log_fieldContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_log_field
	return p
}

func InitEmptyLog_fieldContext(p *Log_fieldContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_log_field
}

func (*Log_fieldContext) IsLog_fieldContext() {}

func NewLog_fieldContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Log_fieldContext {
	var p = new(Log_fieldContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = selectlangParserRULE_log_field

	return p
}

func (s *Log_fieldContext) GetParser() antlr.Parser { return s.parser }

func (s *Log_fieldContext) CopyAll(ctx *Log_fieldContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *Log_fieldContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Log_fieldContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type LogFieldAccessContext struct {
	Log_fieldContext
}

func NewLogFieldAccessContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *LogFieldAccessContext {
	var p = new(LogFieldAccessContext)

	InitEmptyLog_fieldContext(&p.Log_fieldContext)
	p.parser = parser
	p.CopyAll(ctx.(*Log_fieldContext))

	return p
}

func (s *LogFieldAccessContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LogFieldAccessContext) LOG() antlr.TerminalNode {
	return s.GetToken(selectlangParserLOG, 0)
}

func (s *LogFieldAccessContext) DOT() antlr.TerminalNode {
	return s.GetToken(selectlangParserDOT, 0)
}

func (s *LogFieldAccessContext) OP_FIELD() antlr.TerminalNode {
	return s.GetToken(selectlangParserOP_FIELD, 0)
}

func (s *LogFieldAccessContext) PATH_FIELD() antlr.TerminalNode {
	return s.GetToken(selectlangParserPATH_FIELD, 0)
}

func (s *LogFieldAccessContext) NAME_FIELD() antlr.TerminalNode {
	return s.GetToken(selectlangParserNAME_FIELD, 0)
}

func (s *LogFieldAccessContext) VAL_FIELD() antlr.TerminalNode {
	return s.GetToken(selectlangParserVAL_FIELD, 0)
}

func (s *LogFieldAccessContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.EnterLogFieldAccess(s)
	}
}

func (s *LogFieldAccessContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.ExitLogFieldAccess(s)
	}
}

func (p *selectlangParser) Log_field() (localctx ILog_fieldContext) {
	localctx = NewLog_fieldContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, selectlangParserRULE_log_field)
	var _la int

	localctx = NewLogFieldAccessContext(p, localctx)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(107)
		p.Match(selectlangParserLOG)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(108)
		p.Match(selectlangParserDOT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(109)
		_la = p.GetTokenStream().LA(1)

		if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&62914560) != 0) {
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

// IValueContext is an interface to support dynamic dispatch.
type IValueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsValueContext differentiates from other interfaces.
	IsValueContext()
}

type ValueContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyValueContext() *ValueContext {
	var p = new(ValueContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_value
	return p
}

func InitEmptyValueContext(p *ValueContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_value
}

func (*ValueContext) IsValueContext() {}

func NewValueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ValueContext {
	var p = new(ValueContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = selectlangParserRULE_value

	return p
}

func (s *ValueContext) GetParser() antlr.Parser { return s.parser }

func (s *ValueContext) CopyAll(ctx *ValueContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *ValueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ValueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type NumberValueContext struct {
	ValueContext
}

func NewNumberValueContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *NumberValueContext {
	var p = new(NumberValueContext)

	InitEmptyValueContext(&p.ValueContext)
	p.parser = parser
	p.CopyAll(ctx.(*ValueContext))

	return p
}

func (s *NumberValueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NumberValueContext) NUMBER() antlr.TerminalNode {
	return s.GetToken(selectlangParserNUMBER, 0)
}

func (s *NumberValueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.EnterNumberValue(s)
	}
}

func (s *NumberValueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.ExitNumberValue(s)
	}
}

type StringValueContext struct {
	ValueContext
}

func NewStringValueContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *StringValueContext {
	var p = new(StringValueContext)

	InitEmptyValueContext(&p.ValueContext)
	p.parser = parser
	p.CopyAll(ctx.(*ValueContext))

	return p
}

func (s *StringValueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StringValueContext) STRING() antlr.TerminalNode {
	return s.GetToken(selectlangParserSTRING, 0)
}

func (s *StringValueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.EnterStringValue(s)
	}
}

func (s *StringValueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.ExitStringValue(s)
	}
}

func (p *selectlangParser) Value() (localctx IValueContext) {
	localctx = NewValueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, selectlangParserRULE_value)
	p.SetState(113)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case selectlangParserNUMBER:
		localctx = NewNumberValueContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(111)
			p.Match(selectlangParserNUMBER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case selectlangParserSTRING:
		localctx = NewStringValueContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(112)
			p.Match(selectlangParserSTRING)
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

// IComp_operatorContext is an interface to support dynamic dispatch.
type IComp_operatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsComp_operatorContext differentiates from other interfaces.
	IsComp_operatorContext()
}

type Comp_operatorContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyComp_operatorContext() *Comp_operatorContext {
	var p = new(Comp_operatorContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_comp_operator
	return p
}

func InitEmptyComp_operatorContext(p *Comp_operatorContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_comp_operator
}

func (*Comp_operatorContext) IsComp_operatorContext() {}

func NewComp_operatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Comp_operatorContext {
	var p = new(Comp_operatorContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = selectlangParserRULE_comp_operator

	return p
}

func (s *Comp_operatorContext) GetParser() antlr.Parser { return s.parser }

func (s *Comp_operatorContext) CopyAll(ctx *Comp_operatorContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *Comp_operatorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Comp_operatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type EqualsOpContext struct {
	Comp_operatorContext
}

func NewEqualsOpContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *EqualsOpContext {
	var p = new(EqualsOpContext)

	InitEmptyComp_operatorContext(&p.Comp_operatorContext)
	p.parser = parser
	p.CopyAll(ctx.(*Comp_operatorContext))

	return p
}

func (s *EqualsOpContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EqualsOpContext) EQ() antlr.TerminalNode {
	return s.GetToken(selectlangParserEQ, 0)
}

func (s *EqualsOpContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.EnterEqualsOp(s)
	}
}

func (s *EqualsOpContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.ExitEqualsOp(s)
	}
}

type GreaterThanOpContext struct {
	Comp_operatorContext
}

func NewGreaterThanOpContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *GreaterThanOpContext {
	var p = new(GreaterThanOpContext)

	InitEmptyComp_operatorContext(&p.Comp_operatorContext)
	p.parser = parser
	p.CopyAll(ctx.(*Comp_operatorContext))

	return p
}

func (s *GreaterThanOpContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *GreaterThanOpContext) GT() antlr.TerminalNode {
	return s.GetToken(selectlangParserGT, 0)
}

func (s *GreaterThanOpContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.EnterGreaterThanOp(s)
	}
}

func (s *GreaterThanOpContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.ExitGreaterThanOp(s)
	}
}

type LessOrEqualOpContext struct {
	Comp_operatorContext
}

func NewLessOrEqualOpContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *LessOrEqualOpContext {
	var p = new(LessOrEqualOpContext)

	InitEmptyComp_operatorContext(&p.Comp_operatorContext)
	p.parser = parser
	p.CopyAll(ctx.(*Comp_operatorContext))

	return p
}

func (s *LessOrEqualOpContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LessOrEqualOpContext) LE() antlr.TerminalNode {
	return s.GetToken(selectlangParserLE, 0)
}

func (s *LessOrEqualOpContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.EnterLessOrEqualOp(s)
	}
}

func (s *LessOrEqualOpContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.ExitLessOrEqualOp(s)
	}
}

type LessThanOpContext struct {
	Comp_operatorContext
}

func NewLessThanOpContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *LessThanOpContext {
	var p = new(LessThanOpContext)

	InitEmptyComp_operatorContext(&p.Comp_operatorContext)
	p.parser = parser
	p.CopyAll(ctx.(*Comp_operatorContext))

	return p
}

func (s *LessThanOpContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LessThanOpContext) LT() antlr.TerminalNode {
	return s.GetToken(selectlangParserLT, 0)
}

func (s *LessThanOpContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.EnterLessThanOp(s)
	}
}

func (s *LessThanOpContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.ExitLessThanOp(s)
	}
}

type NotEqualsOpContext struct {
	Comp_operatorContext
}

func NewNotEqualsOpContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *NotEqualsOpContext {
	var p = new(NotEqualsOpContext)

	InitEmptyComp_operatorContext(&p.Comp_operatorContext)
	p.parser = parser
	p.CopyAll(ctx.(*Comp_operatorContext))

	return p
}

func (s *NotEqualsOpContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NotEqualsOpContext) NE() antlr.TerminalNode {
	return s.GetToken(selectlangParserNE, 0)
}

func (s *NotEqualsOpContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.EnterNotEqualsOp(s)
	}
}

func (s *NotEqualsOpContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.ExitNotEqualsOp(s)
	}
}

type GreaterOrEqualOpContext struct {
	Comp_operatorContext
}

func NewGreaterOrEqualOpContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *GreaterOrEqualOpContext {
	var p = new(GreaterOrEqualOpContext)

	InitEmptyComp_operatorContext(&p.Comp_operatorContext)
	p.parser = parser
	p.CopyAll(ctx.(*Comp_operatorContext))

	return p
}

func (s *GreaterOrEqualOpContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *GreaterOrEqualOpContext) GE() antlr.TerminalNode {
	return s.GetToken(selectlangParserGE, 0)
}

func (s *GreaterOrEqualOpContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.EnterGreaterOrEqualOp(s)
	}
}

func (s *GreaterOrEqualOpContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.ExitGreaterOrEqualOp(s)
	}
}

func (p *selectlangParser) Comp_operator() (localctx IComp_operatorContext) {
	localctx = NewComp_operatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, selectlangParserRULE_comp_operator)
	p.SetState(121)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case selectlangParserEQ:
		localctx = NewEqualsOpContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(115)
			p.Match(selectlangParserEQ)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case selectlangParserNE:
		localctx = NewNotEqualsOpContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(116)
			p.Match(selectlangParserNE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case selectlangParserGT:
		localctx = NewGreaterThanOpContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(117)
			p.Match(selectlangParserGT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case selectlangParserLT:
		localctx = NewLessThanOpContext(p, localctx)
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(118)
			p.Match(selectlangParserLT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case selectlangParserGE:
		localctx = NewGreaterOrEqualOpContext(p, localctx)
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(119)
			p.Match(selectlangParserGE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case selectlangParserLE:
		localctx = NewLessOrEqualOpContext(p, localctx)
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(120)
			p.Match(selectlangParserLE)
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

// IRegex_operatorContext is an interface to support dynamic dispatch.
type IRegex_operatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsRegex_operatorContext differentiates from other interfaces.
	IsRegex_operatorContext()
}

type Regex_operatorContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRegex_operatorContext() *Regex_operatorContext {
	var p = new(Regex_operatorContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_regex_operator
	return p
}

func InitEmptyRegex_operatorContext(p *Regex_operatorContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_regex_operator
}

func (*Regex_operatorContext) IsRegex_operatorContext() {}

func NewRegex_operatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Regex_operatorContext {
	var p = new(Regex_operatorContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = selectlangParserRULE_regex_operator

	return p
}

func (s *Regex_operatorContext) GetParser() antlr.Parser { return s.parser }

func (s *Regex_operatorContext) CopyAll(ctx *Regex_operatorContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *Regex_operatorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Regex_operatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type RegexOpContext struct {
	Regex_operatorContext
}

func NewRegexOpContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *RegexOpContext {
	var p = new(RegexOpContext)

	InitEmptyRegex_operatorContext(&p.Regex_operatorContext)
	p.parser = parser
	p.CopyAll(ctx.(*Regex_operatorContext))

	return p
}

func (s *RegexOpContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RegexOpContext) REGEX_OP() antlr.TerminalNode {
	return s.GetToken(selectlangParserREGEX_OP, 0)
}

func (s *RegexOpContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.EnterRegexOp(s)
	}
}

func (s *RegexOpContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.ExitRegexOp(s)
	}
}

func (p *selectlangParser) Regex_operator() (localctx IRegex_operatorContext) {
	localctx = NewRegex_operatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, selectlangParserRULE_regex_operator)
	localctx = NewRegexOpContext(p, localctx)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(123)
		p.Match(selectlangParserREGEX_OP)
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

// IRegex_valueContext is an interface to support dynamic dispatch.
type IRegex_valueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsRegex_valueContext differentiates from other interfaces.
	IsRegex_valueContext()
}

type Regex_valueContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRegex_valueContext() *Regex_valueContext {
	var p = new(Regex_valueContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_regex_value
	return p
}

func InitEmptyRegex_valueContext(p *Regex_valueContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = selectlangParserRULE_regex_value
}

func (*Regex_valueContext) IsRegex_valueContext() {}

func NewRegex_valueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Regex_valueContext {
	var p = new(Regex_valueContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = selectlangParserRULE_regex_value

	return p
}

func (s *Regex_valueContext) GetParser() antlr.Parser { return s.parser }

func (s *Regex_valueContext) CopyAll(ctx *Regex_valueContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *Regex_valueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Regex_valueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type RegexValueContext struct {
	Regex_valueContext
}

func NewRegexValueContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *RegexValueContext {
	var p = new(RegexValueContext)

	InitEmptyRegex_valueContext(&p.Regex_valueContext)
	p.parser = parser
	p.CopyAll(ctx.(*Regex_valueContext))

	return p
}

func (s *RegexValueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RegexValueContext) STRING() antlr.TerminalNode {
	return s.GetToken(selectlangParserSTRING, 0)
}

func (s *RegexValueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.EnterRegexValue(s)
	}
}

func (s *RegexValueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(selectlangListener); ok {
		listenerT.ExitRegexValue(s)
	}
}

func (p *selectlangParser) Regex_value() (localctx IRegex_valueContext) {
	localctx = NewRegex_valueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, selectlangParserRULE_regex_value)
	localctx = NewRegexValueContext(p, localctx)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(125)
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

func (p *selectlangParser) Sempred(localctx antlr.RuleContext, ruleIndex, predIndex int) bool {
	switch ruleIndex {
	case 4:
		var t *ExpressionContext = nil
		if localctx != nil {
			t = localctx.(*ExpressionContext)
		}
		return p.Expression_Sempred(t, predIndex)

	case 5:
		var t *And_exprContext = nil
		if localctx != nil {
			t = localctx.(*And_exprContext)
		}
		return p.And_expr_Sempred(t, predIndex)

	default:
		panic("No predicate with index: " + fmt.Sprint(ruleIndex))
	}
}

func (p *selectlangParser) Expression_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 0:
		return p.Precpred(p.GetParserRuleContext(), 2)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}

func (p *selectlangParser) And_expr_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 1:
		return p.Precpred(p.GetParserRuleContext(), 2)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}
