package worker

import (
	"context"

	"github.com/spyzhov/ajson"
	"github.com/w-haibara/kuirejo/compiler"
)

func (w Workflow) evalMap(ctx context.Context, state *compiler.MapState, input *ajson.Node) (*ajson.Node, error) {
	return input, nil
}
