package str_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/mariotoffia/godeviceshadow/loggers/str"
	"github.com/mariotoffia/godeviceshadow/merge"
	"github.com/mariotoffia/godeviceshadow/model"
	"github.com/stretchr/testify/require"
)

type Value struct {
	Value     any       `json:"v"`
	UpdatedAt time.Time `json:"ts"`
}

type Shadow map[string]Value

func (tmv *Value) GetTimestamp() time.Time {
	return tmv.UpdatedAt
}
func (tmv *Value) SetTimestamp(t time.Time) {
	tmv.UpdatedAt = t
}
func (tmv *Value) GetValue() any {
	return tmv.Value
}

func TestStringLogger(t *testing.T) {
	sl := str.NewStringLogger()

	oldShadow := Shadow{
		"test": {
			Value:     22.16,
			UpdatedAt: time.Now().UTC().Add(-1 * time.Hour),
		},
		"test2": {
			Value:     true,
			UpdatedAt: time.Now().UTC().Add(-1 * time.Hour),
		},
	}

	newShadow := Shadow{
		"test": {
			Value:     "to-string",
			UpdatedAt: time.Now().UTC(),
		},
		"test2": {
			Value:     false,
			UpdatedAt: time.Now().UTC(),
		},
		"test3": {
			Value:     21.16,
			UpdatedAt: time.Now().UTC(),
		},
	}

	_, err := merge.Merge(oldShadow, newShadow, merge.MergeOptions{
		Mode:    merge.ClientIsMaster,
		Loggers: []model.MergeLogger{sl},
	})

	require.NoError(t, err)
	fmt.Println(sl.String())
}
