package statemachine

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/google/uuid"
	"github.com/spyzhov/ajson"
)

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

type resource struct {
	typ  string
	name string
}

func (s *TaskState) Transition(ctx context.Context, r *ajson.Node) (next string, w *ajson.Node, err error) {
	if s == nil {
		return "", nil, nil
	}

	select {
	case <-ctx.Done():
		return "", nil, ErrStoppedStateMachine
	default:
	}

	node, err := filterByInputPath(r, s.InputPath)
	if err != nil {
		return "", nil, err
	}

	res, err := s.parseResource()
	if err != nil {
		return "", nil, err
	}

	out, err := res.exec(ctx, node)
	if err != nil {
		// Task failed
		return "", nil, err
	}

	r, err = filterByResultPath(r, out, s.ResultPath)
	if err != nil {
		return "", nil, err
	}

	r, err = filterByOutputPath(r, s.OutputPath)
	if err != nil {
		return "", nil, err
	}

	if s.End {
		return "", r, ErrEndStateMachine
	}

	if strings.TrimSpace(s.Next) == "" {
		return "", nil, ErrNextStateIsBrank
	}

	return s.Next, r, nil
}

func (s *TaskState) parseResource() (*resource, error) {
	v := strings.Split(s.Resource, ":")

	if len(v) < 2 {
		return nil, fmt.Errorf("invalid resource")
	}

	if strings.Trim(v[0], "") == "" || strings.Trim(v[1], "") == "" {
		return nil, fmt.Errorf("invalid resource")
	}

	return &resource{
		typ:  v[0],
		name: v[1],
	}, nil
}

func (res *resource) exec(ctx context.Context, input *ajson.Node) (*ajson.Node, error) {
	switch res.typ {
	case "script":
		args := make([]string, 0)
		if !input.IsArray() {
			return nil, ErrInvalidInputPath
		}
		v := input.MustArray()
		for _, v := range v {
			if !v.IsString() {
				return nil, ErrInvalidInputPath
			}
			args = append(args, v.MustString())
		}

		out, err := res.execScript(ctx, args...)
		if err != nil {
			return nil, err
		}

		node := ajson.ObjectNode(uuid.New().String(), map[string]*ajson.Node{
			"result": ajson.StringNode(uuid.New().String(), string(out)),
		})

		return node, nil
	case "command":
		// TODO
	case "curl":
		// TODO
	}

	return nil, fmt.Errorf("invalid resource type")
}

func (res *resource) execScript(ctx context.Context, args ...string) ([]byte, error) {
	exe, err := exec.LookPath(res.name)
	if err != nil {
		return nil, err
	}

	cmd := exec.CommandContext(ctx, exe, args...) // #nosec G204
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return out, nil
}
