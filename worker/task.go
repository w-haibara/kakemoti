package worker

import (
	"context"

	"github.com/spyzhov/ajson"
	"github.com/w-haibara/kuirejo/compiler"
)

func (w Workflow) evalTask(ctx context.Context, state *compiler.TaskState, input *ajson.Node) (*ajson.Node, error) {
	return input, nil
}
