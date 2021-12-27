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

type ASL struct {
	Comment        string                     `json:"Comment"`
	StartAt        *string                    `json:"StartAt"`
	TimeoutSeconds *int64                     `json:"TimeoutSeconds"`
	Version        *string                    `json:"Version"`
	States         map[string]json.RawMessage `json:"States"`
}

type Workflow struct {
	Comment        string
	StartAt        string
	TimeoutSeconds int64
	Version        string
	States         States
}

type States []State

type State struct {
	Type    string
	Name    string
	Next    string
	Body    StateBody
	Choices map[string]States
}

func makeStateMachine(s *States, state State, states map[string]State) error {
	if state.Type == "Choice" {
		if err := setChoices(s, state, states); err != nil {
			return err
		}
		return nil
	}

	if state.Next == "" {
		return nil
	}

	cur, ok := states[state.Next]
	if !ok {
		return fmt.Errorf("Next state is not found: %s", cur.Next)
	}
	*s = append(*s, cur)
	return makeStateMachine(s, cur, states)
}

func setChoices(s *States, state State, states map[string]State) error {
	body, ok := state.Body.(*ChoiceState)
	if !ok {
		return fmt.Errorf("can't covert to type ChoiceState")
	}

	choices := make(map[string]States)
	for _, choice := range body.Choices {
		if choice.Next == "" {
			continue
		}

		state, ok := states[choice.Next]
		if !ok {
			return fmt.Errorf("Next state is not found: %s", choice.Next)
		}
		choices[state.Name] = States{state}

		s := choices[state.Name]
		if err := makeStateMachine(&s, state, states); err != nil {
			return err
		}
	}

	state.Choices = choices

	var s1 []State = *s
	s1[len(s1)-1] = state

	return nil
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
	workflow := &Workflow{
		Comment:        asl.Comment,
		StartAt:        *asl.StartAt,
		TimeoutSeconds: *asl.TimeoutSeconds,
		Version:        *asl.Version,
	}

	states, err := asl.makeStates()
	if err != nil {
		return nil, err
	}

	workflow.States = make([]State, 1, len(states))
	workflow.States[0] = states[workflow.StartAt]
	if err := makeStateMachine(&workflow.States, states[workflow.StartAt], states); err != nil {
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
