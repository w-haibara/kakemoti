package statemachine

import (
	"bytes"
	"context"
)

type State interface {
	Log(v interface{})
	StateStartLog(name string)
	StateEndLog(name string)
	StateType() string
	String() string
	Transition(ctx context.Context, r, w *bytes.Buffer) (next string, err error)
}
