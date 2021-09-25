package main

import (
	"io"
	"log"
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

func (s State) Transition(r io.Reader, w io.Writer) (next string, err error) {
	log.Println("State:", s.Name, "( Type =", s.Type, ")")

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

	return "", UnknownStateType
}
