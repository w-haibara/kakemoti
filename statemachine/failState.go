package statemachine

import (
	"bytes"
)

type FailState struct {
	CommonState
	Cause string `json:"Cause"`
	Error string `json:"Error"`
}

func (s *FailState) Transition(r, w *bytes.Buffer) (next string, err error) {
	if s == nil {
		return "", nil
	}

	return "", ErrFailedStateMachine
}
