package worker

import (
	"context"
	"errors"
	"sync"

	"github.com/w-haibara/kakemoti/compiler"
	"golang.org/x/sync/errgroup"
)

type outputs struct {
	mu sync.Mutex
	v  []interface{}
}

func (w Workflow) evalParallel(ctx context.Context, state *compiler.ParallelState, input interface{}) (interface{}, statesError) {
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
			if !errors.Is(err, ErrStateMachineTerminated) && err != nil {
				return err
			}

			outputs.mu.Lock()
			outputs.v[i] = o
			outputs.mu.Unlock()

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, NewStatesError("", err)
	}

	return outputs.v, NewStatesError("", nil)
}
