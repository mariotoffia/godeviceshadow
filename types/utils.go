package types

import (
	"reflect"

	"github.com/mariotoffia/godeviceshadow/model"
)

func toEntry(t any, name string, meta []map[string]string) model.TypeEntry {
	rt := reflect.TypeOf(t)

	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	if name == "" {
		name = rt.PkgPath() + "." + rt.Name()
	}

	te := model.TypeEntry{
		Model: rt,
		Name:  name,
	}

	if len(meta) > 0 {
		te.Meta = meta[0]
	}

	return te
}
