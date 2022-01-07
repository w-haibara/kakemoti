package worker

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
	return e.err.Error()
}

func (e statesError) IsEmpty() bool {
	return e.statesErr == "" && e.err == nil
}
