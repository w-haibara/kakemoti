package worker

import (
	"context"
	"fmt"

	"github.com/spyzhov/ajson"
	"github.com/w-haibara/kuirejo/compiler"
)

func (w Workflow) evalSucceed(ctx context.Context, state *compiler.SucceedState, input *ajson.Node) (*ajson.Node, error) {
	return nil, fmt.Errorf("state machine failed: %w", ErrStateMachineTerminated)
}
