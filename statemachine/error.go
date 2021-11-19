package statemachine

import (
	"errors"
	"fmt"
)

type aslError struct {
	name  string
	cause error
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
	ErrStatesALL                    = newASLError(statesAll)
	ErrStatesHeartbeatTimeout       = newASLError(statesHeartbeatTimeout)
	ErrStatesTimeout                = newASLError(statesTimeout)
	ErrStatesTaskFailed             = newASLError(statesTaskFailed)
	ErrStatesPermissions            = newASLError(statesPermissions)
	ErrStatesResultPathMatchFailure = newASLError(statesResultPathMatchFailure)
	ErrStatesParameterPathFailure   = newASLError(statesParameterPathFailure)
	ErrStatesBranchFailed           = newASLError(statesBranchFailed)
	ErrStatesNoChoiceMatched        = newASLError(statesNoChoiceMatched)
	ErrStatesIntrinsicFailure       = newASLError(statesIntrinsicFailure)
)

func newASLError(name string) aslError {
	return aslError{
		name: name,
	}
}

func newASLErrorWithCause(name string, cause error) aslError {
	return aslError{
		name:  name,
		cause: cause,
	}
}

func (e aslError) Error() string {
	if e.cause == nil {
		return fmt.Sprintf("ASL Error [%s]", e.name)
	}

	return fmt.Sprintf("ASL Error [%s], %v", e.name, e.cause)
}

func (e aslError) Unwrap() error {
	return e.cause
}

func (e aslError) Is(target error) bool {
	switch v := target.(type) {
	case aslError:
		if v.name == statesAll {
			return true
		}

		return e.name == v.name
	}

	return errors.Is(e, errors.Unwrap(target))
}
