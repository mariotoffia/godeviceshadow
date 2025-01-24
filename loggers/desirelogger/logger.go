package desirelogger

import (
	"github.com/mariotoffia/godeviceshadow/model"
	"github.com/mariotoffia/godeviceshadow/utils/reutils"
)

// New Implements the `model.CreatableDesiredLogger` interface.
func (d *DesireLogger) New() model.DesiredLogger {
	return New()
}

// Acknowledge Implements the `model.DesiredLogger` interface and keeps a record on the acknowledged
// values.
func (d *DesireLogger) Acknowledge(path string, value model.ValueAndTimestamp) {
	d.acknowledged[path] = value
}

func (d *DesireLogger) Acknowledged() map[string]model.ValueAndTimestamp {
	return d.acknowledged
}

// FromPath accepts a _path_ regexp that selects on the path of the acknowledged values.
//
// The _path_ regex is caches so no additional compile cost, except the first time, is paid.
func (d *DesireLogger) FromPath(path string) (map[string]model.ValueAndTimestamp, error) {
	re, err := reutils.Shared.GetOrCompile(path)

	if err != nil {
		return nil, err
	}

	m := make(map[string]model.ValueAndTimestamp, len(d.acknowledged))

	for k, v := range d.acknowledged {
		if re.MatchString(k) {
			m[k] = v
		}
	}

	return m, nil
}
