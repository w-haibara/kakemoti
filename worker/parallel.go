package worker

import (
	"context"
	"errors"
	"sync"

	"github.com/w-haibara/kakemoti/controller/compiler"
	"golang.org/x/sync/errgroup"
)

type parallelOutputs struct {
	mu sync.Mutex
	v  []interface{}
}

func (w Workflow) evalParallel(ctx context.Context, coj *compiler.CtxObj, state compiler.ParallelState, input interface{}) (interface{}, statesError) {
	var eg errgroup.Group
	var outputs parallelOutputs
	outputs.v = make([]interface{}, len(state.Branches))
	for i := range state.Branches {
		i := i
		eg.Go(func() error {
			w, err := NewWorkflow(&state.Branches[i])
			if err != nil {
				return err
			}

			o, err := w.Exec(ctx, coj, input)
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
		return nil, NewStatesError(StatesErrorBranchFailed, err)
	}

	return outputs.v, NewStatesError("", nil)
}
