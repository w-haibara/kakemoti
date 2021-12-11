package worker

import (
	"context"
	"fmt"

	"github.com/spyzhov/ajson"
	"github.com/w-haibara/kuirejo/compiler"
)

func (w Workflow) evalFail(ctx context.Context, state *compiler.FailState, input *ajson.Node) (*ajson.Node, error) {
	return nil, fmt.Errorf("state machine failed: %w", ErrStateMachineTerminated)
}
