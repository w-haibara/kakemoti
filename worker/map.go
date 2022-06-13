package worker

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/w-haibara/kakemoti/controller/compiler"
	"golang.org/x/sync/errgroup"
)

type mapOutputs struct {
	mu sync.Mutex
	v  []interface{}
}

func (w Workflow) evalMap(ctx context.Context, coj *compiler.CtxObj, state compiler.MapState, input interface{}) (interface{}, statesError) {
	iter, err := NewWorkflow(&state.Iterator)
	if err != nil {
		return nil, statesError{"", err}
	}

	v, err := compiler.UnjoinByPath(coj, input, &state.ItemsPath.Path)
	if err != nil {
		return nil, NewStatesError("", err)
	}
	items, ok := v.([]interface{})
	if !ok {
		return nil, NewStatesError("", fmt.Errorf("input for Map must be an array: [%v]", v))
	}

	var outputs mapOutputs
	outputs.v = make([]interface{}, len(items))
	eg := new(errgroup.Group)
	count := 0
	for i := range items {
		i := i
		eg.Go(func() error {
			c := new(compiler.CtxObj)
			c1, err := c.SetByString("$.Map.Item.Index", i)
			if err != nil {
				return err
			}
			c2, err := c1.SetByString("$.Map.Item.Value", items[i])
			if err != nil {
				return err
			}
			c3, err := c2.SetAll(coj.GetAll())
			if err != nil {
				return err
			}

			o, err := iter.Exec(ctx, c3, items[i])
			if !errors.Is(err, ErrStateMachineTerminated) && err != nil {
				return err
			}

			outputs.mu.Lock()
			outputs.v[i] = o
			outputs.mu.Unlock()

			return nil
		})

		count++
		if count > state.MaxConcurrency {
			count = 0
			if err := eg.Wait(); err != nil {
				return nil, NewStatesError("", err)
			}
		}
	}

	if err := eg.Wait(); err != nil {
		return nil, NewStatesError("", err)
	}

	return outputs.v, NewStatesError("", nil)
}
