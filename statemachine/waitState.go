package statemachine

import (
	"bytes"
	"context"
	"fmt"
	"strings"
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

func (s *WaitState) Dulation(r *bytes.Buffer) (time.Duration, error) {
	if s == nil {
		return time.Duration(0), nil
	}

	if s.Seconds > 0 {
		return time.Duration(s.Seconds) * time.Second, nil
	}

	if s.Timestamp != "" {
		return parseTimestamp(s.Timestamp)
	}

	root, err := ajson.Unmarshal(r.Bytes())
	if err != nil {
		return time.Duration(0), err
	}

	if s.SecondsPath != "" {
		nodes, err := root.JSONPath(s.SecondsPath)
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
		nodes, err := root.JSONPath(s.TimestampPath)
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

func (s *WaitState) Transition(ctx context.Context, r, w *bytes.Buffer) (next string, err error) {
	if s == nil {
		return "", nil
	}

	select {
	case <-ctx.Done():
		return "", ErrStoppedStateMachine
	default:
	}

	d, err := s.Dulation(r)
	if err != nil {
		return "", err
	}
	time.Sleep(d)

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
