package main

import (
	"github.com/k0kubun/pp"
)

type SucceedState struct {
}

func (s SucceedState) Transition() (next string, err error) {
	return "", EndStateMachine
}

func (s SucceedState) Print() {
	pp.Print(s)
}
