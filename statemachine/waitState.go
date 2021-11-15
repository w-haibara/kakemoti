package statemachine

import (
	"context"
	"fmt"
	"time"

	"github.com/spyzhov/ajson"
)

type WaitState struct {
	CommonState
	Seconds       int64  `json:"Seconds"`
	Timestamp     string `json:"Timestamp"`
	SecondsPath   string `json:"SecondsPath"`
	TimestampPath string `json:"TimestampPath"`
}

func parseTimestamp(timestamp string) (time.Duration, error) {
	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return time.Duration(0), err
	}

	return time.Until(t), nil
}

func (s *WaitState) dulation(r *ajson.Node) (time.Duration, error) {
	if s == nil {
		return time.Duration(0), nil
	}

	if s.Seconds > 0 {
		return time.Duration(s.Seconds) * time.Second, nil
	}

	if s.Timestamp != "" {
		return parseTimestamp(s.Timestamp)
	}

	if s.SecondsPath != "" {
		nodes, err := r.JSONPath(s.SecondsPath)
		if err != nil {
			return time.Duration(0), err
		}

		if len(nodes) < 1 {
			return time.Duration(0), fmt.Errorf("JSONPath result is empty")
		}

		v, err := nodes[0].GetNumeric()
		if err != nil {
			return time.Duration(0), err
		}

		return time.Duration(v) * time.Second, nil
	}

	if s.TimestampPath != "" {
		nodes, err := r.JSONPath(s.TimestampPath)
		if err != nil {
			return time.Duration(0), err
		}

		if len(nodes) < 1 {
			return time.Duration(0), fmt.Errorf("JSONPath result is empty")
		}

		v, err := nodes[0].GetString()
		if err != nil {
			return time.Duration(0), err
		}

		return parseTimestamp(v)

	}

	return time.Duration(0), fmt.Errorf("wait dulation is not set")
}

func (s *WaitState) Transition(ctx context.Context, r *ajson.Node) (next string, w *ajson.Node, err error) {
	return s.CommonState.TransitionWithEndNext(ctx, r, func(ctx context.Context, r *ajson.Node) (string, *ajson.Node, error) {
		d, err := s.dulation(r)
		if err != nil {
			return "", nil, err
		}

		time.Sleep(d)

		return "", nil, nil
	})
}
