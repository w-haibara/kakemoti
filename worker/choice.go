package worker

import (
	"context"

	"github.com/spyzhov/ajson"
	"github.com/w-haibara/kuirejo/compiler"
)

func (w Workflow) evalChoice(ctx context.Context, state *compiler.ChoiceState, input *ajson.Node) (string, *ajson.Node, error) {
	next := "Yes"
	return next, input, nil
}
