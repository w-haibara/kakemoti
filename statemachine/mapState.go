package statemachine

import (
	"context"

	"github.com/spyzhov/ajson"
)

type MapState struct {
	CommonState
	Iterator       StateMachine `json:"Iterator"`
	ItemsPath      string       `json:"ItemsPath"`
	MaxConcurrency int64        `json:"MaxConcurrency"`
	ResultPath     string       `json:"ResultPath"`
	ResultSelector string       `json:"ResultSelector"`
	Retry          string       `json:"Retry"`
	Catch          string       `json:"Catch"`
}

func (s *MapState) Transition(ctx context.Context, r *ajson.Node) (next string, w *ajson.Node, err error) {
	return s.CommonState.TransitionWithEndNext(ctx, r, nil)
}
