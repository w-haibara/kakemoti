package worker

import (
	"context"

	"github.com/spyzhov/ajson"
	"github.com/w-haibara/kuirejo/compiler"
)

func (w Workflow) evalPass(ctx context.Context, state *compiler.PassState, input *ajson.Node) (*ajson.Node, error) {
	return input, nil
}
