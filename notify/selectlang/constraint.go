package selectlang

import (
	"fmt"
	"strconv"

	"github.com/mariotoffia/godeviceshadow/loggers/changelogger"
	"github.com/mariotoffia/godeviceshadow/model/notifiermodel"
	"github.com/mariotoffia/godeviceshadow/utils/reutils"
)

type Constraint struct {
	Variable  string
	CompareOp string
	Value     any // string, number, or regex
	ValueType ConstrainValueType
	And       []*Constraint
	Or        []*Constraint
}

func (c Constraint) String() string {
	return fmt.Sprintf("%s %s %v", c.Variable, c.CompareOp, c.Value)
}

// IsScoped will return `true` if the constraint has it's own scope, i.e. it is within '(' and ')'.
func (c Constraint) IsScoped() bool {
	return len(c.And) > 0 || len(c.Or) > 0
}

// IsOnlyScope is when the constraint itself do not have any variable, compare operation or value and
// but have `And` or `Or` constraints. Therefore it is simply a scope and nothing else.
func (c Constraint) IsOnlyScope() bool {
	return !c.HasConstrainValues() && c.IsScoped()
}

// HasConstrainValues returns `true` if the constraint has any value constraints. Independent of
// `And` or `Or` constraints.
func (c Constraint) HasConstrainValues() bool {
	return c.Variable != "" && c.CompareOp != "" && c.Value != nil
}

type ConstraintLogicalOp int

const (
	// ConstraintLogicalLHS is no logical operation, instead it is the left hand side
	// of a logical operation.
	ConstraintLogicalLHS ConstraintLogicalOp = iota
	ConstraintLogicalOpAnd
	ConstraintLogicalOpOr
	// ConstraintLogicalOpNot can only be used by `Scope` and not `Constraint`
	ConstraintLogicalOpNot
)

type ConstrainValueType int

const (
	// ConstrainValueString is a plain string
	ConstrainValueString ConstrainValueType = iota
	// ConstrainValueNumber is a float64 number
	ConstrainValueNumber
	// is a string that represents a regex
	ConstrainValueRegex
)

func (cvt ConstrainValueType) String() string {
	switch cvt {
	case ConstrainValueString:
		return "string"
	case ConstrainValueNumber:
		return "number"
	case ConstrainValueRegex:
		return "regex"
	}

	return "unknown"
}

type selectResult struct {
	res map[string]notifiermodel.SelectedValue
}

func (c *Constraint) match(mv []changelogger.ManagedValue, pl, ack []changelogger.PlainValue, sr *selectResult) (*selectResult, bool) {
	var match bool

	if c.Value != nil {
		for _, v := range mv {
			vv := v.NewValue.GetValue()

			if vv == nil {
				vv = v.OldValue.GetValue()
			}

			if vv == nil {
				continue
			}

			if _, ok := vv.(map[string]any); ok {
				continue
			}

			switch c.ValueType {
			case ConstrainValueString:
				vvs, ok := anyToString(vv)

				if !ok {
					continue
				}

				switch c.CompareOp {
				case "==":
					if vvs == c.Value.(string) {
						match = true
					}
				case "!=":
					if vvs != c.Value.(string) {
						match = true
					}
				}
			case ConstrainValueRegex:
				vvs, ok := anyToString(vv)

				if !ok {
					continue
				}

				var rematch bool

				if regex, ok := AsRegex(c.Value.(string)); !ok {
					if re, err := reutils.Shared.GetOrCompile(regex); (err == nil && re.MatchString(vvs)) || err != nil {
						rematch = true
					}
				}

				switch c.CompareOp {
				case "==":
					if rematch {
						match = true
					}
				case "!=":
					if !rematch {
						match = true
					}
				}
			case ConstrainValueNumber:
				vvn, ok := anyToFloat64(vv)

				if !ok {
					continue
				}

				switch c.CompareOp {
				case "==":
					if vvn == c.Value.(float64) {
						match = true
					}
				case "!=":
					if vvn != c.Value.(float64) {
						match = true
					}
				case ">":
					if vvn > c.Value.(float64) {
						match = true
					}
				case "<":
					if vvn < c.Value.(float64) {
						match = true
					}
				case ">=":
					if vvn >= c.Value.(float64) {
						match = true
					}
				case "<=":
					if vvn <= c.Value.(float64) {
						match = true
					}
				}
			}

			if match {
				sr.res[v.Path] = notifiermodel.SelectedValue{
					Path:         v.Path,
					NewValue:     v.NewValue,
					OldValue:     v.OldValue,
					OldTimeStamp: v.OldTimeStamp,
					NewTimeStamp: v.NewTimeStamp,
				}
			}
		}
	}

	var am, om, m bool

	for _, or := range c.Or {
		if sr, m = or.match(mv, pl, ack, sr); m {
			om = true
		}
	}

	for _, and := range c.And {
		if sr, match = and.match(mv, pl, ack, sr); !match {
			am = false
			break
		} else {
			am = true
		}
	}

	return sr, am || om || match
}

func anyToString(v any) (string, bool) {
	if v == nil {
		return "", false
	}

	if s, ok := v.(string); ok {
		return s, true
	}

	if b, ok := v.([]byte); ok {
		return string(b), true
	}

	return fmt.Sprintf("%v", v), true
}
func anyToFloat64(v any) (float64, bool) {
	switch t := v.(type) {
	case float64:
		return t, true
	case int:
		return float64(t), true
	case int64:
		return float64(t), true
	case int32:
		return float64(t), true
	case float32:
		return float64(t), true
	case uint:
		return float64(t), true
	case uint64:
		return float64(t), true
	case uint32:
		return float64(t), true
	case uint8:
		return float64(t), true
	case uint16:
		return float64(t), true
	case int8:
		return float64(t), true
	case int16:
		return float64(t), true
	case string:
		if f, err := strconv.ParseFloat(t, 64); err == nil {
			return f, true
		}
	case []byte:
		if f, err := strconv.ParseFloat(string(t), 64); err == nil {
			return f, true
		}
	}

	return -1, false
}
