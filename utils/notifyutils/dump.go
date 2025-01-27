package notifyutils

import (
	"io"

	"github.com/mariotoffia/godeviceshadow/model/notifiermodel"
)

func Dump(wr io.Writer, sel notifiermodel.Selection) {

	spaces := func(indent int) string {
		s := ""

		for i := 0; i < indent; i++ {
			s += "  "
		}

		return s
	}

	var dump func(sel notifiermodel.Selection, indent int)

	dump = func(sel notifiermodel.Selection, indent int) {
		io.WriteString(wr, spaces(indent))

		switch t := sel.(type) {
		case *notifiermodel.AndSelection:
			io.WriteString(wr, "AND\n")

			for _, s := range t.Selections {
				dump(s, indent+1)
			}

		case *notifiermodel.OrSelection:
			io.WriteString(wr, "OR\n")

			for _, s := range t.Selections {
				dump(s, indent+1)
			}

		case *notifiermodel.NotSelection:
			io.WriteString(wr, "NOT\n")

			dump(t.Negated, indent+1)

		default:
			io.WriteString(wr, "Selection\n")
		}
	}

	dump(sel, 0)
}
