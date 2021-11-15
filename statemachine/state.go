package statemachine

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/spyzhov/ajson"
	"github.com/w-haibara/kuirejo/log"
)

type State interface {
	SetName(name string)
	SetID(id string)
	StateType() string
	String() string
	FilterInput(ctx context.Context, input *ajson.Node) (*ajson.Node, error)
	FilterOutput(ctx context.Context, output *ajson.Node) (*ajson.Node, error)
	Transition(ctx context.Context, r *ajson.Node) (next string, w *ajson.Node, err error)
	SetLogger(v *log.Logger)
	Logger(v logrus.Fields) *log.Logger
}
