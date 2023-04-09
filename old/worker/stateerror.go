package worker

import (
	"errors"
	"fmt"
)

const (
	StatesErrorALL                    = "States.ALL"
	StatesErrorHeartbeatTimeout       = "States.HeartbeatTimeout"
	StatesErrorTimeout                = "States.Timeout"
	StatesErrorTaskFailed             = "States.TaskFailed"
	StatesErrorPermissions            = "States.Permissions"
	StatesErrorResultPathMatchFailure = "States.ResultPathMatchFailure"
	StatesErrorParameterPathFailure   = "States.ParameterPathFailure"
	StatesErrorBranchFailed           = "States.BranchFailed"
	StatesErrorNoChoiceMatched        = "States.NoChoiceMatched"
	StatesErrorIntrinsicFailure       = "States.IntrinsicFailure"
)

type statesError struct {
	statesErr string
	err       error
}

func NewStatesError(statesErr string, err error) statesError {
	return statesError{
		statesErr: statesErr,
		err:       err,
	}
}

func (e statesError) StatesError() string {
	return e.statesErr
}

func (e statesError) Error() string {
	if e.IsEmpty() {
		return "nil"
	}
	s1 := e.statesErr
	s2 := "nil"
	if e.err != nil {
		s2 = e.err.Error()
	}
	return fmt.Sprintf("StateError=[%s], Error=[%s]", s1, s2)
}

func (e *statesError) IsEmpty() bool {
	if e == nil {
		return true
	}
	return e.statesErr == "" && e.err == nil
}

func (e statesError) Is(target error) bool {
	return errors.Is(e.err, target)
}

func (e statesError) As(target interface{}) bool {
	return errors.As(e.err, target)
}
