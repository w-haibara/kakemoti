package worker

import (
	"context"

	"github.com/w-haibara/kakemoti/compiler"
)

func (w Workflow) evalMap(ctx context.Context, state *compiler.MapState, input interface{}) (interface{}, statesError) {
	return input, NewStatesError("", nil)
}
