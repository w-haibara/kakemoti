package statemachine

import (
	"errors"
)

type statemachineError struct {
	name string
}

const (
	statesAll                    = "States.ALL"
	statesHeartbeatTimeout       = "States.HeartbeatTimeout"
	statesTimeout                = "States.Timeout"
	statesTaskFailed             = "States.TaskFailed"
	statesPermissions            = "States.Permissions"
	statesResultPathMatchFailure = "States.ResultPathMatchFailure"
	statesParameterPathFailure   = "States.ParameterPathFailure"
	statesBranchFailed           = "States.BranchFailed"
	statesNoChoiceMatched        = "States.NoChoiceMatched"
	statesIntrinsicFailure       = "States.IntrinsicFailure"
)

var (
	ErrStatesALL                    = newStateMachineError(statesAll)
	ErrStatesHeartbeatTimeout       = newStateMachineError(statesHeartbeatTimeout)
	ErrStatesTimeout                = newStateMachineError(statesTimeout)
	ErrStatesTaskFailed             = newStateMachineError(statesTaskFailed)
	ErrStatesPermissions            = newStateMachineError(statesPermissions)
	ErrStatesResultPathMatchFailure = newStateMachineError(statesResultPathMatchFailure)
	ErrStatesParameterPathFailure   = newStateMachineError(statesParameterPathFailure)
	ErrStatesBranchFailed           = newStateMachineError(statesBranchFailed)
	ErrStatesNoChoiceMatched        = newStateMachineError(statesNoChoiceMatched)
	ErrStatesIntrinsicFailure       = newStateMachineError(statesIntrinsicFailure)
)

func newStateMachineError(name string) statemachineError {
	return statemachineError{
		name: name,
	}
}

func (e statemachineError) Error() string {
	return e.name
}

func (e statemachineError) Is(target error) bool {
	switch v := target.(type) {
	case statemachineError:
		if v.name == statesAll {
			return true
		}

		return e.name == v.name
	}

	return errors.Is(e, errors.Unwrap(target))
}
