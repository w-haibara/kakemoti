package statemachine

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/google/uuid"
	"github.com/spyzhov/ajson"
)

var (
	ErrInvalidTaskResource     = fmt.Errorf("invalid resource")
	ErrInvalidTaskResourceType = fmt.Errorf("invalid resource type")
	ErrInvalidTaskInput        = fmt.Errorf("invalid task input")
)

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

	node, err = filterByParameters(node, s.Parameters)
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

	out, err = filterByResultSelector(out, s.ResultSelector)
	if err != nil {
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

	switch {
	case len(v) < 2:
		return nil, ErrInvalidTaskResource
	case len(v) > 2:
		v[1] = strings.Join(v[1:], "")
	}

	if strings.Trim(v[0], "") == "" || strings.Trim(v[1], "") == "" {
		return nil, ErrInvalidTaskResource
	}

	return &resource{
		typ:  v[0],
		name: v[1],
	}, nil
}

func (res *resource) exec(ctx context.Context, input *ajson.Node) (*ajson.Node, error) {
	switch res.typ {
	case "script":
		if !input.IsObject() {
			return nil, ErrInvalidTaskInput
		}
		input = input.MustObject()["args"]

		args := make([]string, 0)
		if !input.IsArray() {
			return nil, ErrInvalidTaskInput
		}
		for _, v := range input.MustArray() {
			if !v.IsString() {
				return nil, ErrInvalidTaskInput
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
		panic("not implemented: Task state, resource type = command")
	case "curl":
		// TODO
		panic("not implemented: Task state, resource type = curl")
	}

	return nil, ErrInvalidTaskResourceType
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
