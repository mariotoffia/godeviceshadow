package dynamodbpersistence

import (
	"context"
	"sync"

	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
	"github.com/mariotoffia/godeviceshadow/utils/persistutils"
)

type parallelWriteResult struct {
	index int
	data  []persistencemodel.WriteResult
}

func parallelWriteWorker(
	ctx context.Context,
	wg *sync.WaitGroup,
	queue <-chan int,
	results chan<- parallelWriteResult,
	p *Persistence,
	opt persistencemodel.WriteOptions,
	groups []persistutils.GroupedWriteOperation,
) {
	defer wg.Done()

	for {
		select {
		case idx, ok := <-queue:
			if !ok {
				return // Exit if the queue is closed
			}

			group := groups[idx]

			// Handle validation errors
			if group.Error != nil {
				errorResults := make([]persistencemodel.WriteResult, len(group.Operations))
				for i, op := range group.Operations {
					errorResults[i] = persistencemodel.WriteResult{
						ID:      op.ID,
						Version: op.Version,
						Error:   group.Error,
					}
				}
				select {
				case results <- parallelWriteResult{index: idx, data: errorResults}:
				case <-ctx.Done():
					return // Exit early if the context is cancelled
				}
				continue
			}

			// Process the group
			writeResults := p.WriteOperationGroup(ctx, opt, group)
			select {
			case results <- parallelWriteResult{index: idx, data: writeResults}:
			case <-ctx.Done():
				return // Exit early if the context is cancelled
			}

		case <-ctx.Done():
			return // Exit early if the context is cancelled
		}
	}
}

func (p *Persistence) parallelWrite(
	ctx context.Context,
	opt persistencemodel.WriteOptions,
	groups []persistutils.GroupedWriteOperation,
	maxParallelism int,
) []persistencemodel.WriteResult {
	queue := make(chan int, len(groups))
	results := make(chan parallelWriteResult, len(groups))

	var wg sync.WaitGroup

	// Start workers
	for w := 0; w < maxParallelism; w++ {
		wg.Add(1)
		go parallelWriteWorker(ctx, &wg, queue, results, p, opt, groups)
	}

	// Enqueue work
	for i := range groups {
		select {
		case queue <- i:
		case <-ctx.Done():
			close(queue) // Close queue to prevent further enqueuing
			wg.Wait()    // Wait for workers to finish processing
			close(results)
			return nil // Return an empty result set if context is cancelled
		}
	}

	close(queue) // Close the queue after all tasks are enqueued

	// Wait for workers to finish and close results channel
	go func() {
		wg.Wait()
		close(results)
	}()

	// Dynamically collect results
	var res []persistencemodel.WriteResult
	for r := range results {
		res = append(res, r.data...)
	}

	return res
}
