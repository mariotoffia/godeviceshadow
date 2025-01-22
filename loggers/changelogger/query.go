package changelogger

import (
	"github.com/mariotoffia/godeviceshadow/model"
	"github.com/mariotoffia/godeviceshadow/utils/reutils"
)

// ManagedFromPath will get all entries from managed log entries that has the specified path. The path is
// always a regexp pattern that will be evaluated.
//
// NOTE: Since the path is cached, this function can be invoked many times with the same path without paying the
// cost of compiling the regexp each time.
//
// If no operation is specified, all operations will be included.
func (cl *ChangeMergeLogger) ManagedFromPath(path string, operation ...model.MergeOperation) (ManagedLogMap, error) {
	re, err := reutils.Shared.GetOrCompile(path)

	if err != nil {
		return nil, err
	}

	lm := make(ManagedLogMap, len(cl.ManagedLog))

	for op, v := range cl.ManagedLog {
		if len(v) == 0 || (len(operation) > 0 && !op.In(operation...)) {
			continue
		}

		values := make([]ManagedValue, 0, len(v))

		for _, mv := range v {
			if re.MatchString(mv.Path) {
				values = append(values, mv)
			}
		}

		if len(values) > 0 {
			lm[op] = values
		}
	}

	return lm, nil
}

// PlainFromPath will get all entries from plain log entries that has the specified path. The path is
// always a regexp pattern that will be evaluated.
//
// NOTE: Since the path is cached, this function can be invoked many times with the same path without paying the
// cost of compiling the regexp each time.
func (cl *ChangeMergeLogger) PlainFromPath(path string, operation ...model.MergeOperation) (PlainLogMap, error) {
	re, err := reutils.Shared.GetOrCompile(path)

	if err != nil {
		return nil, err
	}

	lm := make(PlainLogMap, len(cl.PlainLog))

	for op, v := range cl.PlainLog {
		if len(v) == 0 || (len(operation) > 0 && !op.In(operation...)) {
			continue
		}

		values := make([]PlainValue, 0, len(v))

		for _, pv := range v {
			if re.MatchString(pv.Path) {
				values = append(values, pv)
			}
		}

		if len(values) > 0 {
			lm[op] = values
		}
	}

	return lm, nil
}
