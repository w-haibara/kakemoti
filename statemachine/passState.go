package statemachine

import (
	"bytes"
	"context"
	"strings"
)

type PassState struct {
	CommonState
	Result     string `json:"Result"`
	ResultPath string `json:"ResultPath"`
	Parameters string `json:"Parameters"`
}

func (s *PassState) Transition(ctx context.Context, r, w *bytes.Buffer) (next string, err error) {
	if s == nil {
		return "", nil
	}

	select {
	case <-ctx.Done():
		return "", ErrStoppedStateMachine
	default:
	}

	if _, err := r.WriteTo(w); err != nil {
		return "", err
	}

	if s.End {
		return "", ErrEndStateMachine
	}

	if strings.TrimSpace(s.Next) == "" {
		return "", ErrNextStateIsBrank
	}

	return s.Next, nil
}
