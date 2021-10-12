package statemachine

import (
	"bytes"
	"context"

	"github.com/sirupsen/logrus"
)

type State interface {
	SetName(name string)
	SetID(id string)
	StateType() string
	String() string
	Transition(ctx context.Context, r, w *bytes.Buffer) (next string, err error)
	SetLogger(l *logrus.Entry)
	Logger() *logrus.Entry
}
