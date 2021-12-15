package worker

import (
	"context"
	"fmt"

	"github.com/spyzhov/ajson"
	"github.com/w-haibara/kuirejo/compiler"
	"github.com/w-haibara/kuirejo/task"
)

func (w Workflow) evalTask(ctx context.Context, state *compiler.TaskState, input *ajson.Node) (*ajson.Node, error) {
	out, err := task.Do(ctx, state.Resouce.Type, state.Resouce.Path, input)
	if err != nil {
		return nil, fmt.Errorf("task.Do() failed: %w", err)
	}

	return out, nil
}
