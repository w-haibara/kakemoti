package worker

import (
	"context"
	"fmt"

	"github.com/w-haibara/kakemoti/compiler"
)

func (w Workflow) evalSucceed(ctx context.Context, state compiler.SucceedState, input interface{}) (interface{}, statesError) {
	return input, NewStatesError("", fmt.Errorf("Succeed: %w", ErrStateMachineTerminated))
}
