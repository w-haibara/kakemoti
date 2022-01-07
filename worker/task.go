package worker

import (
	"context"

	"github.com/w-haibara/kakemoti/compiler"
	"github.com/w-haibara/kakemoti/task"
)

func (w Workflow) evalTask(ctx context.Context, state *compiler.TaskState, input interface{}) (interface{}, statesError) {
	out, err := task.Do(ctx, state.Resouce.Type, state.Resouce.Path, input)
	if err != nil {
		return nil, NewStatesError(StatesErrorTaskFailed, err)
	}

	return out, NewStatesError("", nil)
}
