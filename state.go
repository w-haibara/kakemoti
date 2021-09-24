package main

import (
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

func (s State) Transition() (next string, err error) {
	log.Println("State:", s.Name, "( Type =", s.Type, ")")

	switch s.Type {
	case "Pass":
		return s.Pass.Transition()
	case "Task":
		return s.Task.Transition()
	case "Choice":
		return s.Choice.Transition()
	case "Wait":
		return s.Wait.Transition()
	case "Succeed":
		return s.Succeed.Transition()
	case "Fail":
		return s.Fail.Transition()
	case "Parallel":
		return s.Parallel.Transition()
	case "Map":
		return s.Map.Transition()
	}

	return "", UnknownStateName
}
