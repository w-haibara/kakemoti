package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/spyzhov/ajson"
	"github.com/w-haibara/kuirejo/compiler"
)

var timeformat = "2006-01-02T15:04:05Z"

func (w Workflow) evalWait(ctx context.Context, state *compiler.WaitState, input *ajson.Node) (*ajson.Node, error) {
	d, err := getDulation(state, input)
	if err != nil {
		return nil, err
	}

	w.loggerWithInfo().Printf("Wait %s from %s", d, time.Now())
	time.Sleep(d)

	return input, nil
}

func getDulation(state *compiler.WaitState, input *ajson.Node) (time.Duration, error) {
	seconds := float64(state.Seconds)
	if state.SecondsPath != "" {
		nodes, err := input.JSONPath(state.SecondsPath)
		if err != nil {
			return 0, fmt.Errorf("input.JSONPath(path) failed: %w", err)
		}

		if len(nodes) != 1 {
			return 0, fmt.Errorf("invalid length of input.JSONPath(path) result")
		}

		v, err := nodes[0].GetNumeric()
		if err != nil {
			return 0, fmt.Errorf("invalid type of input.JSONPath(path) result")
		}

		seconds = v
	}

	if seconds != 0 {
		return time.Duration(seconds) * time.Second, nil
	}

	timestamp := state.Timestamp
	if state.TimestampPath != "" {
		nodes, err := input.JSONPath(state.TimestampPath)
		if err != nil {
			return 0, fmt.Errorf("input.JSONPath(path) failed: %w", err)
		}

		if len(nodes) != 1 {
			return 0, fmt.Errorf("invalid length of input.JSONPath(path) result")
		}

		v, err := nodes[0].GetString()
		if err != nil {
			return 0, fmt.Errorf("invalid type of input.JSONPath(path) result")
		}

		timestamp = v
	}

	t, err := time.ParseInLocation(timeformat, timestamp, time.Now().Location())
	if err != nil {
		return 0, err
	}

	return time.Until(t), nil
}
