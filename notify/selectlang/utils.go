package selectlang

import (
	"slices"

	"github.com/antlr4-go/antlr/v4"
)

// FirstChild returns the first available child that implements the `T` interface.
//
// If none are found it will return the zero value of `T`.
func FirstChild[T any](ctx antlr.ParserRuleContext) T {
	if ctx == nil {
		var zero T

		return zero
	}

	for i := 0; i < ctx.GetChildCount(); i++ {
		child := ctx.GetChild(i)
		if child != nil {
			if t, ok := child.(T); ok {
				return t
			}
		}
	}

	var zero T

	return zero
}

func ToStringList(ctx antlr.ParserRuleContext, filter ...string) []string {
	if ctx == nil {
		return nil
	}

	list := make([]string, 0, ctx.GetChildCount())

	for i := 0; i < ctx.GetChildCount(); i++ {
		child := ctx.GetChild(i)
		if child != nil {
			if pt, ok := child.(antlr.ParseTree); ok {
				txt := pt.GetText()

				if slices.Contains(filter, txt) {
					continue
				}

				list = append(list, txt)
			}
		}
	}

	return list
}
