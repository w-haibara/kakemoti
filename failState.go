package main

import (
	"bytes"
)

type FailState struct {
	CommonState
	Cause string `json:"Cause"`
	Error string `json:"Error"`
}

func (s FailState) Transition(r, w *bytes.Buffer) (next string, err error) {
	return "", ErrFailedStateMachine
}
