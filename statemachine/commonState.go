package statemachine

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/k0kubun/pp"
	"github.com/sirupsen/logrus"
	"github.com/spyzhov/ajson"
	"github.com/w-haibara/kuirejo/log"
)

type CommonState struct {
	Name           string      `json:"-"`
	StateMachineID string      `json:"-"`
	Type           string      `json:"Type"`
	Next           string      `json:"Next"`
	End            bool        `json:"End"`
	Comment        string      `json:"Comment"`
	InputPath      string      `json:"InputPath"`
	OutputPath     string      `json:"OutputPath"`
	logger         *log.Logger `json:"-"`
}

func (s *CommonState) SetName(name string) {
	s.Name = name
}

func (s *CommonState) SetID(id string) {
	s.StateMachineID = id
}

func (s *CommonState) StateType() string {
	if s == nil {
		return ""
	}

	return s.Type
}

func (s *CommonState) String() string {
	if s == nil {
		return ""
	}

	return pp.Sprintln(s)
}

func (s *CommonState) FilterInput(ctx context.Context, input *ajson.Node) (*ajson.Node, error) {
	node, err := filterByInputPath(input, s.InputPath)
	if err != nil {
		return nil, err
	}

	return node, nil
}

func (s *CommonState) FilterOutput(ctx context.Context, output *ajson.Node) (*ajson.Node, error) {
	node, err := filterByOutputPath(output, s.OutputPath)
	if err != nil {
		return nil, err
	}

	return node, nil
}

type TransitionFunc func(ctx context.Context, r *ajson.Node) (string, *ajson.Node, error)

func (fn TransitionFunc) do(ctx context.Context, r *ajson.Node) (string, *ajson.Node, error) {
	if fn == nil {
		return "", nil, nil
	}

	next, w, err := fn(ctx, r)
	if err != nil {
		return "", r, err
	}

	if w != nil {
		r = w
	}

	return next, r, nil
}

func (s *CommonState) Transition(ctx context.Context, r *ajson.Node, fn TransitionFunc) (string, *ajson.Node, error) {
	select {
	case <-ctx.Done():
		return "", nil, ErrStateMachineTerminated
	default:
	}

	if fn != nil {
		next, w, err := fn.do(ctx, r)
		if err != nil {
			return next, r, err
		}

		if w != nil {
			r = w
		}

		if strings.TrimSpace(next) != "" {
			return next, r, nil
		}
	}

	return s.Next, r, nil
}

func (s *CommonState) TransitionWithIO(ctx context.Context, r *ajson.Node, fn TransitionFunc) (string, *ajson.Node, error) {
	if node, err := s.FilterInput(ctx, r); err != nil {
		return "", nil, fmt.Errorf("failed to FilterInput(): %v", err)
	} else {
		r = node
	}

	next, w, err := s.Transition(ctx, r, fn)

	if node, err := s.FilterOutput(ctx, w); err != nil {
		return "", nil, fmt.Errorf("failed to FilterOutput(): %v", err)
	} else {
		w = node
	}

	return next, w, err
}

func (s *CommonState) TransitionWithEndNext(ctx context.Context, r *ajson.Node, fn TransitionFunc) (string, *ajson.Node, error) {
	next, w, err := s.TransitionWithIO(ctx, r, fn)
	if err != nil {
		return next, w, err
	}

	if s.End {
		return "", r, ErrStateMachineTerminated
	}

	return next, w, nil
}

func (s *CommonState) TransitionWithResultpathParameters(ctx context.Context, r *ajson.Node, parameters *json.RawMessage, resultPath string, fn TransitionFunc) (string, *ajson.Node, error) {
	return s.TransitionWithEndNext(ctx, r,
		func(ctx context.Context, r *ajson.Node) (string, *ajson.Node, error) {
			r, err := replaceByParameters(r, parameters)
			if err != nil {
				return "", nil, err
			}

			if fn == nil {
				return s.Next, r, nil
			}

			next, w, err := fn(ctx, r)
			if next != "" {
				s.Next = next
			}

			if w != nil {
				node, err := filterByResultPath(r, w, resultPath)
				if err != nil {
					return "", nil, err
				}
				if node != nil {
					w = node
				}
			}

			return s.Next, w, err
		})
}

func (s *CommonState) TransitionWithResultselectorRetry(ctx context.Context, r *ajson.Node, parameters *json.RawMessage, resultPath string, resultSelector *json.RawMessage, retry, catch string, fn TransitionFunc) (string, *ajson.Node, error) {
	// TODO: Implement Retry & Catch
	return s.TransitionWithResultpathParameters(ctx, r,
		parameters, resultPath,
		func(ctx context.Context, r *ajson.Node) (string, *ajson.Node, error) {
			if fn == nil {
				return s.Next, r, nil
			}

			next, w, err := fn.do(ctx, r)
			if next != "" {
				s.Next = next
			}

			if w != nil {
				node, err := replaceByResultSelector(w, resultSelector)
				if err != nil {
					return "", nil, err
				}
				if node != nil {
					w = node
				}
			}

			return s.Next, w, err
		})
}

func (s *CommonState) SetLogger(v *log.Logger) {
	s.logger = v
}

func (s *CommonState) Logger(v logrus.Fields) *log.Logger {
	l := log.Logger{
		Entry: s.logger.WithFields(logrus.Fields{
			"name": s.Name,
			"type": s.Type,
			"next": s.Next,
			"end":  s.End,
			"line": log.Line(),
		}),
	}

	if v == nil {
		return &l
	}

	return &log.Logger{
		Entry: l.WithFields(v),
	}
}
