package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/k0kubun/pp"
)

type CommonState struct {
	Type       string `json:"Type"`
	Next       string `json:"Next"`
	End        bool   `json:"End"`
	Comment    string `json:"Comment"`
	InputPath  string `json:"InputPath"`
	OutputPath string `json:"OutputPath"`
}

type PassState struct {
	CommonState
	Result     string `json:"Result"`
	ResultPath string `json:"ResultPath"`
	Parameters string `json:"Parameters"`
}

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

type ChoiceState struct {
	CommonState
	Choices []map[string]interface{} `json:"Choices"`
	Default string                   `json:"Default"`
}

type WaitState struct {
	CommonState
	Seconds       int64  `json:"Seconds"`
	Timestamp     string `json:"Timestamp"`
	SecondsPath   string `json:"SecondsPath"`
	TimestampPath string `json:"TimestampPath"`
}

type SucceedState struct {
}

type FailState struct {
	CommonState
	Cause string `json:"Cause"`
	Error string `json:"Error"`
}

type ParallelState struct {
	CommonState
	Branches       []StateMachine `json:"Branches"`
	ResultPath     string         `json:"ResultPath"`
	ResultSelector string         `json:"ResultSelector"`
	Retry          string         `json:"Retry"`
	Catch          string         `json:"Catch"`
}

type MapState struct {
	CommonState
	Iterator       StateMachine `json:"Iterator"`
	ItemsPath      string       `json:"ItemsPath"`
	MaxConcurrency int64        `json:"MaxConcurrency"`
	ResultPath     string       `json:"ResultPath"`
	ResultSelector string       `json:"ResultSelector"`
	Retry          string       `json:"Retry"`
	Catch          string       `json:"Catch"`
}

type State struct {
	Type     string
	Pass     *PassState
	Task     *TaskState
	Choice   *ChoiceState
	Wait     *WaitState
	Succeed  *SucceedState
	Fail     *FailState
	Parallel *ParallelState
	Map      *MapState
}

type States map[string]State

type StateMachine struct {
	Comment        string                 `json:"Comment"`
	StartAt        string                 `json:"StartAt"`
	TimeoutSeconds int64                  `json:"TimeoutSeconds"`
	Version        int64                  `json:"Version"`
	States         map[string]interface{} `json:"States"`
}

func main() {
	f, err := os.Open("./workflow.json")
	if err != nil {
		log.Panic("error:", err)
	}

	dec := json.NewDecoder(f)

	sm := StateMachine{}
	if err := dec.Decode(&sm); err != nil {
		log.Panic("error:", err)
	}

	sm.PrintInfo()

	states := sm.GenerateStates()
	states.Print()
}

func (sm StateMachine) PrintInfo() {
	fmt.Println("====== StateMachine Info ======")
	pp.Println("Comment", sm.Comment)
	pp.Println("StartAt", sm.StartAt)
	pp.Println("TimeoutSeconds", sm.TimeoutSeconds)
	pp.Println("Version", sm.Version)
	fmt.Println("===============================")
}

func (sm StateMachine) GenerateStates() States {
	states := map[string]State{}
	for name, state := range sm.States {
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
				Pass: &PassState{},
			}
			if err := convert(s, states[name].Pass); err != nil {
				log.Println("error:", err)
				continue
			}
		case "Task":
			states[name] = State{
				Type: "Task",
				Task: &TaskState{},
			}
			if err := convert(s, states[name].Task); err != nil {
				log.Println("error:", err)
				continue
			}
		case "Choice":
			states[name] = State{
				Type:   "Choice",
				Choice: &ChoiceState{},
			}
			if err := convert(s, states[name].Choice); err != nil {
				log.Println("error:", err)
				continue
			}
		case "Wait":
			states[name] = State{
				Type: "Wait",
				Wait: &WaitState{},
			}
			if err := convert(s, states[name].Wait); err != nil {
				log.Println("error:", err)
				continue
			}
		case "Succeed":
			states[name] = State{
				Type:    "Succeed",
				Succeed: &SucceedState{},
			}
			if err := convert(s, states[name].Succeed); err != nil {
				log.Println("error:", err)
				continue
			}
		case "Fail":
			states[name] = State{
				Type: "Fail",
				Fail: &FailState{},
			}
			if err := convert(s, states[name].Fail); err != nil {
				log.Println("error:", err)
				continue
			}
		case "Parallel":
			states[name] = State{
				Type:     "Parallel",
				Parallel: &ParallelState{},
			}
			if err := convert(s, states[name].Parallel); err != nil {
				log.Println("error:", err)
				continue
			}
		case "Map":
			states[name] = State{
				Type: "Map",
				Map:  &MapState{},
			}
			if err := convert(s, states[name].Map); err != nil {
				log.Println("error:", err)
				continue
			}
		default:
			states[name] = State{
				Type: t,
			}
		}
	}

	return states
}

func (s States) Print() {
	fmt.Println("=========== States  ===========")
	for k, v := range s {
		pp.Println(k)

		switch v.Type {
		case "Pass":
			pp.Println(v.Pass)
		case "Task":
			pp.Println(v.Task)
		case "Choice":
			pp.Println(v.Choice)
		case "Wait":
			pp.Println(v.Wait)
		case "Succeed":
			pp.Println(v.Succeed)
		case "Fail":
			pp.Println(v.Fail)
		case "Parallel":
			pp.Println(v.Parallel)
		case "Map":
			pp.Println(v.Map)
		}

		println()
	}
	fmt.Println("===============================")
}
