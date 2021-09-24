package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/k0kubun/pp"
)

var (
	InvalidStartAtValue = fmt.Errorf("invalid StateAt value")
	UnknownStateName    = fmt.Errorf("unknown state name")
	UnknownStateType    = fmt.Errorf("unknown state type")
	NextStateIsBrank    = fmt.Errorf("next state is brank")
	EndStateMachine     = fmt.Errorf("end state machine")
)

type StateMachine struct {
	Comment        string                 `json:"Comment"`
	StartAt        string                 `json:"StartAt"`
	TimeoutSeconds int64                  `json:"TimeoutSeconds"`
	Version        int64                  `json:"Version"`
	RawStates      map[string]interface{} `json:"States"`
	States         map[string]State       `json:"-"`
}

func NewStateMachine(path string) (*StateMachine, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(f)

	sm := new(StateMachine)
	if err := dec.Decode(sm); err != nil {
		return nil, err
	}

	sm.SetStates()

	return sm, nil
}

func (sm *StateMachine) SetStates() {
	states := map[string]State{}
	for name, state := range sm.RawStates {
		s, ok := state.(map[string]interface{})
		if !ok {
			log.Println("invalid state definition:", name)
			continue
		}

		t, ok := s["Type"].(string)
		if !ok {
			log.Println("invalid type value:", s["Type"])
			continue
		}

		convert := func(src, dst interface{}) error {
			var buf bytes.Buffer
			enc := json.NewEncoder(&buf)
			if err := enc.Encode(src); err != nil {
				return err
			}

			dec := json.NewDecoder(&buf)
			if err := dec.Decode(&dst); err != nil {
				return err
			}

			return nil
		}

		switch t {
		case "Pass":
			states[name] = State{
				Type: "Pass",
				Name: name,
				Pass: &PassState{},
			}
			if err := convert(s, states[name].Pass); err != nil {
				log.Println("error:", err)
				continue
			}
		case "Task":
			states[name] = State{
				Type: "Task",
				Name: name,
				Task: &TaskState{},
			}
			if err := convert(s, states[name].Task); err != nil {
				log.Println("error:", err)
				continue
			}
		case "Choice":
			states[name] = State{
				Type:   "Choice",
				Name:   name,
				Choice: &ChoiceState{},
			}
			if err := convert(s, states[name].Choice); err != nil {
				log.Println("error:", err)
				continue
			}
		case "Wait":
			states[name] = State{
				Type: "Wait",
				Name: name,
				Wait: &WaitState{},
			}
			if err := convert(s, states[name].Wait); err != nil {
				log.Println("error:", err)
				continue
			}
		case "Succeed":
			states[name] = State{
				Type:    "Succeed",
				Name:    name,
				Succeed: &SucceedState{},
			}
			if err := convert(s, states[name].Succeed); err != nil {
				log.Println("error:", err)
				continue
			}
		case "Fail":
			states[name] = State{
				Type: "Fail",
				Name: name,
				Fail: &FailState{},
			}
			if err := convert(s, states[name].Fail); err != nil {
				log.Println("error:", err)
				continue
			}
		case "Parallel":
			states[name] = State{
				Type:     "Parallel",
				Name:     name,
				Parallel: &ParallelState{},
			}
			if err := convert(s, states[name].Parallel); err != nil {
				log.Println("error:", err)
				continue
			}
			for i := range states[name].Parallel.Branches {
				states[name].Parallel.Branches[i].SetStates()
			}
		case "Map":
			states[name] = State{
				Type: "Map",
				Name: name,
				Map:  &MapState{},
			}
			if err := convert(s, states[name].Map); err != nil {
				log.Println("error:", err)
				continue
			}
			states[name].Map.Iterator.SetStates()
		default:
			states[name] = State{
				Type: t,
			}
		}
	}

	sm.States = states
}

func (sm StateMachine) PrintInfo() {
	fmt.Println("====== StateMachine Info ======")
	pp.Println("Comment", sm.Comment)
	pp.Println("StartAt", sm.StartAt)
	pp.Println("TimeoutSeconds", sm.TimeoutSeconds)
	pp.Println("Version", sm.Version)
	fmt.Println("===============================")
}

func (sm StateMachine) PrintStates() {
	s := sm.States
	fmt.Println("=========== States  ===========")
	for k, v := range s {
		pp.Println(k)

		switch v.Type {
		case "Pass":
			v.Pass.Print()
		case "Task":
			v.Task.Print()
		case "Choice":
			v.Choice.Print()
		case "Wait":
			v.Wait.Print()
		case "Succeed":
			v.Succeed.Print()
		case "Fail":
			v.Fail.Print()
		case "Parallel":
			v.Parallel.Print()
		case "Map":
			v.Map.Print()
		}

		println()
	}
	fmt.Println("===============================")
}

func (sm StateMachine) Start() error {
	if _, ok := sm.States[sm.StartAt]; !ok {
		return InvalidStartAtValue
	}

	next := sm.StartAt
	var err error
	for {
		s, ok := sm.States[next]
		if !ok {
			log.Println("UnknownStateName:", next)
			return UnknownStateName
		}

		next, err = s.Transition()
		switch {
		case err == UnknownStateType:
			log.Println("UnknownStateType:", next)
			goto End
		case err == EndStateMachine:
			log.Println("EndStateMachine")
			goto End
		case err != nil:
			return err
		}

		if _, ok := sm.States[next]; !ok {
			log.Println("UnknownStateName: [", next, "]")
			goto End
		}
	}

End:
	return nil
}
