package statemachine

import (
	"context"
	"fmt"

	"github.com/spyzhov/ajson"
)

type FailState struct {
	CommonState
	Cause string `json:"Cause"`
	Error string `json:"Error"`
}

func (s *FailState) Transition(ctx context.Context, r *ajson.Node) (next string, w *ajson.Node, err error) {
	return s.CommonState.Transition(ctx, r, func(ctx context.Context, r *ajson.Node) (string, *ajson.Node, error) {
		return "", nil, fmt.Errorf("state machine failed: %w", ErrStateMachineTerminated)
	})
}
