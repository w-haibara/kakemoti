package worker

import (
	"context"

	"github.com/w-haibara/kakemoti/compiler"
)

func (w Workflow) evalChoice(ctx context.Context, state *compiler.ChoiceState, input interface{}) (string, interface{}, statesError) {
	for _, choice := range state.Choices {
		ok, err := choice.Condition.Eval(ctx, input)
		if err != nil {
			return "", nil, NewStatesError("", err)
		}
		if ok {
			return choice.Next, input, NewStatesError("", nil)
		}
	}

	return state.Default, input, NewStatesError("", nil)
}
