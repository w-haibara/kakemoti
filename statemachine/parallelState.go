package statemachine

import (
	"context"
	"sync"

	"github.com/spyzhov/ajson"
	"golang.org/x/sync/errgroup"
)

type ParallelState struct {
	CommonState
	Branches       []StateMachine `json:"Branches"`
	ResultPath     string         `json:"ResultPath"`
	ResultSelector string         `json:"ResultSelector"`
	Retry          string         `json:"Retry"`
	Catch          string         `json:"Catch"`
}

type outputs struct {
	mu sync.Mutex
	v  []*ajson.Node
}

func (s *ParallelState) Transition(ctx context.Context, r *ajson.Node) (next string, w *ajson.Node, err error) {
	return s.CommonState.TransitionWithResultselectorRetry(ctx, r,
		nil, s.ResultPath,
		nil, s.Retry, s.Catch,
		func(ctx context.Context, r *ajson.Node) (string, *ajson.Node, error) {
			var eg errgroup.Group
			var outputs outputs
			outputs.v = make([]*ajson.Node, len(s.Branches))

			for i, sm := range s.Branches {
				i, sm := i, sm
				sm.Logger = s.logger

				eg.Go(func() error {
					w, err := sm.start(ctx, r)
					if err != nil {
						return err
					}

					outputs.mu.Lock()
					outputs.v[i] = w.Clone()
					outputs.mu.Unlock()

					return nil
				})
			}

			if err := eg.Wait(); err != nil {
				return "", nil, err
			}

			return "", ajson.ArrayNode(s.StateMachineID, outputs.v), nil
		})
}
