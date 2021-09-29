package statemachine

import (
	"bytes"
	"context"
)

type State interface {
	StateType() string
	String() string
	Transition(ctx context.Context, r, w *bytes.Buffer) (next string, err error)
}
