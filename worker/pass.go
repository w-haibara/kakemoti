package worker

import (
	"context"

	"github.com/w-haibara/kakemoti/controller/compiler"
)

func (w Workflow) evalPass(ctx context.Context, state compiler.PassState, input interface{}) (interface{}, statesError) {
	output := state.Result
	if output == nil {
		output = input
	}
	return output, NewStatesError("", nil)
}
