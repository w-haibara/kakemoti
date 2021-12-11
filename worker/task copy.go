package worker

import (
	"context"

	"github.com/spyzhov/ajson"
	"github.com/w-haibara/kuirejo/compiler"
)

func (w Workflow) evalParallel(ctx context.Context, state *compiler.ParallelState, input *ajson.Node) (*ajson.Node, error) {
	return input, nil
}
