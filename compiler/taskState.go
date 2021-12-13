package compiler

import (
	"encoding/json"
	"fmt"
	"strings"
)

var (
	ErrInvalidTaskResource     = fmt.Errorf("invalid resource")
	ErrInvalidTaskResourceType = fmt.Errorf("invalid resource type")
)

type RawTaskState struct {
	CommonState
	RawResource          string           `json:"Resource"`
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

func (s *RawTaskState) decode() (*TaskState, error) {
	v := strings.SplitN(s.RawResource, ":", 2)

	if len(v) != 2 {
		return nil, ErrInvalidTaskResource
	}

	return &TaskState{
		RawTaskState: s,
		Resouce: TaskResouce{
			Type: v[0],
			Path: v[1],
		},
	}, nil
}

type TaskState struct {
	*RawTaskState
	Resouce TaskResouce
}

type TaskResouce struct {
	Type string
	Path string
}
