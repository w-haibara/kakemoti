package worker

import (
	"context"
	"errors"
	"fmt"

	"github.com/w-haibara/kakemoti/compiler"
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

	for i, item := range items {
		o, err := iter.Exec(ctx, item)
		if !errors.Is(err, ErrStateMachineTerminated) && err != nil {
			return nil, statesError{"", err}
		}

		outputs.mu.Lock()
		outputs.v[i] = o
		outputs.mu.Unlock()
	}

	return outputs.v, NewStatesError("", nil)
}
