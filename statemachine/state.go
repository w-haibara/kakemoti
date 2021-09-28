package statemachine

import (
	"bytes"
)

type State interface {
	StateType() string
	String() string
	Transition(r, w *bytes.Buffer) (next string, err error)
}
