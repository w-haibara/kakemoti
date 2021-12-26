package compiler

import (
	"encoding/json"
)

const (
	FieldsType1 = iota
	FieldsType2
	FieldsType3
	FieldsType4
	FieldsType5
)

type CommonState1 struct {
	Type    string `json:"Type"`
	Comment string `json:"Comment"`
}

func (state CommonState1) GetNext() string {
	return ""
}

func (state CommonState1) FieldsType() int {
	return FieldsType1
}

func (state CommonState1) Common() CommonState5 {
	return CommonState5{
		CommonState4: CommonState4{
			CommonState3: CommonState3{
				CommonState2: CommonState2{
					CommonState1: state,
				},
			},
		},
	}
}

type CommonState2 struct {
	CommonState1
	InputPath  string `json:"InputPath"`
	OutputPath string `json:"OutputPath"`
}

func (state CommonState2) FieldsType() int {
	return FieldsType2
}

func (state CommonState2) Common() CommonState5 {
	return CommonState5{
		CommonState4: CommonState4{
			CommonState3: CommonState3{
				CommonState2: state,
			},
		},
	}
}

type CommonState3 struct {
	CommonState2
	Next string `json:"Next"`
	End  bool   `json:"End"`
}

func (state CommonState3) GetNext() string {
	return state.Next
}

func (state CommonState3) FieldsType() int {
	return FieldsType3
}

func (state CommonState3) Common() CommonState5 {
	return CommonState5{
		CommonState4: CommonState4{
			CommonState3: state,
		},
	}
}

type CommonState4 struct {
	CommonState3
	ResultPath string           `json:"ResultPath"`
	Parameters *json.RawMessage `json:"Parameters"`
}

func (state CommonState4) FieldsType() int {
	return FieldsType4
}

func (state CommonState4) Common() CommonState5 {
	return CommonState5{
		CommonState4: state,
	}
}

type CommonState5 struct {
	CommonState4
	ResultSelector *json.RawMessage `json:"ResultSelector"`
	Retry          string           `json:"Retry"`
	Catch          string           `json:"Catch"`
}

func (state CommonState5) FieldsType() int {
	return FieldsType5
}

func (state CommonState5) Common() CommonState5 {
	return state
}
