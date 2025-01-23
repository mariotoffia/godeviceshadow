package mempersistence

import (
	"fmt"

	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
)

func renderSortKey(mt persistencemodel.ModelType, sk string) string {
	switch mt {
	case 0:
		return "DSC#" + sk
	case persistencemodel.ModelTypeReported:
		return "DSR#" + sk
	case persistencemodel.ModelTypeDesired:
		return "DSD#" + sk
	}

	panic(fmt.Sprintf("unknown model type: %s", mt.String()))
}
