package worker

import (
	"context"

	"github.com/w-haibara/kuirejo/compiler"
)

func (w Workflow) evalMap(ctx context.Context, state *compiler.MapState, input interface{}) (interface{}, error) {
	return input, nil
}
