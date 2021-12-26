package worker

import (
	"context"
	"fmt"

	"github.com/w-haibara/kuirejo/compiler"
	"github.com/w-haibara/kuirejo/task"
)

func (w Workflow) evalTask(ctx context.Context, state *compiler.TaskState, input interface{}) (interface{}, error) {
	out, err := task.Do(ctx, state.Resouce.Type, state.Resouce.Path, input)
	if err != nil {
		return nil, fmt.Errorf("task.Do() failed: %v", err)
	}

	return out, nil
}
