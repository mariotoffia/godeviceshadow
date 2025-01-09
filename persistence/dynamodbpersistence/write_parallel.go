package dynamodbpersistence

import (
	"context"
	"sync"

	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
	"github.com/mariotoffia/godeviceshadow/utils"
	"github.com/mariotoffia/godeviceshadow/utils/persistutils"
)

func (p *Persistence) writeParallel(ctx context.Context, opt persistencemodel.WriteOptions, groups []persistutils.GroupedWriteOperation, maxParallelism int) []persistencemodel.WriteResult {
	batches := utils.ToBatch(groups, maxParallelism)
	res := make([]persistencemodel.WriteResult, 0, len(groups)*2)
	mtx := sync.Mutex{}
	wg := sync.WaitGroup{}

	for i := range batches {
		for _, current := range batches[i] {
			wg.Add(1)

			go func(op persistutils.GroupedWriteOperation) {
				var wop []persistencemodel.WriteResult

				if op.Error != nil {
					for j := range op.Operations {
						wop = append(wop, persistencemodel.WriteResult{
							ID:      op.Operations[j].ID,
							Version: op.Operations[j].Version,
							Error:   op.Error,
						})
					}
				} else {
					wop = p.WriteOperationGroup(ctx, opt, op)
				}

				mtx.Lock()
				defer mtx.Unlock()

				res = append(res, wop...)
				wg.Done()
			}(current)
		}

		wg.Wait()
	}

	return res
}
