package worker

import (
	"context"
	"fmt"

	"github.com/w-haibara/kuirejo/compiler"
)

func (w Workflow) evalFail(ctx context.Context, state *compiler.FailState, input interface{}) (interface{}, error) {
	return nil, fmt.Errorf("state machine failed: %w", ErrStateMachineTerminated)
}
