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
	if s == nil {
		return "", nil, nil
	}

	select {
	case <-ctx.Done():
		return "", nil, ErrStoppedStateMachine
	default:
	}

	return "", r, ErrFailedStateMachine
}
