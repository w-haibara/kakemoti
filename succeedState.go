package main

import (
	"io"

	"github.com/k0kubun/pp"
)

type SucceedState struct {
}

func (s SucceedState) Transition(r io.Reader, w io.Writer) (next string, err error) {
	return "", EndStateMachine
}

func (s SucceedState) Print() {
	pp.Print(s)
}
