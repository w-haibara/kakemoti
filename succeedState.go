package main

import (
	"bytes"

	"github.com/k0kubun/pp"
)

type SucceedState struct {
}

func (s SucceedState) Transition(r, w *bytes.Buffer) (next string, err error) {
	return "", EndStateMachine
}

func (s SucceedState) Print() {
	pp.Print(s)
}
