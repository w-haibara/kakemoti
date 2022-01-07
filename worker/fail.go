package worker

import (
	"context"
	"fmt"

	"github.com/w-haibara/kakemoti/compiler"
)

func (w Workflow) evalFail(ctx context.Context, state *compiler.FailState, input interface{}) (interface{}, statesError) {
	return input, NewStatesError("", fmt.Errorf("Fail: %w", ErrStateMachineTerminated))
}
