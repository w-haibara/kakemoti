package worker

type statesError string

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
