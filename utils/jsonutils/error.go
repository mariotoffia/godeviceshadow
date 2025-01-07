package jsonutils

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// HighlightSyntaxError provides context around the error offset in the data.
//
// If no _ctxSize_ is provided, it will default to _40_.
func HighlightSyntaxError(data []byte, se *json.SyntaxError, ctxSize ...int) string {
	var contextSize int

	if len(ctxSize) > 0 {
		contextSize = ctxSize[0]
	}

	if contextSize <= 0 {
		contextSize = 40
	}

	start := int(se.Offset) - contextSize
	if start < 0 {
		start = 0
	}

	end := int(se.Offset) + contextSize
	if end > len(data) {
		end = len(data)
	}

	pre := data[start:int(se.Offset)]
	post := data[int(se.Offset):end]

	arrowLine := bytes.Repeat([]byte(" "), len(pre))
	arrowLine = append(arrowLine, []byte("^ "+se.Error())...)

	return fmt.Sprintf(
		"Error near offset %d:\n%s\n%s", se.Offset, string(pre)+string(post), string(arrowLine),
	)
}
