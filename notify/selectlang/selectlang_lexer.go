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
		"", "'SELECT'", "'FROM'", "'WHERE'", "'AND'", "'OR'", "'IN'", "'HAS'",
		"'*'", "'.'", "'=='", "'!='", "'>'", "'<'", "'>='", "'<='", "'~='",
		"'~!='", "'('", "')'", "','", "'obj'", "'log'", "'ID'", "'Name'", "'Operation'",
		"'Path'", "'Value'",
	}
	staticData.SymbolicNames = []string{
		"", "SELECT", "FROM", "WHERE", "AND", "OR", "IN", "HAS", "STAR", "DOT",
		"EQ", "NE", "GT", "LT", "GE", "LE", "REGEX_OP", "REGEX_NOT_OP", "LPAREN",
		"RPAREN", "COMMA", "OBJ", "LOG", "ID_FIELD", "NAME_FIELD", "OP_FIELD",
		"PATH_FIELD", "VAL_FIELD", "IDENTIFIER", "NUMBER", "STRING", "WS",
	}
	staticData.RuleNames = []string{
		"SELECT", "FROM", "WHERE", "AND", "OR", "IN", "HAS", "STAR", "DOT",
		"EQ", "NE", "GT", "LT", "GE", "LE", "REGEX_OP", "REGEX_NOT_OP", "LPAREN",
		"RPAREN", "COMMA", "OBJ", "LOG", "ID_FIELD", "NAME_FIELD", "OP_FIELD",
		"PATH_FIELD", "VAL_FIELD", "IDENTIFIER", "NUMBER", "STRING", "WS",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 0, 31, 203, 6, -1, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2,
		4, 7, 4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2,
		10, 7, 10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15,
		7, 15, 2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 2, 19, 7, 19, 2, 20, 7,
		20, 2, 21, 7, 21, 2, 22, 7, 22, 2, 23, 7, 23, 2, 24, 7, 24, 2, 25, 7, 25,
		2, 26, 7, 26, 2, 27, 7, 27, 2, 28, 7, 28, 2, 29, 7, 29, 2, 30, 7, 30, 1,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 3, 1, 3, 1, 3, 1, 3, 1, 4, 1, 4, 1,
		4, 1, 5, 1, 5, 1, 5, 1, 6, 1, 6, 1, 6, 1, 6, 1, 7, 1, 7, 1, 8, 1, 8, 1,
		9, 1, 9, 1, 9, 1, 10, 1, 10, 1, 10, 1, 11, 1, 11, 1, 12, 1, 12, 1, 13,
		1, 13, 1, 13, 1, 14, 1, 14, 1, 14, 1, 15, 1, 15, 1, 15, 1, 16, 1, 16, 1,
		16, 1, 16, 1, 17, 1, 17, 1, 18, 1, 18, 1, 19, 1, 19, 1, 20, 1, 20, 1, 20,
		1, 20, 1, 21, 1, 21, 1, 21, 1, 21, 1, 22, 1, 22, 1, 22, 1, 23, 1, 23, 1,
		23, 1, 23, 1, 23, 1, 24, 1, 24, 1, 24, 1, 24, 1, 24, 1, 24, 1, 24, 1, 24,
		1, 24, 1, 24, 1, 25, 1, 25, 1, 25, 1, 25, 1, 25, 1, 26, 1, 26, 1, 26, 1,
		26, 1, 26, 1, 26, 1, 27, 1, 27, 5, 27, 168, 8, 27, 10, 27, 12, 27, 171,
		9, 27, 1, 28, 4, 28, 174, 8, 28, 11, 28, 12, 28, 175, 1, 28, 1, 28, 4,
		28, 180, 8, 28, 11, 28, 12, 28, 181, 3, 28, 184, 8, 28, 1, 29, 1, 29, 1,
		29, 1, 29, 5, 29, 190, 8, 29, 10, 29, 12, 29, 193, 9, 29, 1, 29, 1, 29,
		1, 30, 4, 30, 198, 8, 30, 11, 30, 12, 30, 199, 1, 30, 1, 30, 0, 0, 31,
		1, 1, 3, 2, 5, 3, 7, 4, 9, 5, 11, 6, 13, 7, 15, 8, 17, 9, 19, 10, 21, 11,
		23, 12, 25, 13, 27, 14, 29, 15, 31, 16, 33, 17, 35, 18, 37, 19, 39, 20,
		41, 21, 43, 22, 45, 23, 47, 24, 49, 25, 51, 26, 53, 27, 55, 28, 57, 29,
		59, 30, 61, 31, 1, 0, 5, 3, 0, 65, 90, 95, 95, 97, 122, 4, 0, 48, 57, 65,
		90, 95, 95, 97, 122, 1, 0, 48, 57, 2, 0, 39, 39, 92, 92, 3, 0, 9, 10, 13,
		13, 32, 32, 209, 0, 1, 1, 0, 0, 0, 0, 3, 1, 0, 0, 0, 0, 5, 1, 0, 0, 0,
		0, 7, 1, 0, 0, 0, 0, 9, 1, 0, 0, 0, 0, 11, 1, 0, 0, 0, 0, 13, 1, 0, 0,
		0, 0, 15, 1, 0, 0, 0, 0, 17, 1, 0, 0, 0, 0, 19, 1, 0, 0, 0, 0, 21, 1, 0,
		0, 0, 0, 23, 1, 0, 0, 0, 0, 25, 1, 0, 0, 0, 0, 27, 1, 0, 0, 0, 0, 29, 1,
		0, 0, 0, 0, 31, 1, 0, 0, 0, 0, 33, 1, 0, 0, 0, 0, 35, 1, 0, 0, 0, 0, 37,
		1, 0, 0, 0, 0, 39, 1, 0, 0, 0, 0, 41, 1, 0, 0, 0, 0, 43, 1, 0, 0, 0, 0,
		45, 1, 0, 0, 0, 0, 47, 1, 0, 0, 0, 0, 49, 1, 0, 0, 0, 0, 51, 1, 0, 0, 0,
		0, 53, 1, 0, 0, 0, 0, 55, 1, 0, 0, 0, 0, 57, 1, 0, 0, 0, 0, 59, 1, 0, 0,
		0, 0, 61, 1, 0, 0, 0, 1, 63, 1, 0, 0, 0, 3, 70, 1, 0, 0, 0, 5, 75, 1, 0,
		0, 0, 7, 81, 1, 0, 0, 0, 9, 85, 1, 0, 0, 0, 11, 88, 1, 0, 0, 0, 13, 91,
		1, 0, 0, 0, 15, 95, 1, 0, 0, 0, 17, 97, 1, 0, 0, 0, 19, 99, 1, 0, 0, 0,
		21, 102, 1, 0, 0, 0, 23, 105, 1, 0, 0, 0, 25, 107, 1, 0, 0, 0, 27, 109,
		1, 0, 0, 0, 29, 112, 1, 0, 0, 0, 31, 115, 1, 0, 0, 0, 33, 118, 1, 0, 0,
		0, 35, 122, 1, 0, 0, 0, 37, 124, 1, 0, 0, 0, 39, 126, 1, 0, 0, 0, 41, 128,
		1, 0, 0, 0, 43, 132, 1, 0, 0, 0, 45, 136, 1, 0, 0, 0, 47, 139, 1, 0, 0,
		0, 49, 144, 1, 0, 0, 0, 51, 154, 1, 0, 0, 0, 53, 159, 1, 0, 0, 0, 55, 165,
		1, 0, 0, 0, 57, 173, 1, 0, 0, 0, 59, 185, 1, 0, 0, 0, 61, 197, 1, 0, 0,
		0, 63, 64, 5, 83, 0, 0, 64, 65, 5, 69, 0, 0, 65, 66, 5, 76, 0, 0, 66, 67,
		5, 69, 0, 0, 67, 68, 5, 67, 0, 0, 68, 69, 5, 84, 0, 0, 69, 2, 1, 0, 0,
		0, 70, 71, 5, 70, 0, 0, 71, 72, 5, 82, 0, 0, 72, 73, 5, 79, 0, 0, 73, 74,
		5, 77, 0, 0, 74, 4, 1, 0, 0, 0, 75, 76, 5, 87, 0, 0, 76, 77, 5, 72, 0,
		0, 77, 78, 5, 69, 0, 0, 78, 79, 5, 82, 0, 0, 79, 80, 5, 69, 0, 0, 80, 6,
		1, 0, 0, 0, 81, 82, 5, 65, 0, 0, 82, 83, 5, 78, 0, 0, 83, 84, 5, 68, 0,
		0, 84, 8, 1, 0, 0, 0, 85, 86, 5, 79, 0, 0, 86, 87, 5, 82, 0, 0, 87, 10,
		1, 0, 0, 0, 88, 89, 5, 73, 0, 0, 89, 90, 5, 78, 0, 0, 90, 12, 1, 0, 0,
		0, 91, 92, 5, 72, 0, 0, 92, 93, 5, 65, 0, 0, 93, 94, 5, 83, 0, 0, 94, 14,
		1, 0, 0, 0, 95, 96, 5, 42, 0, 0, 96, 16, 1, 0, 0, 0, 97, 98, 5, 46, 0,
		0, 98, 18, 1, 0, 0, 0, 99, 100, 5, 61, 0, 0, 100, 101, 5, 61, 0, 0, 101,
		20, 1, 0, 0, 0, 102, 103, 5, 33, 0, 0, 103, 104, 5, 61, 0, 0, 104, 22,
		1, 0, 0, 0, 105, 106, 5, 62, 0, 0, 106, 24, 1, 0, 0, 0, 107, 108, 5, 60,
		0, 0, 108, 26, 1, 0, 0, 0, 109, 110, 5, 62, 0, 0, 110, 111, 5, 61, 0, 0,
		111, 28, 1, 0, 0, 0, 112, 113, 5, 60, 0, 0, 113, 114, 5, 61, 0, 0, 114,
		30, 1, 0, 0, 0, 115, 116, 5, 126, 0, 0, 116, 117, 5, 61, 0, 0, 117, 32,
		1, 0, 0, 0, 118, 119, 5, 126, 0, 0, 119, 120, 5, 33, 0, 0, 120, 121, 5,
		61, 0, 0, 121, 34, 1, 0, 0, 0, 122, 123, 5, 40, 0, 0, 123, 36, 1, 0, 0,
		0, 124, 125, 5, 41, 0, 0, 125, 38, 1, 0, 0, 0, 126, 127, 5, 44, 0, 0, 127,
		40, 1, 0, 0, 0, 128, 129, 5, 111, 0, 0, 129, 130, 5, 98, 0, 0, 130, 131,
		5, 106, 0, 0, 131, 42, 1, 0, 0, 0, 132, 133, 5, 108, 0, 0, 133, 134, 5,
		111, 0, 0, 134, 135, 5, 103, 0, 0, 135, 44, 1, 0, 0, 0, 136, 137, 5, 73,
		0, 0, 137, 138, 5, 68, 0, 0, 138, 46, 1, 0, 0, 0, 139, 140, 5, 78, 0, 0,
		140, 141, 5, 97, 0, 0, 141, 142, 5, 109, 0, 0, 142, 143, 5, 101, 0, 0,
		143, 48, 1, 0, 0, 0, 144, 145, 5, 79, 0, 0, 145, 146, 5, 112, 0, 0, 146,
		147, 5, 101, 0, 0, 147, 148, 5, 114, 0, 0, 148, 149, 5, 97, 0, 0, 149,
		150, 5, 116, 0, 0, 150, 151, 5, 105, 0, 0, 151, 152, 5, 111, 0, 0, 152,
		153, 5, 110, 0, 0, 153, 50, 1, 0, 0, 0, 154, 155, 5, 80, 0, 0, 155, 156,
		5, 97, 0, 0, 156, 157, 5, 116, 0, 0, 157, 158, 5, 104, 0, 0, 158, 52, 1,
		0, 0, 0, 159, 160, 5, 86, 0, 0, 160, 161, 5, 97, 0, 0, 161, 162, 5, 108,
		0, 0, 162, 163, 5, 117, 0, 0, 163, 164, 5, 101, 0, 0, 164, 54, 1, 0, 0,
		0, 165, 169, 7, 0, 0, 0, 166, 168, 7, 1, 0, 0, 167, 166, 1, 0, 0, 0, 168,
		171, 1, 0, 0, 0, 169, 167, 1, 0, 0, 0, 169, 170, 1, 0, 0, 0, 170, 56, 1,
		0, 0, 0, 171, 169, 1, 0, 0, 0, 172, 174, 7, 2, 0, 0, 173, 172, 1, 0, 0,
		0, 174, 175, 1, 0, 0, 0, 175, 173, 1, 0, 0, 0, 175, 176, 1, 0, 0, 0, 176,
		183, 1, 0, 0, 0, 177, 179, 5, 46, 0, 0, 178, 180, 7, 2, 0, 0, 179, 178,
		1, 0, 0, 0, 180, 181, 1, 0, 0, 0, 181, 179, 1, 0, 0, 0, 181, 182, 1, 0,
		0, 0, 182, 184, 1, 0, 0, 0, 183, 177, 1, 0, 0, 0, 183, 184, 1, 0, 0, 0,
		184, 58, 1, 0, 0, 0, 185, 191, 5, 39, 0, 0, 186, 190, 8, 3, 0, 0, 187,
		188, 5, 92, 0, 0, 188, 190, 9, 0, 0, 0, 189, 186, 1, 0, 0, 0, 189, 187,
		1, 0, 0, 0, 190, 193, 1, 0, 0, 0, 191, 189, 1, 0, 0, 0, 191, 192, 1, 0,
		0, 0, 192, 194, 1, 0, 0, 0, 193, 191, 1, 0, 0, 0, 194, 195, 5, 39, 0, 0,
		195, 60, 1, 0, 0, 0, 196, 198, 7, 4, 0, 0, 197, 196, 1, 0, 0, 0, 198, 199,
		1, 0, 0, 0, 199, 197, 1, 0, 0, 0, 199, 200, 1, 0, 0, 0, 200, 201, 1, 0,
		0, 0, 201, 202, 6, 30, 0, 0, 202, 62, 1, 0, 0, 0, 8, 0, 169, 175, 181,
		183, 189, 191, 199, 1, 6, 0, 0,
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
	selectlangLexerSELECT       = 1
	selectlangLexerFROM         = 2
	selectlangLexerWHERE        = 3
	selectlangLexerAND          = 4
	selectlangLexerOR           = 5
	selectlangLexerIN           = 6
	selectlangLexerHAS          = 7
	selectlangLexerSTAR         = 8
	selectlangLexerDOT          = 9
	selectlangLexerEQ           = 10
	selectlangLexerNE           = 11
	selectlangLexerGT           = 12
	selectlangLexerLT           = 13
	selectlangLexerGE           = 14
	selectlangLexerLE           = 15
	selectlangLexerREGEX_OP     = 16
	selectlangLexerREGEX_NOT_OP = 17
	selectlangLexerLPAREN       = 18
	selectlangLexerRPAREN       = 19
	selectlangLexerCOMMA        = 20
	selectlangLexerOBJ          = 21
	selectlangLexerLOG          = 22
	selectlangLexerID_FIELD     = 23
	selectlangLexerNAME_FIELD   = 24
	selectlangLexerOP_FIELD     = 25
	selectlangLexerPATH_FIELD   = 26
	selectlangLexerVAL_FIELD    = 27
	selectlangLexerIDENTIFIER   = 28
	selectlangLexerNUMBER       = 29
	selectlangLexerSTRING       = 30
	selectlangLexerWS           = 31
)
