package worker

import (
	"context"

	"github.com/w-haibara/kakemoti/compiler"
)

func (w Workflow) evalPass(ctx context.Context, state *compiler.PassState, input interface{}) (interface{}, error) {
	output := state.Result
	if output == nil {
		output = input
	}
	return output, nil
}
