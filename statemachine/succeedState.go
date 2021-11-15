package statemachine

import (
	"context"

	"github.com/spyzhov/ajson"
)

type SucceedState struct {
	CommonState
}

func (s *SucceedState) Transition(ctx context.Context, r *ajson.Node) (next string, w *ajson.Node, err error) {
	return s.CommonState.Transition(ctx, r, func(ctx context.Context, r *ajson.Node) (string, *ajson.Node, error) {
		return "", r, ErrSucceededStateMachine
	})
}
