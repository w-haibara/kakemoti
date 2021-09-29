package statemachine

import (
	"bytes"
	"context"
)

type FailState struct {
	CommonState
	Cause string `json:"Cause"`
	Error string `json:"Error"`
}

func (s *FailState) Transition(ctx context.Context, r, w *bytes.Buffer) (next string, err error) {
	if s == nil {
		return "", nil
	}

	select {
	case <-ctx.Done():
		return "", ErrStoppedStateMachine
	default:
	}

	return "", ErrFailedStateMachine
}
