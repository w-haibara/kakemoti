package statemachine

import (
	"bytes"
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

func (s *CommonState) StateType() string {
	return s.Type
}

func (s *CommonState) String() string {
	return pp.Sprintln(s)
}

func (s *CommonState) Transition(r, w *bytes.Buffer) (next string, err error) {
	if s.End {
		return "", ErrEndStateMachine
	}

	if strings.TrimSpace(s.Next) == "" {
		return "", ErrNextStateIsBrank
	}

	return s.Next, nil
}
