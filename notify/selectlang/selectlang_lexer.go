// Code generated from selectlang.g4 by ANTLR 4.13.2. DO NOT EDIT.

package selectlang

import (
	"fmt"
	"github.com/antlr4-go/antlr/v4"
	"sync"
	"unicode"
)

// Suppress unused import error
var _ = fmt.Printf
var _ = sync.Once{}
var _ = unicode.IsLetter

type selectlangLexer struct {
	*antlr.BaseLexer
	channelNames []string
	modeNames    []string
	// TODO: EOF string
}

var SelectlangLexerLexerStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	ChannelNames           []string
	ModeNames              []string
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func selectlanglexerLexerInit() {
	staticData := &SelectlangLexerLexerStaticData
	staticData.ChannelNames = []string{
		"DEFAULT_TOKEN_CHANNEL", "HIDDEN",
	}
	staticData.ModeNames = []string{
		"DEFAULT_MODE",
	}
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
		"SELECT", "FROM", "WHERE", "AND", "OR", "IN", "STAR", "DOT", "EQ", "NE",
		"GT", "LT", "GE", "LE", "REGEX_OP", "LPAREN", "RPAREN", "COMMA", "OBJ",
		"LOG", "ID_FIELD", "NAME_FIELD", "OP_FIELD", "PATH_FIELD", "VAL_FIELD",
		"IDENTIFIER", "NUMBER", "STRING", "WS",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 0, 29, 191, 6, -1, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2,
		4, 7, 4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2,
		10, 7, 10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15,
		7, 15, 2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 2, 19, 7, 19, 2, 20, 7,
		20, 2, 21, 7, 21, 2, 22, 7, 22, 2, 23, 7, 23, 2, 24, 7, 24, 2, 25, 7, 25,
		2, 26, 7, 26, 2, 27, 7, 27, 2, 28, 7, 28, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0,
		1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2,
		1, 2, 1, 3, 1, 3, 1, 3, 1, 3, 1, 4, 1, 4, 1, 4, 1, 5, 1, 5, 1, 5, 1, 6,
		1, 6, 1, 7, 1, 7, 1, 8, 1, 8, 1, 8, 1, 9, 1, 9, 1, 9, 1, 10, 1, 10, 1,
		11, 1, 11, 1, 12, 1, 12, 1, 12, 1, 13, 1, 13, 1, 13, 1, 14, 1, 14, 1, 14,
		1, 15, 1, 15, 1, 16, 1, 16, 1, 17, 1, 17, 1, 18, 1, 18, 1, 18, 1, 18, 1,
		19, 1, 19, 1, 19, 1, 19, 1, 20, 1, 20, 1, 20, 1, 21, 1, 21, 1, 21, 1, 21,
		1, 21, 1, 22, 1, 22, 1, 22, 1, 22, 1, 22, 1, 22, 1, 22, 1, 22, 1, 22, 1,
		22, 1, 23, 1, 23, 1, 23, 1, 23, 1, 23, 1, 24, 1, 24, 1, 24, 1, 24, 1, 24,
		1, 24, 1, 25, 1, 25, 5, 25, 156, 8, 25, 10, 25, 12, 25, 159, 9, 25, 1,
		26, 4, 26, 162, 8, 26, 11, 26, 12, 26, 163, 1, 26, 1, 26, 4, 26, 168, 8,
		26, 11, 26, 12, 26, 169, 3, 26, 172, 8, 26, 1, 27, 1, 27, 1, 27, 1, 27,
		5, 27, 178, 8, 27, 10, 27, 12, 27, 181, 9, 27, 1, 27, 1, 27, 1, 28, 4,
		28, 186, 8, 28, 11, 28, 12, 28, 187, 1, 28, 1, 28, 0, 0, 29, 1, 1, 3, 2,
		5, 3, 7, 4, 9, 5, 11, 6, 13, 7, 15, 8, 17, 9, 19, 10, 21, 11, 23, 12, 25,
		13, 27, 14, 29, 15, 31, 16, 33, 17, 35, 18, 37, 19, 39, 20, 41, 21, 43,
		22, 45, 23, 47, 24, 49, 25, 51, 26, 53, 27, 55, 28, 57, 29, 1, 0, 5, 3,
		0, 65, 90, 95, 95, 97, 122, 4, 0, 48, 57, 65, 90, 95, 95, 97, 122, 1, 0,
		48, 57, 2, 0, 39, 39, 92, 92, 3, 0, 9, 10, 13, 13, 32, 32, 197, 0, 1, 1,
		0, 0, 0, 0, 3, 1, 0, 0, 0, 0, 5, 1, 0, 0, 0, 0, 7, 1, 0, 0, 0, 0, 9, 1,
		0, 0, 0, 0, 11, 1, 0, 0, 0, 0, 13, 1, 0, 0, 0, 0, 15, 1, 0, 0, 0, 0, 17,
		1, 0, 0, 0, 0, 19, 1, 0, 0, 0, 0, 21, 1, 0, 0, 0, 0, 23, 1, 0, 0, 0, 0,
		25, 1, 0, 0, 0, 0, 27, 1, 0, 0, 0, 0, 29, 1, 0, 0, 0, 0, 31, 1, 0, 0, 0,
		0, 33, 1, 0, 0, 0, 0, 35, 1, 0, 0, 0, 0, 37, 1, 0, 0, 0, 0, 39, 1, 0, 0,
		0, 0, 41, 1, 0, 0, 0, 0, 43, 1, 0, 0, 0, 0, 45, 1, 0, 0, 0, 0, 47, 1, 0,
		0, 0, 0, 49, 1, 0, 0, 0, 0, 51, 1, 0, 0, 0, 0, 53, 1, 0, 0, 0, 0, 55, 1,
		0, 0, 0, 0, 57, 1, 0, 0, 0, 1, 59, 1, 0, 0, 0, 3, 66, 1, 0, 0, 0, 5, 71,
		1, 0, 0, 0, 7, 77, 1, 0, 0, 0, 9, 81, 1, 0, 0, 0, 11, 84, 1, 0, 0, 0, 13,
		87, 1, 0, 0, 0, 15, 89, 1, 0, 0, 0, 17, 91, 1, 0, 0, 0, 19, 94, 1, 0, 0,
		0, 21, 97, 1, 0, 0, 0, 23, 99, 1, 0, 0, 0, 25, 101, 1, 0, 0, 0, 27, 104,
		1, 0, 0, 0, 29, 107, 1, 0, 0, 0, 31, 110, 1, 0, 0, 0, 33, 112, 1, 0, 0,
		0, 35, 114, 1, 0, 0, 0, 37, 116, 1, 0, 0, 0, 39, 120, 1, 0, 0, 0, 41, 124,
		1, 0, 0, 0, 43, 127, 1, 0, 0, 0, 45, 132, 1, 0, 0, 0, 47, 142, 1, 0, 0,
		0, 49, 147, 1, 0, 0, 0, 51, 153, 1, 0, 0, 0, 53, 161, 1, 0, 0, 0, 55, 173,
		1, 0, 0, 0, 57, 185, 1, 0, 0, 0, 59, 60, 5, 83, 0, 0, 60, 61, 5, 69, 0,
		0, 61, 62, 5, 76, 0, 0, 62, 63, 5, 69, 0, 0, 63, 64, 5, 67, 0, 0, 64, 65,
		5, 84, 0, 0, 65, 2, 1, 0, 0, 0, 66, 67, 5, 70, 0, 0, 67, 68, 5, 82, 0,
		0, 68, 69, 5, 79, 0, 0, 69, 70, 5, 77, 0, 0, 70, 4, 1, 0, 0, 0, 71, 72,
		5, 87, 0, 0, 72, 73, 5, 72, 0, 0, 73, 74, 5, 69, 0, 0, 74, 75, 5, 82, 0,
		0, 75, 76, 5, 69, 0, 0, 76, 6, 1, 0, 0, 0, 77, 78, 5, 65, 0, 0, 78, 79,
		5, 78, 0, 0, 79, 80, 5, 68, 0, 0, 80, 8, 1, 0, 0, 0, 81, 82, 5, 79, 0,
		0, 82, 83, 5, 82, 0, 0, 83, 10, 1, 0, 0, 0, 84, 85, 5, 73, 0, 0, 85, 86,
		5, 78, 0, 0, 86, 12, 1, 0, 0, 0, 87, 88, 5, 42, 0, 0, 88, 14, 1, 0, 0,
		0, 89, 90, 5, 46, 0, 0, 90, 16, 1, 0, 0, 0, 91, 92, 5, 61, 0, 0, 92, 93,
		5, 61, 0, 0, 93, 18, 1, 0, 0, 0, 94, 95, 5, 33, 0, 0, 95, 96, 5, 61, 0,
		0, 96, 20, 1, 0, 0, 0, 97, 98, 5, 62, 0, 0, 98, 22, 1, 0, 0, 0, 99, 100,
		5, 60, 0, 0, 100, 24, 1, 0, 0, 0, 101, 102, 5, 62, 0, 0, 102, 103, 5, 61,
		0, 0, 103, 26, 1, 0, 0, 0, 104, 105, 5, 60, 0, 0, 105, 106, 5, 61, 0, 0,
		106, 28, 1, 0, 0, 0, 107, 108, 5, 126, 0, 0, 108, 109, 5, 61, 0, 0, 109,
		30, 1, 0, 0, 0, 110, 111, 5, 40, 0, 0, 111, 32, 1, 0, 0, 0, 112, 113, 5,
		41, 0, 0, 113, 34, 1, 0, 0, 0, 114, 115, 5, 44, 0, 0, 115, 36, 1, 0, 0,
		0, 116, 117, 5, 111, 0, 0, 117, 118, 5, 98, 0, 0, 118, 119, 5, 106, 0,
		0, 119, 38, 1, 0, 0, 0, 120, 121, 5, 108, 0, 0, 121, 122, 5, 111, 0, 0,
		122, 123, 5, 103, 0, 0, 123, 40, 1, 0, 0, 0, 124, 125, 5, 73, 0, 0, 125,
		126, 5, 68, 0, 0, 126, 42, 1, 0, 0, 0, 127, 128, 5, 78, 0, 0, 128, 129,
		5, 97, 0, 0, 129, 130, 5, 109, 0, 0, 130, 131, 5, 101, 0, 0, 131, 44, 1,
		0, 0, 0, 132, 133, 5, 79, 0, 0, 133, 134, 5, 112, 0, 0, 134, 135, 5, 101,
		0, 0, 135, 136, 5, 114, 0, 0, 136, 137, 5, 97, 0, 0, 137, 138, 5, 116,
		0, 0, 138, 139, 5, 105, 0, 0, 139, 140, 5, 111, 0, 0, 140, 141, 5, 110,
		0, 0, 141, 46, 1, 0, 0, 0, 142, 143, 5, 80, 0, 0, 143, 144, 5, 97, 0, 0,
		144, 145, 5, 116, 0, 0, 145, 146, 5, 104, 0, 0, 146, 48, 1, 0, 0, 0, 147,
		148, 5, 86, 0, 0, 148, 149, 5, 97, 0, 0, 149, 150, 5, 108, 0, 0, 150, 151,
		5, 117, 0, 0, 151, 152, 5, 101, 0, 0, 152, 50, 1, 0, 0, 0, 153, 157, 7,
		0, 0, 0, 154, 156, 7, 1, 0, 0, 155, 154, 1, 0, 0, 0, 156, 159, 1, 0, 0,
		0, 157, 155, 1, 0, 0, 0, 157, 158, 1, 0, 0, 0, 158, 52, 1, 0, 0, 0, 159,
		157, 1, 0, 0, 0, 160, 162, 7, 2, 0, 0, 161, 160, 1, 0, 0, 0, 162, 163,
		1, 0, 0, 0, 163, 161, 1, 0, 0, 0, 163, 164, 1, 0, 0, 0, 164, 171, 1, 0,
		0, 0, 165, 167, 5, 46, 0, 0, 166, 168, 7, 2, 0, 0, 167, 166, 1, 0, 0, 0,
		168, 169, 1, 0, 0, 0, 169, 167, 1, 0, 0, 0, 169, 170, 1, 0, 0, 0, 170,
		172, 1, 0, 0, 0, 171, 165, 1, 0, 0, 0, 171, 172, 1, 0, 0, 0, 172, 54, 1,
		0, 0, 0, 173, 179, 5, 39, 0, 0, 174, 178, 8, 3, 0, 0, 175, 176, 5, 92,
		0, 0, 176, 178, 9, 0, 0, 0, 177, 174, 1, 0, 0, 0, 177, 175, 1, 0, 0, 0,
		178, 181, 1, 0, 0, 0, 179, 177, 1, 0, 0, 0, 179, 180, 1, 0, 0, 0, 180,
		182, 1, 0, 0, 0, 181, 179, 1, 0, 0, 0, 182, 183, 5, 39, 0, 0, 183, 56,
		1, 0, 0, 0, 184, 186, 7, 4, 0, 0, 185, 184, 1, 0, 0, 0, 186, 187, 1, 0,
		0, 0, 187, 185, 1, 0, 0, 0, 187, 188, 1, 0, 0, 0, 188, 189, 1, 0, 0, 0,
		189, 190, 6, 28, 0, 0, 190, 58, 1, 0, 0, 0, 8, 0, 157, 163, 169, 171, 177,
		179, 187, 1, 6, 0, 0,
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

// selectlangLexerInit initializes any static state used to implement selectlangLexer. By default the
// static state used to implement the lexer is lazily initialized during the first call to
// NewselectlangLexer(). You can call this function if you wish to initialize the static state ahead
// of time.
func SelectlangLexerInit() {
	staticData := &SelectlangLexerLexerStaticData
	staticData.once.Do(selectlanglexerLexerInit)
}

// NewselectlangLexer produces a new lexer instance for the optional input antlr.CharStream.
func NewselectlangLexer(input antlr.CharStream) *selectlangLexer {
	SelectlangLexerInit()
	l := new(selectlangLexer)
	l.BaseLexer = antlr.NewBaseLexer(input)
	staticData := &SelectlangLexerLexerStaticData
	l.Interpreter = antlr.NewLexerATNSimulator(l, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	l.channelNames = staticData.ChannelNames
	l.modeNames = staticData.ModeNames
	l.RuleNames = staticData.RuleNames
	l.LiteralNames = staticData.LiteralNames
	l.SymbolicNames = staticData.SymbolicNames
	l.GrammarFileName = "selectlang.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// selectlangLexer tokens.
const (
	selectlangLexerSELECT     = 1
	selectlangLexerFROM       = 2
	selectlangLexerWHERE      = 3
	selectlangLexerAND        = 4
	selectlangLexerOR         = 5
	selectlangLexerIN         = 6
	selectlangLexerSTAR       = 7
	selectlangLexerDOT        = 8
	selectlangLexerEQ         = 9
	selectlangLexerNE         = 10
	selectlangLexerGT         = 11
	selectlangLexerLT         = 12
	selectlangLexerGE         = 13
	selectlangLexerLE         = 14
	selectlangLexerREGEX_OP   = 15
	selectlangLexerLPAREN     = 16
	selectlangLexerRPAREN     = 17
	selectlangLexerCOMMA      = 18
	selectlangLexerOBJ        = 19
	selectlangLexerLOG        = 20
	selectlangLexerID_FIELD   = 21
	selectlangLexerNAME_FIELD = 22
	selectlangLexerOP_FIELD   = 23
	selectlangLexerPATH_FIELD = 24
	selectlangLexerVAL_FIELD  = 25
	selectlangLexerIDENTIFIER = 26
	selectlangLexerNUMBER     = 27
	selectlangLexerSTRING     = 28
	selectlangLexerWS         = 29
)
