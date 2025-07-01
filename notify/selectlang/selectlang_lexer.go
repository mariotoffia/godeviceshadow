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
		"'('", "')'", "','", "'obj'", "'log'", "'ID'", "'Name'", "'Operation'",
		"'Path'", "'Value'",
	}
	staticData.SymbolicNames = []string{
		"", "SELECT", "FROM", "WHERE", "AND", "OR", "IN", "HAS", "STAR", "DOT",
		"EQ", "NE", "GT", "LT", "GE", "LE", "REGEX_OP", "LPAREN", "RPAREN",
		"COMMA", "OBJ", "LOG", "ID_FIELD", "NAME_FIELD", "OP_FIELD", "PATH_FIELD",
		"VAL_FIELD", "IDENTIFIER", "NUMBER", "STRING", "WS",
	}
	staticData.RuleNames = []string{
		"SELECT", "FROM", "WHERE", "AND", "OR", "IN", "HAS", "STAR", "DOT",
		"EQ", "NE", "GT", "LT", "GE", "LE", "REGEX_OP", "LPAREN", "RPAREN",
		"COMMA", "OBJ", "LOG", "ID_FIELD", "NAME_FIELD", "OP_FIELD", "PATH_FIELD",
		"VAL_FIELD", "IDENTIFIER", "NUMBER", "STRING", "WS",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 0, 30, 197, 6, -1, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2,
		4, 7, 4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2,
		10, 7, 10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15,
		7, 15, 2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 2, 19, 7, 19, 2, 20, 7,
		20, 2, 21, 7, 21, 2, 22, 7, 22, 2, 23, 7, 23, 2, 24, 7, 24, 2, 25, 7, 25,
		2, 26, 7, 26, 2, 27, 7, 27, 2, 28, 7, 28, 2, 29, 7, 29, 1, 0, 1, 0, 1,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 2, 1,
		2, 1, 2, 1, 2, 1, 2, 1, 3, 1, 3, 1, 3, 1, 3, 1, 4, 1, 4, 1, 4, 1, 5, 1,
		5, 1, 5, 1, 6, 1, 6, 1, 6, 1, 6, 1, 7, 1, 7, 1, 8, 1, 8, 1, 9, 1, 9, 1,
		9, 1, 10, 1, 10, 1, 10, 1, 11, 1, 11, 1, 12, 1, 12, 1, 13, 1, 13, 1, 13,
		1, 14, 1, 14, 1, 14, 1, 15, 1, 15, 1, 15, 1, 16, 1, 16, 1, 17, 1, 17, 1,
		18, 1, 18, 1, 19, 1, 19, 1, 19, 1, 19, 1, 20, 1, 20, 1, 20, 1, 20, 1, 21,
		1, 21, 1, 21, 1, 22, 1, 22, 1, 22, 1, 22, 1, 22, 1, 23, 1, 23, 1, 23, 1,
		23, 1, 23, 1, 23, 1, 23, 1, 23, 1, 23, 1, 23, 1, 24, 1, 24, 1, 24, 1, 24,
		1, 24, 1, 25, 1, 25, 1, 25, 1, 25, 1, 25, 1, 25, 1, 26, 1, 26, 5, 26, 162,
		8, 26, 10, 26, 12, 26, 165, 9, 26, 1, 27, 4, 27, 168, 8, 27, 11, 27, 12,
		27, 169, 1, 27, 1, 27, 4, 27, 174, 8, 27, 11, 27, 12, 27, 175, 3, 27, 178,
		8, 27, 1, 28, 1, 28, 1, 28, 1, 28, 5, 28, 184, 8, 28, 10, 28, 12, 28, 187,
		9, 28, 1, 28, 1, 28, 1, 29, 4, 29, 192, 8, 29, 11, 29, 12, 29, 193, 1,
		29, 1, 29, 0, 0, 30, 1, 1, 3, 2, 5, 3, 7, 4, 9, 5, 11, 6, 13, 7, 15, 8,
		17, 9, 19, 10, 21, 11, 23, 12, 25, 13, 27, 14, 29, 15, 31, 16, 33, 17,
		35, 18, 37, 19, 39, 20, 41, 21, 43, 22, 45, 23, 47, 24, 49, 25, 51, 26,
		53, 27, 55, 28, 57, 29, 59, 30, 1, 0, 5, 3, 0, 65, 90, 95, 95, 97, 122,
		4, 0, 48, 57, 65, 90, 95, 95, 97, 122, 1, 0, 48, 57, 2, 0, 39, 39, 92,
		92, 3, 0, 9, 10, 13, 13, 32, 32, 203, 0, 1, 1, 0, 0, 0, 0, 3, 1, 0, 0,
		0, 0, 5, 1, 0, 0, 0, 0, 7, 1, 0, 0, 0, 0, 9, 1, 0, 0, 0, 0, 11, 1, 0, 0,
		0, 0, 13, 1, 0, 0, 0, 0, 15, 1, 0, 0, 0, 0, 17, 1, 0, 0, 0, 0, 19, 1, 0,
		0, 0, 0, 21, 1, 0, 0, 0, 0, 23, 1, 0, 0, 0, 0, 25, 1, 0, 0, 0, 0, 27, 1,
		0, 0, 0, 0, 29, 1, 0, 0, 0, 0, 31, 1, 0, 0, 0, 0, 33, 1, 0, 0, 0, 0, 35,
		1, 0, 0, 0, 0, 37, 1, 0, 0, 0, 0, 39, 1, 0, 0, 0, 0, 41, 1, 0, 0, 0, 0,
		43, 1, 0, 0, 0, 0, 45, 1, 0, 0, 0, 0, 47, 1, 0, 0, 0, 0, 49, 1, 0, 0, 0,
		0, 51, 1, 0, 0, 0, 0, 53, 1, 0, 0, 0, 0, 55, 1, 0, 0, 0, 0, 57, 1, 0, 0,
		0, 0, 59, 1, 0, 0, 0, 1, 61, 1, 0, 0, 0, 3, 68, 1, 0, 0, 0, 5, 73, 1, 0,
		0, 0, 7, 79, 1, 0, 0, 0, 9, 83, 1, 0, 0, 0, 11, 86, 1, 0, 0, 0, 13, 89,
		1, 0, 0, 0, 15, 93, 1, 0, 0, 0, 17, 95, 1, 0, 0, 0, 19, 97, 1, 0, 0, 0,
		21, 100, 1, 0, 0, 0, 23, 103, 1, 0, 0, 0, 25, 105, 1, 0, 0, 0, 27, 107,
		1, 0, 0, 0, 29, 110, 1, 0, 0, 0, 31, 113, 1, 0, 0, 0, 33, 116, 1, 0, 0,
		0, 35, 118, 1, 0, 0, 0, 37, 120, 1, 0, 0, 0, 39, 122, 1, 0, 0, 0, 41, 126,
		1, 0, 0, 0, 43, 130, 1, 0, 0, 0, 45, 133, 1, 0, 0, 0, 47, 138, 1, 0, 0,
		0, 49, 148, 1, 0, 0, 0, 51, 153, 1, 0, 0, 0, 53, 159, 1, 0, 0, 0, 55, 167,
		1, 0, 0, 0, 57, 179, 1, 0, 0, 0, 59, 191, 1, 0, 0, 0, 61, 62, 5, 83, 0,
		0, 62, 63, 5, 69, 0, 0, 63, 64, 5, 76, 0, 0, 64, 65, 5, 69, 0, 0, 65, 66,
		5, 67, 0, 0, 66, 67, 5, 84, 0, 0, 67, 2, 1, 0, 0, 0, 68, 69, 5, 70, 0,
		0, 69, 70, 5, 82, 0, 0, 70, 71, 5, 79, 0, 0, 71, 72, 5, 77, 0, 0, 72, 4,
		1, 0, 0, 0, 73, 74, 5, 87, 0, 0, 74, 75, 5, 72, 0, 0, 75, 76, 5, 69, 0,
		0, 76, 77, 5, 82, 0, 0, 77, 78, 5, 69, 0, 0, 78, 6, 1, 0, 0, 0, 79, 80,
		5, 65, 0, 0, 80, 81, 5, 78, 0, 0, 81, 82, 5, 68, 0, 0, 82, 8, 1, 0, 0,
		0, 83, 84, 5, 79, 0, 0, 84, 85, 5, 82, 0, 0, 85, 10, 1, 0, 0, 0, 86, 87,
		5, 73, 0, 0, 87, 88, 5, 78, 0, 0, 88, 12, 1, 0, 0, 0, 89, 90, 5, 72, 0,
		0, 90, 91, 5, 65, 0, 0, 91, 92, 5, 83, 0, 0, 92, 14, 1, 0, 0, 0, 93, 94,
		5, 42, 0, 0, 94, 16, 1, 0, 0, 0, 95, 96, 5, 46, 0, 0, 96, 18, 1, 0, 0,
		0, 97, 98, 5, 61, 0, 0, 98, 99, 5, 61, 0, 0, 99, 20, 1, 0, 0, 0, 100, 101,
		5, 33, 0, 0, 101, 102, 5, 61, 0, 0, 102, 22, 1, 0, 0, 0, 103, 104, 5, 62,
		0, 0, 104, 24, 1, 0, 0, 0, 105, 106, 5, 60, 0, 0, 106, 26, 1, 0, 0, 0,
		107, 108, 5, 62, 0, 0, 108, 109, 5, 61, 0, 0, 109, 28, 1, 0, 0, 0, 110,
		111, 5, 60, 0, 0, 111, 112, 5, 61, 0, 0, 112, 30, 1, 0, 0, 0, 113, 114,
		5, 126, 0, 0, 114, 115, 5, 61, 0, 0, 115, 32, 1, 0, 0, 0, 116, 117, 5,
		40, 0, 0, 117, 34, 1, 0, 0, 0, 118, 119, 5, 41, 0, 0, 119, 36, 1, 0, 0,
		0, 120, 121, 5, 44, 0, 0, 121, 38, 1, 0, 0, 0, 122, 123, 5, 111, 0, 0,
		123, 124, 5, 98, 0, 0, 124, 125, 5, 106, 0, 0, 125, 40, 1, 0, 0, 0, 126,
		127, 5, 108, 0, 0, 127, 128, 5, 111, 0, 0, 128, 129, 5, 103, 0, 0, 129,
		42, 1, 0, 0, 0, 130, 131, 5, 73, 0, 0, 131, 132, 5, 68, 0, 0, 132, 44,
		1, 0, 0, 0, 133, 134, 5, 78, 0, 0, 134, 135, 5, 97, 0, 0, 135, 136, 5,
		109, 0, 0, 136, 137, 5, 101, 0, 0, 137, 46, 1, 0, 0, 0, 138, 139, 5, 79,
		0, 0, 139, 140, 5, 112, 0, 0, 140, 141, 5, 101, 0, 0, 141, 142, 5, 114,
		0, 0, 142, 143, 5, 97, 0, 0, 143, 144, 5, 116, 0, 0, 144, 145, 5, 105,
		0, 0, 145, 146, 5, 111, 0, 0, 146, 147, 5, 110, 0, 0, 147, 48, 1, 0, 0,
		0, 148, 149, 5, 80, 0, 0, 149, 150, 5, 97, 0, 0, 150, 151, 5, 116, 0, 0,
		151, 152, 5, 104, 0, 0, 152, 50, 1, 0, 0, 0, 153, 154, 5, 86, 0, 0, 154,
		155, 5, 97, 0, 0, 155, 156, 5, 108, 0, 0, 156, 157, 5, 117, 0, 0, 157,
		158, 5, 101, 0, 0, 158, 52, 1, 0, 0, 0, 159, 163, 7, 0, 0, 0, 160, 162,
		7, 1, 0, 0, 161, 160, 1, 0, 0, 0, 162, 165, 1, 0, 0, 0, 163, 161, 1, 0,
		0, 0, 163, 164, 1, 0, 0, 0, 164, 54, 1, 0, 0, 0, 165, 163, 1, 0, 0, 0,
		166, 168, 7, 2, 0, 0, 167, 166, 1, 0, 0, 0, 168, 169, 1, 0, 0, 0, 169,
		167, 1, 0, 0, 0, 169, 170, 1, 0, 0, 0, 170, 177, 1, 0, 0, 0, 171, 173,
		5, 46, 0, 0, 172, 174, 7, 2, 0, 0, 173, 172, 1, 0, 0, 0, 174, 175, 1, 0,
		0, 0, 175, 173, 1, 0, 0, 0, 175, 176, 1, 0, 0, 0, 176, 178, 1, 0, 0, 0,
		177, 171, 1, 0, 0, 0, 177, 178, 1, 0, 0, 0, 178, 56, 1, 0, 0, 0, 179, 185,
		5, 39, 0, 0, 180, 184, 8, 3, 0, 0, 181, 182, 5, 92, 0, 0, 182, 184, 9,
		0, 0, 0, 183, 180, 1, 0, 0, 0, 183, 181, 1, 0, 0, 0, 184, 187, 1, 0, 0,
		0, 185, 183, 1, 0, 0, 0, 185, 186, 1, 0, 0, 0, 186, 188, 1, 0, 0, 0, 187,
		185, 1, 0, 0, 0, 188, 189, 5, 39, 0, 0, 189, 58, 1, 0, 0, 0, 190, 192,
		7, 4, 0, 0, 191, 190, 1, 0, 0, 0, 192, 193, 1, 0, 0, 0, 193, 191, 1, 0,
		0, 0, 193, 194, 1, 0, 0, 0, 194, 195, 1, 0, 0, 0, 195, 196, 6, 29, 0, 0,
		196, 60, 1, 0, 0, 0, 8, 0, 163, 169, 175, 177, 183, 185, 193, 1, 6, 0,
		0,
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
	selectlangLexerHAS        = 7
	selectlangLexerSTAR       = 8
	selectlangLexerDOT        = 9
	selectlangLexerEQ         = 10
	selectlangLexerNE         = 11
	selectlangLexerGT         = 12
	selectlangLexerLT         = 13
	selectlangLexerGE         = 14
	selectlangLexerLE         = 15
	selectlangLexerREGEX_OP   = 16
	selectlangLexerLPAREN     = 17
	selectlangLexerRPAREN     = 18
	selectlangLexerCOMMA      = 19
	selectlangLexerOBJ        = 20
	selectlangLexerLOG        = 21
	selectlangLexerID_FIELD   = 22
	selectlangLexerNAME_FIELD = 23
	selectlangLexerOP_FIELD   = 24
	selectlangLexerPATH_FIELD = 25
	selectlangLexerVAL_FIELD  = 26
	selectlangLexerIDENTIFIER = 27
	selectlangLexerNUMBER     = 28
	selectlangLexerSTRING     = 29
	selectlangLexerWS         = 30
)
