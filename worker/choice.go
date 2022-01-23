package worker

import (
	"context"

	"github.com/w-haibara/kakemoti/compiler"
)

func (w Workflow) evalChoice(ctx context.Context, coj *compiler.CtxObj, state compiler.ChoiceState, input interface{}) (string, interface{}, statesError) {
	for _, choice := range state.Choices {
		ok, err := choice.Condition.Eval(coj, input)
		if err != nil {
			return "", nil, NewStatesError("", err)
		}
		if ok {
			return choice.Next, input, NewStatesError("", nil)
		}
	}

	if state.Default == "" {
		return "", nil, NewStatesError(StatesErrorBranchFailed, nil)
	}

	return state.Default, input, NewStatesError("", nil)
}
