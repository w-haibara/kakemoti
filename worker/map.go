package worker

import (
	"context"
	"errors"
	"fmt"

	"github.com/w-haibara/kakemoti/compiler"
	"golang.org/x/sync/errgroup"
)

func (w Workflow) evalMap(ctx context.Context, state compiler.MapState, input interface{}) (interface{}, statesError) {
	iter, err := NewWorkflow(&state.Iterator, w.Logger)
	if err != nil {
		return nil, statesError{"", err}
	}

	v, err := compiler.UnjoinByPath(ctx, input, &state.ItemsPath)
	if err != nil {
		return nil, NewStatesError("", err)
	}
	items, ok := v.([]interface{})
	if !ok {
		return nil, NewStatesError("", fmt.Errorf("input for Map must be an array: [%v]", v))
	}

	var outputs parallelOutputs
	outputs.v = make([]interface{}, len(items))

	var eg errgroup.Group
	count := state.MaxConcurrency
	for i, item := range items {
		i := i
		item := item
		fn := func() error {
			o, err := iter.Exec(ctx, item)
			if !errors.Is(err, ErrStateMachineTerminated) && err != nil {
				return err
			}

			outputs.mu.Lock()
			outputs.v[i] = o
			outputs.mu.Unlock()

			return nil
		}

		count--
		if count > 0 {
			eg.Go(fn)
		} else {
			if err := fn(); err != nil {
				return nil, statesError{"", err}
			}
		}
	}

	if err := eg.Wait(); err != nil {
		return nil, NewStatesError("", err)
	}

	return outputs.v, NewStatesError("", nil)
}
