package compiler

import "encoding/json"

type CommonState1 struct {
	Type    string `json:"Type"`
	Comment string `json:"Comment"`
}

func (state CommonState1) GetNext() string {
	return ""
}

type CommonState2 struct {
	CommonState1
	InputPath  string `json:"InputPath"`
	OutputPath string `json:"OutputPath"`
}

type CommonState3 struct {
	CommonState2
	Next string `json:"Next"`
	End  bool   `json:"End"`
}

func (state CommonState3) GetNext() string {
	return state.Next
}

type CommonState4 struct {
	CommonState3
	ResultPath string           `json:"ResultPath"`
	Parameters *json.RawMessage `json:"Parameters"`
}

type CommonState5 struct {
	CommonState4
	ResultSelector *json.RawMessage `json:"ResultSelector"`
	Retry          string           `json:"Retry"`
	Catch          string           `json:"Catch"`
}
