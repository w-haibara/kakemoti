package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/w-haibara/kakemoti/compiler"
)

var timeformat = "2006-01-02T15:04:05Z"

func (w Workflow) evalWait(ctx context.Context, state *compiler.WaitState, input interface{}) (interface{}, statesError) {
	d, err := getDulation(ctx, state, input)
	if err != nil {
		return nil, NewStatesError("", err)
	}

	w.loggerWithInfo().Printf("Wait %s from %s", d, time.Now())
	time.Sleep(d)

	return input, NewStatesError("", nil)
}

func getDulation(ctx context.Context, state *compiler.WaitState, input interface{}) (time.Duration, error) {
	switch {
	case state.Seconds != nil || state.SecondsPath != nil:
		var seconds int64
		if state.Seconds != nil {
			seconds = *state.Seconds
		}
		if state.SecondsPath != nil {
			v, err := UnjoinByJsonPath(ctx, input, *state.SecondsPath)
			if err != nil {
				return 0, err
			}

			if v, ok := v.(int64); !ok {
				return 0, fmt.Errorf("invalid type of input.JSONPath(path) result")
			} else {
				seconds = v
			}
		}
		if seconds == 0 {
			return 0, nil
		}
		return time.Duration(seconds) * time.Second, nil
	case state.Timestamp != nil || state.TimestampPath != nil:
		timestamp := ""
		if state.Timestamp != nil {
			timestamp = *state.Timestamp
		}
		if state.TimestampPath != nil {
			v, err := UnjoinByJsonPath(ctx, input, *state.TimestampPath)
			if err != nil {
				return 0, err
			}

			if v, ok := v.(string); !ok {
				return 0, fmt.Errorf("invalid type of input.JSONPath(path) result")
			} else {
				timestamp = v
			}
		}
		t, err := time.ParseInLocation(timeformat, timestamp, time.Now().Location())
		if err != nil {
			return 0, err
		}
		return time.Until(t), nil
	}
	return 0, nil
}
