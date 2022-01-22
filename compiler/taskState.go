package compiler

import (
	"fmt"
	"strings"
)

var (
	ErrInvalidTaskResource     = fmt.Errorf("invalid resource")
	ErrInvalidTaskResourceType = fmt.Errorf("invalid resource type")
)

type RawTaskState struct {
	CommonState5
	RawResource          string `json:"Resource"`
	TimeoutSeconds       string `json:"TimeoutSeconds"`       // TODO
	TimeoutSecondsPath   string `json:"TimeoutSecondsPath"`   // TODO
	HeartbeatSeconds     string `json:"HeartbeatSeconds"`     // TODO
	HeartbeatSecondsPath string `json:"HeartbeatSecondsPath"` // TODO
}

func (raw *RawTaskState) decode(name string) (State, error) {
	s, err := raw.CommonState5.decode(name)
	if err != nil {
		return nil, err
	}
	raw.CommonState5 = s.Common()

	v := strings.SplitN(raw.RawResource, ":", 2)

	if len(v) != 2 {
		return nil, ErrInvalidTaskResource
	}

	return TaskState{
		RawTaskState: raw,
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
