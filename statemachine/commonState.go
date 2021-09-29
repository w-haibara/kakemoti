package statemachine

import (
	"bytes"
	"context"
	"strings"

	"github.com/k0kubun/pp"
)

type CommonState struct {
	Type       string `json:"Type"`
	Next       string `json:"Next"`
	End        bool   `json:"End"`
	Comment    string `json:"Comment"`
	InputPath  string `json:"InputPath"`
	OutputPath string `json:"OutputPath"`
}

func (s CommonState) Log(v interface{}) {
	_, _ = pp.Println(v)
}

func (s CommonState) StateStartLog(name string) {
	s.Log("State Start: " + s.Type + ", " + name)
}

func (s CommonState) StateEndLog(name string) {
	s.Log("State End  : " + s.Type + ", " + name)
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

func (s *CommonState) Transition(ctx context.Context, r, w *bytes.Buffer) (next string, err error) {
	if s == nil {
		return "", nil
	}

	select {
	case <-ctx.Done():
		return "", ErrStoppedStateMachine
	default:
	}

	if s.End {
		return "", ErrEndStateMachine
	}

	if strings.TrimSpace(s.Next) == "" {
		return "", ErrNextStateIsBrank
	}

	return s.Next, nil
}
