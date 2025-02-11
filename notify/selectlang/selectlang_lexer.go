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
		"T__0", "T__1", "T__2", "T__3", "T__4", "T__5", "T__6", "T__7", "T__8",
		"T__9", "T__10", "T__11", "WHERE", "LPAREN", "RPAREN", "EQ", "NE", "GT",
		"LT", "GE", "LE", "BEFORE", "AFTER", "ID", "AND", "OR", "NOT", "NAME",
		"OPERATION", "NUMBER", "STRING", "TIME", "REGEX", "WS",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 0, 34, 276, 6, -1, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2,
		4, 7, 4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2,
		10, 7, 10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15,
		7, 15, 2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 2, 19, 7, 19, 2, 20, 7,
		20, 2, 21, 7, 21, 2, 22, 7, 22, 2, 23, 7, 23, 2, 24, 7, 24, 2, 25, 7, 25,
		2, 26, 7, 26, 2, 27, 7, 27, 2, 28, 7, 28, 2, 29, 7, 29, 2, 30, 7, 30, 2,
		31, 7, 31, 2, 32, 7, 32, 2, 33, 7, 33, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 2, 1,
		2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 3, 1, 3, 1, 3, 1, 3, 1, 4, 1, 4, 1, 5, 1,
		5, 1, 6, 1, 6, 1, 6, 1, 6, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1,
		8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1,
		9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 10, 1, 10, 1, 10, 1, 10, 1, 10,
		1, 10, 1, 10, 1, 10, 1, 10, 1, 10, 1, 11, 1, 11, 1, 11, 1, 11, 1, 11, 1,
		11, 1, 12, 1, 12, 1, 12, 1, 12, 1, 12, 1, 12, 1, 13, 1, 13, 1, 14, 1, 14,
		1, 15, 1, 15, 1, 15, 1, 16, 1, 16, 1, 16, 1, 17, 1, 17, 1, 18, 1, 18, 1,
		19, 1, 19, 1, 19, 1, 20, 1, 20, 1, 20, 1, 21, 1, 21, 1, 21, 1, 21, 1, 21,
		1, 21, 1, 21, 1, 22, 1, 22, 1, 22, 1, 22, 1, 22, 1, 22, 1, 23, 1, 23, 1,
		23, 1, 23, 1, 24, 1, 24, 1, 24, 1, 24, 1, 25, 1, 25, 1, 25, 1, 26, 1, 26,
		1, 26, 1, 26, 1, 27, 1, 27, 1, 27, 1, 27, 1, 27, 1, 27, 1, 28, 1, 28, 1,
		28, 1, 28, 1, 28, 1, 28, 1, 28, 1, 28, 1, 28, 1, 28, 1, 28, 1, 29, 3, 29,
		218, 8, 29, 1, 29, 4, 29, 221, 8, 29, 11, 29, 12, 29, 222, 1, 29, 1, 29,
		4, 29, 227, 8, 29, 11, 29, 12, 29, 228, 3, 29, 231, 8, 29, 1, 30, 1, 30,
		5, 30, 235, 8, 30, 10, 30, 12, 30, 238, 9, 30, 1, 30, 1, 30, 1, 31, 1,
		31, 1, 31, 1, 31, 1, 31, 1, 31, 1, 31, 1, 31, 1, 31, 1, 31, 1, 31, 1, 31,
		1, 31, 1, 31, 1, 31, 1, 31, 1, 31, 1, 31, 1, 31, 1, 32, 1, 32, 5, 32, 263,
		8, 32, 10, 32, 12, 32, 266, 9, 32, 1, 32, 1, 32, 1, 33, 4, 33, 271, 8,
		33, 11, 33, 12, 33, 272, 1, 33, 1, 33, 2, 236, 264, 0, 34, 1, 1, 3, 2,
		5, 3, 7, 4, 9, 5, 11, 6, 13, 7, 15, 8, 17, 9, 19, 10, 21, 11, 23, 12, 25,
		13, 27, 14, 29, 15, 31, 16, 33, 17, 35, 18, 37, 19, 39, 20, 41, 21, 43,
		22, 45, 23, 47, 24, 49, 25, 51, 26, 53, 27, 55, 28, 57, 29, 59, 30, 61,
		31, 63, 32, 65, 33, 67, 34, 1, 0, 2, 1, 0, 48, 57, 3, 0, 9, 10, 13, 13,
		32, 32, 282, 0, 1, 1, 0, 0, 0, 0, 3, 1, 0, 0, 0, 0, 5, 1, 0, 0, 0, 0, 7,
		1, 0, 0, 0, 0, 9, 1, 0, 0, 0, 0, 11, 1, 0, 0, 0, 0, 13, 1, 0, 0, 0, 0,
		15, 1, 0, 0, 0, 0, 17, 1, 0, 0, 0, 0, 19, 1, 0, 0, 0, 0, 21, 1, 0, 0, 0,
		0, 23, 1, 0, 0, 0, 0, 25, 1, 0, 0, 0, 0, 27, 1, 0, 0, 0, 0, 29, 1, 0, 0,
		0, 0, 31, 1, 0, 0, 0, 0, 33, 1, 0, 0, 0, 0, 35, 1, 0, 0, 0, 0, 37, 1, 0,
		0, 0, 0, 39, 1, 0, 0, 0, 0, 41, 1, 0, 0, 0, 0, 43, 1, 0, 0, 0, 0, 45, 1,
		0, 0, 0, 0, 47, 1, 0, 0, 0, 0, 49, 1, 0, 0, 0, 0, 51, 1, 0, 0, 0, 0, 53,
		1, 0, 0, 0, 0, 55, 1, 0, 0, 0, 0, 57, 1, 0, 0, 0, 0, 59, 1, 0, 0, 0, 0,
		61, 1, 0, 0, 0, 0, 63, 1, 0, 0, 0, 0, 65, 1, 0, 0, 0, 0, 67, 1, 0, 0, 0,
		1, 69, 1, 0, 0, 0, 3, 76, 1, 0, 0, 0, 5, 84, 1, 0, 0, 0, 7, 91, 1, 0, 0,
		0, 9, 95, 1, 0, 0, 0, 11, 97, 1, 0, 0, 0, 13, 99, 1, 0, 0, 0, 15, 103,
		1, 0, 0, 0, 17, 110, 1, 0, 0, 0, 19, 117, 1, 0, 0, 0, 21, 129, 1, 0, 0,
		0, 23, 139, 1, 0, 0, 0, 25, 145, 1, 0, 0, 0, 27, 151, 1, 0, 0, 0, 29, 153,
		1, 0, 0, 0, 31, 155, 1, 0, 0, 0, 33, 158, 1, 0, 0, 0, 35, 161, 1, 0, 0,
		0, 37, 163, 1, 0, 0, 0, 39, 165, 1, 0, 0, 0, 41, 168, 1, 0, 0, 0, 43, 171,
		1, 0, 0, 0, 45, 178, 1, 0, 0, 0, 47, 184, 1, 0, 0, 0, 49, 188, 1, 0, 0,
		0, 51, 192, 1, 0, 0, 0, 53, 195, 1, 0, 0, 0, 55, 199, 1, 0, 0, 0, 57, 205,
		1, 0, 0, 0, 59, 217, 1, 0, 0, 0, 61, 232, 1, 0, 0, 0, 63, 241, 1, 0, 0,
		0, 65, 260, 1, 0, 0, 0, 67, 270, 1, 0, 0, 0, 69, 70, 5, 114, 0, 0, 70,
		71, 5, 101, 0, 0, 71, 72, 5, 112, 0, 0, 72, 73, 5, 111, 0, 0, 73, 74, 5,
		114, 0, 0, 74, 75, 5, 116, 0, 0, 75, 2, 1, 0, 0, 0, 76, 77, 5, 100, 0,
		0, 77, 78, 5, 101, 0, 0, 78, 79, 5, 115, 0, 0, 79, 80, 5, 105, 0, 0, 80,
		81, 5, 114, 0, 0, 81, 82, 5, 101, 0, 0, 82, 83, 5, 100, 0, 0, 83, 4, 1,
		0, 0, 0, 84, 85, 5, 100, 0, 0, 85, 86, 5, 101, 0, 0, 86, 87, 5, 108, 0,
		0, 87, 88, 5, 101, 0, 0, 88, 89, 5, 116, 0, 0, 89, 90, 5, 101, 0, 0, 90,
		6, 1, 0, 0, 0, 91, 92, 5, 97, 0, 0, 92, 93, 5, 108, 0, 0, 93, 94, 5, 108,
		0, 0, 94, 8, 1, 0, 0, 0, 95, 96, 5, 44, 0, 0, 96, 10, 1, 0, 0, 0, 97, 98,
		5, 58, 0, 0, 98, 12, 1, 0, 0, 0, 99, 100, 5, 97, 0, 0, 100, 101, 5, 100,
		0, 0, 101, 102, 5, 100, 0, 0, 102, 14, 1, 0, 0, 0, 103, 104, 5, 114, 0,
		0, 104, 105, 5, 101, 0, 0, 105, 106, 5, 109, 0, 0, 106, 107, 5, 111, 0,
		0, 107, 108, 5, 118, 0, 0, 108, 109, 5, 101, 0, 0, 109, 16, 1, 0, 0, 0,
		110, 111, 5, 117, 0, 0, 111, 112, 5, 112, 0, 0, 112, 113, 5, 100, 0, 0,
		113, 114, 5, 97, 0, 0, 114, 115, 5, 116, 0, 0, 115, 116, 5, 101, 0, 0,
		116, 18, 1, 0, 0, 0, 117, 118, 5, 97, 0, 0, 118, 119, 5, 99, 0, 0, 119,
		120, 5, 107, 0, 0, 120, 121, 5, 110, 0, 0, 121, 122, 5, 111, 0, 0, 122,
		123, 5, 119, 0, 0, 123, 124, 5, 108, 0, 0, 124, 125, 5, 101, 0, 0, 125,
		126, 5, 100, 0, 0, 126, 127, 5, 103, 0, 0, 127, 128, 5, 101, 0, 0, 128,
		20, 1, 0, 0, 0, 129, 130, 5, 110, 0, 0, 130, 131, 5, 111, 0, 0, 131, 132,
		5, 45, 0, 0, 132, 133, 5, 99, 0, 0, 133, 134, 5, 104, 0, 0, 134, 135, 5,
		97, 0, 0, 135, 136, 5, 110, 0, 0, 136, 137, 5, 103, 0, 0, 137, 138, 5,
		101, 0, 0, 138, 22, 1, 0, 0, 0, 139, 140, 5, 118, 0, 0, 140, 141, 5, 97,
		0, 0, 141, 142, 5, 108, 0, 0, 142, 143, 5, 117, 0, 0, 143, 144, 5, 101,
		0, 0, 144, 24, 1, 0, 0, 0, 145, 146, 5, 87, 0, 0, 146, 147, 5, 72, 0, 0,
		147, 148, 5, 69, 0, 0, 148, 149, 5, 82, 0, 0, 149, 150, 5, 69, 0, 0, 150,
		26, 1, 0, 0, 0, 151, 152, 5, 40, 0, 0, 152, 28, 1, 0, 0, 0, 153, 154, 5,
		41, 0, 0, 154, 30, 1, 0, 0, 0, 155, 156, 5, 61, 0, 0, 156, 157, 5, 61,
		0, 0, 157, 32, 1, 0, 0, 0, 158, 159, 5, 33, 0, 0, 159, 160, 5, 61, 0, 0,
		160, 34, 1, 0, 0, 0, 161, 162, 5, 62, 0, 0, 162, 36, 1, 0, 0, 0, 163, 164,
		5, 60, 0, 0, 164, 38, 1, 0, 0, 0, 165, 166, 5, 62, 0, 0, 166, 167, 5, 61,
		0, 0, 167, 40, 1, 0, 0, 0, 168, 169, 5, 60, 0, 0, 169, 170, 5, 61, 0, 0,
		170, 42, 1, 0, 0, 0, 171, 172, 5, 98, 0, 0, 172, 173, 5, 101, 0, 0, 173,
		174, 5, 102, 0, 0, 174, 175, 5, 111, 0, 0, 175, 176, 5, 114, 0, 0, 176,
		177, 5, 101, 0, 0, 177, 44, 1, 0, 0, 0, 178, 179, 5, 97, 0, 0, 179, 180,
		5, 102, 0, 0, 180, 181, 5, 116, 0, 0, 181, 182, 5, 101, 0, 0, 182, 183,
		5, 114, 0, 0, 183, 46, 1, 0, 0, 0, 184, 185, 5, 105, 0, 0, 185, 186, 5,
		100, 0, 0, 186, 187, 5, 58, 0, 0, 187, 48, 1, 0, 0, 0, 188, 189, 5, 65,
		0, 0, 189, 190, 5, 78, 0, 0, 190, 191, 5, 68, 0, 0, 191, 50, 1, 0, 0, 0,
		192, 193, 5, 79, 0, 0, 193, 194, 5, 82, 0, 0, 194, 52, 1, 0, 0, 0, 195,
		196, 5, 78, 0, 0, 196, 197, 5, 79, 0, 0, 197, 198, 5, 84, 0, 0, 198, 54,
		1, 0, 0, 0, 199, 200, 5, 110, 0, 0, 200, 201, 5, 97, 0, 0, 201, 202, 5,
		109, 0, 0, 202, 203, 5, 101, 0, 0, 203, 204, 5, 58, 0, 0, 204, 56, 1, 0,
		0, 0, 205, 206, 5, 111, 0, 0, 206, 207, 5, 112, 0, 0, 207, 208, 5, 101,
		0, 0, 208, 209, 5, 114, 0, 0, 209, 210, 5, 97, 0, 0, 210, 211, 5, 116,
		0, 0, 211, 212, 5, 105, 0, 0, 212, 213, 5, 111, 0, 0, 213, 214, 5, 110,
		0, 0, 214, 215, 5, 58, 0, 0, 215, 58, 1, 0, 0, 0, 216, 218, 5, 45, 0, 0,
		217, 216, 1, 0, 0, 0, 217, 218, 1, 0, 0, 0, 218, 220, 1, 0, 0, 0, 219,
		221, 7, 0, 0, 0, 220, 219, 1, 0, 0, 0, 221, 222, 1, 0, 0, 0, 222, 220,
		1, 0, 0, 0, 222, 223, 1, 0, 0, 0, 223, 230, 1, 0, 0, 0, 224, 226, 5, 46,
		0, 0, 225, 227, 7, 0, 0, 0, 226, 225, 1, 0, 0, 0, 227, 228, 1, 0, 0, 0,
		228, 226, 1, 0, 0, 0, 228, 229, 1, 0, 0, 0, 229, 231, 1, 0, 0, 0, 230,
		224, 1, 0, 0, 0, 230, 231, 1, 0, 0, 0, 231, 60, 1, 0, 0, 0, 232, 236, 5,
		39, 0, 0, 233, 235, 9, 0, 0, 0, 234, 233, 1, 0, 0, 0, 235, 238, 1, 0, 0,
		0, 236, 237, 1, 0, 0, 0, 236, 234, 1, 0, 0, 0, 237, 239, 1, 0, 0, 0, 238,
		236, 1, 0, 0, 0, 239, 240, 5, 39, 0, 0, 240, 62, 1, 0, 0, 0, 241, 242,
		7, 0, 0, 0, 242, 243, 6, 31, 0, 0, 243, 244, 5, 45, 0, 0, 244, 245, 7,
		0, 0, 0, 245, 246, 6, 31, 1, 0, 246, 247, 5, 45, 0, 0, 247, 248, 7, 0,
		0, 0, 248, 249, 6, 31, 2, 0, 249, 250, 5, 84, 0, 0, 250, 251, 7, 0, 0,
		0, 251, 252, 6, 31, 3, 0, 252, 253, 5, 58, 0, 0, 253, 254, 7, 0, 0, 0,
		254, 255, 6, 31, 4, 0, 255, 256, 5, 58, 0, 0, 256, 257, 7, 0, 0, 0, 257,
		258, 6, 31, 5, 0, 258, 259, 5, 90, 0, 0, 259, 64, 1, 0, 0, 0, 260, 264,
		5, 47, 0, 0, 261, 263, 9, 0, 0, 0, 262, 261, 1, 0, 0, 0, 263, 266, 1, 0,
		0, 0, 264, 265, 1, 0, 0, 0, 264, 262, 1, 0, 0, 0, 265, 267, 1, 0, 0, 0,
		266, 264, 1, 0, 0, 0, 267, 268, 5, 47, 0, 0, 268, 66, 1, 0, 0, 0, 269,
		271, 7, 1, 0, 0, 270, 269, 1, 0, 0, 0, 271, 272, 1, 0, 0, 0, 272, 270,
		1, 0, 0, 0, 272, 273, 1, 0, 0, 0, 273, 274, 1, 0, 0, 0, 274, 275, 6, 33,
		6, 0, 275, 68, 1, 0, 0, 0, 8, 0, 217, 222, 228, 230, 236, 264, 272, 7,
		1, 31, 0, 1, 31, 1, 1, 31, 2, 1, 31, 3, 1, 31, 4, 1, 31, 5, 6, 0, 0,
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
	selectlangLexerT__0      = 1
	selectlangLexerT__1      = 2
	selectlangLexerT__2      = 3
	selectlangLexerT__3      = 4
	selectlangLexerT__4      = 5
	selectlangLexerT__5      = 6
	selectlangLexerT__6      = 7
	selectlangLexerT__7      = 8
	selectlangLexerT__8      = 9
	selectlangLexerT__9      = 10
	selectlangLexerT__10     = 11
	selectlangLexerT__11     = 12
	selectlangLexerWHERE     = 13
	selectlangLexerLPAREN    = 14
	selectlangLexerRPAREN    = 15
	selectlangLexerEQ        = 16
	selectlangLexerNE        = 17
	selectlangLexerGT        = 18
	selectlangLexerLT        = 19
	selectlangLexerGE        = 20
	selectlangLexerLE        = 21
	selectlangLexerBEFORE    = 22
	selectlangLexerAFTER     = 23
	selectlangLexerID        = 24
	selectlangLexerAND       = 25
	selectlangLexerOR        = 26
	selectlangLexerNOT       = 27
	selectlangLexerNAME      = 28
	selectlangLexerOPERATION = 29
	selectlangLexerNUMBER    = 30
	selectlangLexerSTRING    = 31
	selectlangLexerTIME      = 32
	selectlangLexerREGEX     = 33
	selectlangLexerWS        = 34
)