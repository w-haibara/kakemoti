package statemachine

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"karage/log"

	"github.com/k0kubun/pp"
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
	Logger         *log.Logger `json:"-"`
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

func (s *CommonState) SetLogger(l *log.Logger) {
	s.Logger = l
}

func (s *CommonState) GetLogger() *log.Logger {
	return s.Logger
}

func (s *CommonState) Log(v ...interface{}) {
	s.Logger.Println(s.StateMachineID, s.Name, s.Type, fmt.Sprint(v...))
}

func (s *CommonState) StateStartLog() {
	s.Log("START")
}

func (s *CommonState) StateEndLog() {
	s.Log("END")
}
