package main

import ()

type SucceedState struct {
}

func (s SucceedState) Transition() (next string, err error) {
	return "", EndStateMachine
}
