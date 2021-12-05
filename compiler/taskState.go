package compiler

import "encoding/json"

type TaskState struct {
	CommonState
	Resource             string           `json:"Resource"`
	Parameters           *json.RawMessage `json:"Parameters"`
	ResultPath           string           `json:"ResultPath"`
	ResultSelector       *json.RawMessage `json:"ResultSelector"`
	Retry                string           `json:"Retry"`                // TODO
	Catch                string           `json:"Catch"`                // TODO
	TimeoutSeconds       string           `json:"TimeoutSeconds"`       // TODO
	TimeoutSecondsPath   string           `json:"TimeoutSecondsPath"`   // TODO
	HeartbeatSeconds     string           `json:"HeartbeatSeconds"`     // TODO
	HeartbeatSecondsPath string           `json:"HeartbeatSecondsPath"` // TODO

}
