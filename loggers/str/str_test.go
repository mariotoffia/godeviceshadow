package str_test

import (
	"strings"
	"testing"
	"time"

	"github.com/mariotoffia/godeviceshadow/loggers/str"
	"github.com/mariotoffia/godeviceshadow/merge"
	"github.com/mariotoffia/godeviceshadow/model"
	"github.com/stretchr/testify/assert"
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
func (tmv *Value) GetValue() any {
	return tmv.Value
}

func TestStringLogger(t *testing.T) {
	sl := str.NewStringLogger()
	now := time.Now().UTC()
	nowMinusOneHour := now.Add(-1 * time.Hour)

	oldShadow := Shadow{
		"test": {
			Value:     22.16,
			UpdatedAt: nowMinusOneHour,
		},
		"test2": {
			Value:     true,
			UpdatedAt: nowMinusOneHour,
		},
	}

	newShadow := Shadow{
		"test": {
			Value:     "to-string",
			UpdatedAt: now,
		},
		"test2": {
			Value:     false,
			UpdatedAt: now,
		},
		"test3": {
			Value:     21.16,
			UpdatedAt: now,
		},
	}

	_, err := merge.Merge(oldShadow, newShadow, merge.MergeOptions{
		Mode:    merge.ClientIsMaster,
		Loggers: []model.MergeLogger{sl},
	})

	require.NoError(t, err)

	res := strings.Split(sl.String(), "\n")

	var lines []string

	for _, l := range res {
		if l != "" {
			lines = append(lines, l)
		}
	}

	require.Equal(t, 4, len(lines))

	items := func(s string) []string {
		var res []string

		for _, item := range strings.Split(s, " ") {
			if item != "" {
				res = append(res, item)
			}
		}

		return res
	}

	find := func(s string) string {
		for _, line := range lines {
			if strings.Contains(line, "Z "+s+" ") {
				return line
			}
		}

		require.Fail(t, "could not find line with prefix: "+s)
		return ""
	}

	s := items(lines[0])
	require.Equal(t, 8, len(s))
	assert.Equal(t, "Operation", s[0])
	assert.Equal(t, "Old Timestamp", s[1]+" "+s[2])
	assert.Equal(t, "New Timestamp", s[3]+" "+s[4])
	assert.Equal(t, "Path", s[5])
	assert.Equal(t, "OldValue", s[6])
	assert.Equal(t, "NewValue", s[7])

	s = items(find("test2"))
	require.Equal(t, 6, len(s))
	assert.Equal(t, "update", s[0])
	assert.Equal(t, nowMinusOneHour.Format(time.RFC3339), s[1])
	assert.Equal(t, now.Format(time.RFC3339), s[2])
	assert.Equal(t, "test2", s[3])
	assert.Equal(t, "true", s[4])
	assert.Equal(t, "false", s[5])

	s = items(find("test"))
	require.Equal(t, 6, len(s))
	assert.Equal(t, "update", s[0])
	assert.Equal(t, nowMinusOneHour.Format(time.RFC3339), s[1])
	assert.Equal(t, now.Format(time.RFC3339), s[2])
	assert.Equal(t, "test", s[3])
	assert.Equal(t, "22.16", s[4])
	assert.Equal(t, "to-string", s[5])

	s = items(find("test3"))
	require.Equal(t, 6, len(s))
	assert.Equal(t, "add", s[0])
	assert.Equal(t, "0001-01-01T00:00:00Z", s[1])
	assert.Equal(t, now.Format(time.RFC3339), s[2])
	assert.Equal(t, "test3", s[3])
	assert.Equal(t, "nil", s[4])
	assert.Equal(t, "21.16", s[5])
}
