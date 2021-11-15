package statemachine

import (
	"context"

	"github.com/spyzhov/ajson"
)

type FailState struct {
	CommonState
	Cause string `json:"Cause"`
	Error string `json:"Error"`
}

func (s *FailState) Transition(ctx context.Context, r *ajson.Node) (next string, w *ajson.Node, err error) {
	return s.CommonState.TransitionWithoutPostCheck(ctx, r, func(ctx context.Context, r *ajson.Node) (string, *ajson.Node, error) {
		return "", nil, ErrFailedStateMachine
	})
}
