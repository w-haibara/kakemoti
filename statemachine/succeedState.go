package statemachine

import (
	"context"

	"github.com/spyzhov/ajson"
)

type SucceedState struct {
	CommonState
}

func (s *SucceedState) Transition(ctx context.Context, r *ajson.Node) (next string, w *ajson.Node, err error) {
	if s == nil {
		return "", nil, nil
	}

	select {
	case <-ctx.Done():
		return "", nil, ErrStoppedStateMachine
	default:
	}

	return "", r, ErrSucceededStateMachine
}
