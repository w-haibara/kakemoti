package worker

import (
	"context"
	"errors"
	"os"

	"github.com/w-haibara/kakemoti/controller/compiler"
	"github.com/w-haibara/kakemoti/task"
)

func (w Workflow) evalTask(ctx context.Context, state compiler.TaskState, input interface{}) (interface{}, statesError) {
	out, stateserr, err := task.Do(ctx, state.Resouce.Type, state.Resouce.Path, input)
	if stateserr != "" {
		return nil, NewStatesError(stateserr, err)
	}
	if err != nil {
		if errors.Is(err, os.ErrPermission) {
			return nil, NewStatesError(StatesErrorPermissions, err)
		}
		return nil, NewStatesError(StatesErrorTaskFailed, err)
	}

	return out, NewStatesError(stateserr, nil)
}
