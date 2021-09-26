package statemachine

import (
	"bytes"
)

type State struct {
	Type     string
	Name     string
	Pass     *PassState
	Task     *TaskState
	Choice   *ChoiceState
	Wait     *WaitState
	Succeed  *SucceedState
	Fail     *FailState
	Parallel *ParallelState
	Map      *MapState
}

func (s State) Transition(r, w *bytes.Buffer) (next string, err error) {
	switch s.Type {
	case "Pass":
		return s.Pass.Transition(r, w)
	case "Task":
		return s.Task.Transition(r, w)
	case "Choice":
		return s.Choice.Transition(r, w)
	case "Wait":
		return s.Wait.Transition(r, w)
	case "Succeed":
		return s.Succeed.Transition(r, w)
	case "Fail":
		return s.Fail.Transition(r, w)
	case "Parallel":
		return s.Parallel.Transition(r, w)
	case "Map":
		return s.Map.Transition(r, w)
	}

	return "", ErrUnknownStateType
}
