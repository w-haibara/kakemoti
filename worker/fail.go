package worker

import (
	"context"
	"fmt"

	"github.com/w-haibara/kakemoti/compiler"
)

func (w Workflow) evalFail(ctx context.Context, state *compiler.FailState, input interface{}) (interface{}, error) {
	return input, fmt.Errorf("Fail: %w", ErrStateMachineTerminated)
}
