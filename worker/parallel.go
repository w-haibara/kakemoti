package worker

import (
	"context"
	"sync"

	"github.com/w-haibara/kuirejo/compiler"
	"golang.org/x/sync/errgroup"
)

type outputs struct {
	mu sync.Mutex
	v  []interface{}
}

func (w Workflow) evalParallel(ctx context.Context, state *compiler.ParallelState, input interface{}) (interface{}, error) {
	var eg errgroup.Group
	var outputs outputs
	outputs.v = make([]interface{}, len(state.Branches))
	for i := range state.Branches {
		i := i
		eg.Go(func() error {
			w, err := NewWorkflow(&state.Branches[i], w.Logger)
			if err != nil {
				return err
			}

			o, err := w.Exec(ctx, input)
			if err != nil {
				return err
			}

			outputs.mu.Lock()
			outputs.v[i] = o
			outputs.mu.Unlock()

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return outputs.v, nil
}
