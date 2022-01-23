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
	RawResource          string  `json:"Resource"`
	TimeoutSeconds       *int    `json:"TimeoutSeconds"`
	TimeoutSecondsPath   *string `json:"TimeoutSecondsPath"`
	HeartbeatSeconds     *int    `json:"HeartbeatSeconds"`
	HeartbeatSecondsPath *string `json:"HeartbeatSecondsPath"`
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

	var timeoutSecondsPath *Path
	if raw.TimeoutSecondsPath != nil {
		v, err := NewPath(*raw.TimeoutSecondsPath)
		if err != nil {
			return nil, err
		}
		timeoutSecondsPath = &v
	}

	var heartbeatSecondsPath *Path
	if raw.HeartbeatSecondsPath != nil {
		v, err := NewPath(*raw.HeartbeatSecondsPath)
		if err != nil {
			return nil, err
		}
		heartbeatSecondsPath = &v
	}

	return TaskState{
		CommonState5: s.Common(),
		Resouce: TaskResouce{
			Type: v[0],
			Path: v[1],
		},
		TimeoutSeconds:       raw.TimeoutSeconds,
		TimeoutSecondsPath:   timeoutSecondsPath,
		HeartbeatSeconds:     raw.HeartbeatSeconds,
		HeartbeatSecondsPath: heartbeatSecondsPath,
	}, nil
}

type TaskState struct {
	CommonState5
	Resouce              TaskResouce
	TimeoutSeconds       *int
	TimeoutSecondsPath   *Path
	HeartbeatSeconds     *int
	HeartbeatSecondsPath *Path
}

type TaskResouce struct {
	Type string
	Path string
}
