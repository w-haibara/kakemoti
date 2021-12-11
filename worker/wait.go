package worker

import (
	"context"

	"github.com/spyzhov/ajson"
	"github.com/w-haibara/kuirejo/compiler"
)

func (w Workflow) evalWait(ctx context.Context, state *compiler.WaitState, input *ajson.Node) (*ajson.Node, error) {
	return input, nil
}
