package statemachine

import (
	"bytes"
	"context"
)

type SucceedState struct {
	CommonState
}

func (s *SucceedState) Transition(ctx context.Context, r, w *bytes.Buffer) (next string, err error) {
	if s == nil {
		return "", nil
	}

	select {
	case <-ctx.Done():
		return "", ErrStoppedStateMachine
	default:
	}

	return "", ErrSucceededStateMachine
}
