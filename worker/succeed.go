package worker

import (
	"context"
	"fmt"

	"github.com/w-haibara/kuirejo/compiler"
)

func (w Workflow) evalSucceed(ctx context.Context, state *compiler.SucceedState, input interface{}) (interface{}, error) {
	return nil, fmt.Errorf("state machine failed: %w", ErrStateMachineTerminated)
}
