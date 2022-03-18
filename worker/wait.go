package worker

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/w-haibara/kakemoti/compiler"
)

var timeformat = "2006-01-02T15:04:05Z"

func (w Workflow) evalWait(ctx context.Context, coj *compiler.CtxObj, state compiler.WaitState, input interface{}) (interface{}, statesError) {
	d, err := getDulation(ctx, coj, state, input)
	if err != nil {
		return nil, NewStatesError("", err)
	}

	log.WithFields(w.infoFields()).Printf("Wait %s from %s", d, time.Now())
	time.Sleep(d)

	return input, NewStatesError("", nil)
}

func getDulation(ctx context.Context, coj *compiler.CtxObj, state compiler.WaitState, input interface{}) (time.Duration, error) {
	switch {
	case state.Seconds != nil:
		if *state.Seconds == 0 {
			return 0, nil
		}

		return time.Duration(*state.Seconds) * time.Second, nil
	case state.SecondsPath != nil:
		v, err := compiler.UnjoinByPath(coj, input, state.SecondsPath)
		if err != nil {
			return 0, err
		}

		seconds, ok := v.(int)
		if !ok {
			return 0, fmt.Errorf("invalid type of input.Path(path) result")
		}

		if seconds == 0 {
			return 0, nil
		}

		return time.Duration(seconds) * time.Second, nil
	case state.Timestamp != nil:
		return time.Until(state.Timestamp.Time), nil
	case state.TimestampPath != nil:
		v, err := compiler.UnjoinByPath(coj, input, state.TimestampPath)
		if err != nil {
			return 0, err
		}

		timestamp, ok := v.(string)
		if !ok {
			return 0, fmt.Errorf("invalid type of input.Path(path) result")
		}

		t, err := time.ParseInLocation(timeformat, timestamp, time.Now().Location())
		if err != nil {
			return 0, err
		}

		return time.Until(t), nil
	}
	return 0, nil
}
