package selectlang

import (
	"maps"
	"slices"
	"time"

	"github.com/mariotoffia/godeviceshadow/loggers/changelogger"
	"github.com/mariotoffia/godeviceshadow/model"
	"github.com/mariotoffia/godeviceshadow/model/notifiermodel"
	"github.com/mariotoffia/godeviceshadow/utils/reutils"
)

type LoggerExpression struct {
	CaptureOperations   []string
	CaptureRegex        string
	CaptureEqMapVarExpr string
	Where               *Constraint
}

type vts struct {
	model.ValueAndTimestamp
	v any
}

func (v *vts) GetValue() any {
	return v.v
}

var noTime = time.Time{}

func (v *vts) GetTimeStamp() time.Time {
	return noTime
}

func (le *LoggerExpression) ToMatchFunc() func(op notifiermodel.NotifierOperation, value bool) (bool, []notifiermodel.SelectedValue) {
	//
	return func(op notifiermodel.NotifierOperation, value bool) (bool, []notifiermodel.SelectedValue) {
		ml := make([]changelogger.ManagedValue, 0, op.MergeLogger.ManagedLog.Size())
		pl := make([]changelogger.PlainValue, 0, op.MergeLogger.PlainLog.Size())
		al := make([]changelogger.PlainValue, 0, len(op.DesireLogger.Acknowledged()))
		all := slices.Contains(le.CaptureOperations, "all")

		captureEqMapVarExpr, _ := AsString(le.CaptureEqMapVarExpr)
		captureRegex, _ := AsRegex(le.CaptureRegex)

		var (
			acknowledged []changelogger.PlainValue
		)

		convertDesired := func() []changelogger.PlainValue {
			ack := op.DesireLogger.Acknowledged()
			if len(ack) == 0 {
				return nil
			}

			res := make([]changelogger.PlainValue, 0, len(ack))

			for k, v := range ack {
				acknowledged = append(acknowledged, changelogger.PlainValue{
					Path:     k,
					NewValue: v,
				})
			}

			return res
		}

		matchManaged := func(mv changelogger.ManagedValue) (changelogger.ManagedValue, bool) {
			if captureRegex != "" {
				if re, err := reutils.Shared.GetOrCompile(captureRegex); (err == nil && !re.MatchString(mv.Path)) || err != nil {
					return mv, false
				}
			}

			if captureEqMapVarExpr == "" {
				return mv, true
			}

			var ve model.ValueAndTimestamp

			if mv.NewValue != nil {
				ve = mv.NewValue
			} else if mv.OldValue != nil {
				ve = mv.OldValue
			}

			if ve == nil {
				return mv, false
			}

			if m, ok := ve.GetValue().(map[string]any); ok {
				if v, ok := m[captureEqMapVarExpr]; ok {
					return changelogger.ManagedValue{
						NewValue: &model.ValueAndTimestampImpl{Value: v, Timestamp: ve.GetTimestamp()},
					}, true
				}
			}

			return mv, false
		}

		matchPlain := func(pv changelogger.PlainValue) bool {
			if captureRegex != "" {
				if re, err := reutils.Shared.GetOrCompile(captureRegex); (err == nil && !re.MatchString(pv.Path)) || err != nil {
					return false
				}
			}

			if captureEqMapVarExpr == "" {
				return true
			}

			if pv.NewValue != nil {
				if m, ok := pv.NewValue.(map[string]any); ok {
					if v, ok := m[captureEqMapVarExpr]; ok {
						return v == value
					}
				}
			}

			if pv.OldValue != nil {
				if m, ok := pv.OldValue.(map[string]any); ok {
					if v, ok := m[captureEqMapVarExpr]; ok {
						return v == value
					}
				}
			}

			return false
		}

		if all {
			for _, v := range op.MergeLogger.ManagedLog {
				for _, vv := range v {
					if vv, ok := matchManaged(vv); ok {
						ml = append(ml, vv)
					}
				}
			}

			for _, v := range op.MergeLogger.PlainLog {
				for _, vv := range v {
					if matchPlain(vv) {
						pl = append(pl, vv)
					}
				}
			}

			for _, ack := range convertDesired() {
				if matchPlain(ack) {
					al = append(al, ack)
				}
			}
		} else {
			for _, operation := range le.CaptureOperations {
				if m, ok := op.MergeLogger.ManagedLog[model.MergeOperationFromString(operation)]; ok {
					for _, vv := range m {
						if vv, ok := matchManaged(vv); ok {
							ml = append(ml, vv)
						}
					}
				}

				if p, ok := op.MergeLogger.PlainLog[model.MergeOperationFromString(operation)]; ok {
					for _, vv := range p {
						if matchPlain(vv) {
							pl = append(pl, vv)
						}
					}
				}
			}

			if slices.Contains(le.CaptureOperations, "acknowledge") {
				for _, ack := range convertDesired() {
					if matchPlain(ack) {
						al = append(al, ack)
					}
				}
			}
		}

		size := len(ml) + len(pl) + len(al)

		if size == 0 {
			return false, nil
		}

		if le.Where != nil {
			res, match := le.Where.match(ml, pl, al, &selectResult{res: make(map[string]notifiermodel.SelectedValue, len(ml)+len(pl)+len(al))})

			if match && value {
				return match, slices.Collect(maps.Values(res.res))
			}

			return match, nil
		}

		if !value {
			return true, nil
		}

		sv := make([]notifiermodel.SelectedValue, 0, size)

		for _, v := range ml {
			sv = append(sv, notifiermodel.SelectedValue{
				Path:         v.Path,
				OldValue:     v.OldValue,
				NewValue:     v.NewValue,
				OldTimeStamp: v.OldTimeStamp,
				NewTimeStamp: v.NewTimeStamp,
			})
		}

		for _, v := range pl {
			sv = append(sv, notifiermodel.SelectedValue{
				Path:     v.Path,
				OldValue: &vts{v: v.OldValue},
				NewValue: &vts{v: v.NewValue},
			})
		}

		for _, v := range al {
			sv = append(sv, notifiermodel.SelectedValue{
				Path:     v.Path,
				OldValue: &vts{v: v.OldValue},
				NewValue: &vts{v: v.NewValue},
			})
		}

		return true, sv
	}
}
