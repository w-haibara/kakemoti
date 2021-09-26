package statemachine

type TaskState struct {
	CommonState
	Resource             string `json:"Resource"`
	Parameters           string `json:"Parameters"`
	ResultPath           string `json:"ResultPath"`
	ResultSelector       string `json:"ResultSelector"`
	Retry                string `json:"Retry"`
	Catch                string `json:"Catch"`
	TimeoutSeconds       string `json:"TimeoutSeconds"`
	TimeoutSecondsPath   string `json:"TimeoutSecondsPath"`
	HeartbeatSeconds     string `json:"HeartbeatSeconds"`
	HeartbeatSecondsPath string `json:"HeartbeatSecondsPath"`
}
