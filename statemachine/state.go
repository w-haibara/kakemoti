package statemachine

import (
	"bytes"
	"context"
	"karage/log"
)

type State interface {
	SetName(name string)
	SetID(id string)
	StateType() string
	String() string
	Transition(ctx context.Context, r, w *bytes.Buffer) (next string, err error)
	SetLogger(l *log.Logger)
	GetLogger() *log.Logger
	Log(v ...interface{})
	StateStartLog()
	StateEndLog()
}
