package statemachine

import (
	"bytes"
)

type SucceedState struct {
	CommonState
}

func (s *SucceedState) Transition(r, w *bytes.Buffer) (next string, err error) {
	return "", ErrSucceededStateMachine
}
