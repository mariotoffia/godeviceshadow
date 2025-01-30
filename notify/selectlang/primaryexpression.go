package selectlang

import (
	"slices"

	"github.com/mariotoffia/godeviceshadow/model/notifiermodel"
	"github.com/mariotoffia/godeviceshadow/utils/reutils"
)

type PrimaryExpression struct {
	ID        string
	Name      string
	Operation []string
}

// ToMatchFunc will return a function that can be used to match a notifier operation.
func (p *PrimaryExpression) ToMatchFunc() func(op notifiermodel.NotifierOperation, value bool) (bool, []notifiermodel.SelectedValue) {
	//
	return func(op notifiermodel.NotifierOperation, value bool) (bool, []notifiermodel.SelectedValue) {
		if p.ID != "" {
			if str, ok := AsString(p.ID); ok {
				if str != op.ID.ID {
					return false, nil
				}
			} else if re, ok := AsRegex(p.ID); ok {
				if regexp, err := reutils.Shared.GetOrCompile(re); err != nil {
					return false, nil
				} else if !regexp.MatchString(op.ID.ID) {
					return false, nil
				}
			}
		}

		if p.Name != "" {
			if str, ok := AsString(p.Name); ok {
				if str != op.ID.Name {
					return false, nil
				}
			} else if re, ok := AsRegex(p.Name); ok {
				if regexp, err := reutils.Shared.GetOrCompile(re); err != nil {
					return false, nil
				} else if !regexp.MatchString(op.ID.Name) {
					return false, nil
				}
			}
		}

		if len(p.Operation) > 0 {
			if p.Operation[0] != "all" {
				if !slices.Contains(p.Operation, string(op.Operation)) {
					return false, nil
				}
			}
		}

		return true, nil
	}
}
