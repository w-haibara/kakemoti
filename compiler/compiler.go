package compiler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

var ErrStateMachineTerminated = errors.New("state machine terminated")

type States []State

type State struct {
	Type string
	Name string
	Next string
	Body StateBody
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func Compile(ctx context.Context, aslBytes *bytes.Buffer) (*Workflow, error) {
	asl, err := NewASL(aslBytes)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	workflow, err := asl.compile()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return workflow, nil
}

type ASL struct {
	Comment        string                     `json:"Comment"`
	StartAt        *string                    `json:"StartAt"`
	TimeoutSeconds *int64                     `json:"TimeoutSeconds"`
	Version        *string                    `json:"Version"`
	States         map[string]json.RawMessage `json:"States"`
}

func NewASL(aslBytes *bytes.Buffer) (*ASL, error) {
	dec := json.NewDecoder(aslBytes)

	asl := new(ASL)
	if err := dec.Decode(asl); err != nil {
		return nil, err
	}

	if err := asl.validate(); err != nil {
		return nil, err
	}

	return asl, nil
}

func (asl *ASL) validate() error {
	if asl.StartAt == nil {
		return fmt.Errorf("Top-level fields: 'StartAt' is needed")
	}

	if asl.Version == nil {
		asl.Version = new(string)
		*asl.Version = "1.0"
	}

	if asl.TimeoutSeconds == nil {
		asl.TimeoutSeconds = new(int64)
		*asl.TimeoutSeconds = 0
	}

	return nil
}

func (asl *ASL) compile() (*Workflow, error) {
	states, err := asl.makeStates()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for k := range states {
		if err := states[k].Body.DecodePath(); err != nil {
			return nil, err
		}
	}

	workflow, err := asl.makeWorkflow(states)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return workflow, nil
}

func (asl *ASL) makeStates() (map[string]State, error) {
	states := make(map[string]State)
	for name, state := range asl.States {
		v := &struct {
			Type string `json:"type"`
		}{}
		if err := json.Unmarshal(state, v); err != nil {
			log.Println(err)
			return nil, err
		}

		switch v.Type {
		case "Choice":
			var raw RawChoiceState
			if err := json.Unmarshal(state, &raw); err != nil {
				log.Println(err)
				return nil, err
			}

			body, err := raw.decode()
			if err != nil {
				log.Println(err)
				return nil, err
			}

			states[name] = State{
				Type: v.Type,
				Name: name,
				Body: body,
			}
			continue
		case "Parallel":
			var raw RawParallelState
			if err := json.Unmarshal(state, &raw); err != nil {
				log.Println(err)
				return nil, err
			}

			body, err := raw.decode()
			if err != nil {
				log.Println(err)
				return nil, err
			}

			states[name] = State{
				Type: v.Type,
				Name: name,
				Next: body.Next,
				Body: body,
			}
			continue
		case "Task":
			var raw RawTaskState
			if err := json.Unmarshal(state, &raw); err != nil {
				log.Println(err)
				return nil, err
			}

			body, err := raw.decode()
			if err != nil {
				log.Println(err)
				return nil, err
			}

			states[name] = State{
				Type: v.Type,
				Name: name,
				Next: body.Next,
				Body: body,
			}
			continue
		}

		var body StateBody
		switch v.Type {
		case "Pass":
			body = new(PassState)
		case "Task":
			body = new(TaskState)
		case "Wait":
			body = new(WaitState)
		case "Succeed":
			body = new(SucceedState)
		case "Fail":
			body = new(FailState)
		case "Map":
			body = new(MapState)
		default:
			err := fmt.Errorf("Unknown state name: %s", v.Type)
			log.Println(err)
			return nil, err
		}
		if err := json.Unmarshal(state, body); err != nil {
			log.Println(err)
			return nil, err
		}

		states[name] = State{
			Type: v.Type,
			Name: name,
			Next: body.GetNext(),
			Body: body,
		}
	}

	return states, nil
}

func (asl *ASL) makeWorkflow(statesMap map[string]State) (*Workflow, error) {
	workflow := NewWorkflow(*asl)
	if err := workflow.makeStates(statesMap); err != nil {
		return nil, err
	}

	return workflow, nil
}

type Workflow struct {
	Comment        string
	StartAt        string
	TimeoutSeconds int64
	Version        string
	States         []States
	StatesIndexMap map[string][2]int
}

func NewWorkflow(asl ASL) *Workflow {
	m := make(map[string][2]int)
	var (
		startAt              = ""
		timeoutSeconds int64 = 0
		version              = ""
	)
	if asl.StartAt != nil {
		startAt = *asl.StartAt
	}
	if asl.TimeoutSeconds != nil {
		timeoutSeconds = *asl.TimeoutSeconds
	}
	if asl.Version != nil {
		version = *asl.Version
	}
	return &Workflow{
		Comment:        asl.Comment,
		StartAt:        startAt,
		TimeoutSeconds: timeoutSeconds,
		Version:        version,
		StatesIndexMap: m,
	}
}

func (wf *Workflow) makeStates(statesMap map[string]State) error {
	nexts := []string{wf.StartAt}
	for {
		ns, err := wf.makeBranches(nexts, statesMap)
		if err != nil {
			log.Println(err)
			return err
		}

		if len(ns) == 0 {
			return nil
		}

		nexts = ns
	}
}

func (wf *Workflow) makeBranches(starts []string, statesMap map[string]State) ([]string, error) {
	if len(starts) == 0 {
		return []string{}, nil
	}

	nexts := []string{}
	for _, next := range starts {
		if next == "" {
			continue
		}

		if _, ok := statesMap[next]; !ok {
			return nil, fmt.Errorf("invalid state name: %s", next)
		}

		ns1, err := wf.makeBranch(statesMap[next], statesMap)
		if err != nil {
			return nil, err
		}

		ns2, err := wf.makeCatchBranch(statesMap[next], statesMap)
		if err != nil {
			return nil, err
		}

		nexts = append(nexts, ns1...)
		nexts = append(nexts, ns2...)
	}

	return nexts, nil
}

func (wf *Workflow) makeBranch(start State, statesMap map[string]State) ([]string, error) {
	states := make(States, 0)
	cur := start
	for {
		states = append(states, cur)
		if cur.Next == "" {
			if wf.stateIsExistInBranch(cur.Name) {
				return nil, nil
			}
			wf.States = append(wf.States, states)
			for i, state := range states {
				wf.StatesIndexMap[state.Name] = [2]int{len(wf.States) - 1, i}
			}
			if bn := GetNexts(cur.Body); bn != nil {
				nexts := make([]string, 0, len(bn))
				for _, next := range bn {
					if _, ok := wf.StatesIndexMap[next]; !ok {
						nexts = append(nexts, next)
					}
				}
				return nexts, nil
			}
			return nil, nil
		} else if wf.stateIsExistInBranch(cur.Next) {
			wf.States = append(wf.States, states)
			for i, state := range states {
				wf.StatesIndexMap[state.Name] = [2]int{len(wf.States) - 1, i}
			}
			return nil, nil
		}
		var ok bool
		cur, ok = statesMap[cur.Next]
		if !ok {
			return nil, fmt.Errorf("key not found: %v", cur.Next)
		}
	}
}

func (wf *Workflow) makeCatchBranch(state State, statesMap map[string]State) ([]string, error) {
	nexts := []string{}
	if state.Body.FieldsType() >= FieldsType5 {
		for _, catch := range state.Body.Common().Catch {
			if _, ok := wf.StatesIndexMap[catch.Next]; !ok {
				nexts = append(nexts, catch.Next)
			}
		}
	}
	return nexts, nil
}

func (wf *Workflow) stateIsExistInBranch(name string) bool {
	i := wf.StatesIndexMap[name]
	if len(wf.States) > i[0] && len(wf.States[i[0]]) > i[1] {
		if wf.States[i[0]][i[1]].Name == name {
			return true
		}
	}
	return false
}

func GetNexts(body StateBody) []string {
	if choice, ok := body.(*ChoiceState); ok {
		return choice.GetNexts()
	}

	return nil
}
