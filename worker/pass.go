package worker

import (
	"context"

	"github.com/w-haibara/kuirejo/compiler"
)

func (w Workflow) evalPass(ctx context.Context, state *compiler.PassState, input interface{}) (interface{}, error) {
	return input, nil
}
