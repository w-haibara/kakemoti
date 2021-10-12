package statemachine

import (
	"context"
	"strings"

	"github.com/spyzhov/ajson"
)

type PassState struct {
	CommonState
	Result     string `json:"Result"`
	ResultPath string `json:"ResultPath"`
	Parameters string `json:"Parameters"`
}

func (s *PassState) Transition(ctx context.Context, r *ajson.Node) (next string, w *ajson.Node, err error) {
	if s == nil {
		return "", nil, nil
	}

	select {
	case <-ctx.Done():
		return "", nil, ErrStoppedStateMachine
	default:
	}

	if s.End {
		return "", r, ErrEndStateMachine
	}

	if strings.TrimSpace(s.Next) == "" {
		return "", nil, ErrNextStateIsBrank
	}

	return s.Next, r, nil
}
