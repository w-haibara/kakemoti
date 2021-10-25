package statemachine

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/spyzhov/ajson"
)

type State interface {
	SetName(name string)
	SetID(id string)
	StateType() string
	String() string
	FilterInput(ctx context.Context, input *ajson.Node) (*ajson.Node, error)
	FilterOutput(ctx context.Context, output *ajson.Node) (*ajson.Node, error)
	Transition(ctx context.Context, r *ajson.Node) (next string, w *ajson.Node, err error)
	SetLogger(v *logrus.Entry)
	Logger(v logrus.Fields) *logrus.Entry
}
