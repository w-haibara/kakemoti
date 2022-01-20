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
	StateName string
	Type      string `json:"Type"`
	Comment   string `json:"Comment"`
}

func (state CommonState1) decode(name string) (State, error) {
	state.StateName = name
	return state, nil
}

func (state CommonState1) Name() string {
	return state.StateName
}

func (state CommonState1) Next() string {
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
	RawInputPath  *string `json:"InputPath"`
	InputPath     *Path
	RawOutputPath *string `json:"OutputPath"`
	OutputPath    *Path
}

func (state CommonState2) decode(name string) (State, error) {
	s, err := state.CommonState1.decode(name)
	if err != nil {
		return nil, err
	}

	res := CommonState2{CommonState1: s.Common().CommonState1}

	if state.RawInputPath != nil {
		v1, err := NewPath(*state.RawInputPath)
		if err != nil {
			return nil, err
		}
		res.InputPath = &v1
	}

	if state.RawOutputPath != nil {
		v2, err := NewPath(*state.RawOutputPath)
		if err != nil {
			return nil, err
		}
		res.OutputPath = &v2
	}

	return res, nil
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
	NextName string `json:"Next"`
	End      bool   `json:"End"`
}

func (state CommonState3) Next() string {
	return state.NextName
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

func (state CommonState3) decode(name string) (State, error) {
	s, err := state.CommonState2.decode(name)
	if err != nil {
		return nil, err
	}
	state.CommonState2 = s.Common().CommonState2
	return state, nil
}

type CommonState4 struct {
	CommonState3
	RawResultPath *string `json:"ResultPath"`
	ResultPath    *ReferencePath
	Parameters    *json.RawMessage `json:"Parameters"`
}

func (state CommonState4) FieldsType() int {
	return FieldsType4
}

func (state CommonState4) Common() CommonState5 {
	return CommonState5{
		CommonState4: state,
	}
}

func (state CommonState4) decode(name string) (State, error) {
	s, err := state.CommonState3.decode(name)
	if err != nil {
		return nil, err
	}
	state.CommonState3 = s.Common().CommonState3

	if state.RawResultPath != nil {
		v, err := NewReferencePath(*state.RawResultPath)
		if err != nil {
			return nil, err
		}
		state.ResultPath = &v
	}

	return state, nil
}

type CommonState5 struct {
	CommonState4
	ResultSelector *json.RawMessage `json:"ResultSelector"`
	Retry          []Retry          `json:"Retry"`
	Catch          []Catch          `json:"Catch"`
}

func (state CommonState5) FieldsType() int {
	return FieldsType5
}

type Retry struct {
	ErrorEquals     []string
	IntervalSeconds *int
	MaxAttempts     *int
	BackoffRate     *float64
}

type Catch struct {
	ErrorEquals []string
	ResultPath  *Path
	Next        string
}

func (state CommonState5) Common() CommonState5 {
	return state
}

func (state CommonState5) decode(name string) (State, error) {
	s, err := state.CommonState4.decode(name)
	if err != nil {
		return nil, err
	}
	state.CommonState4 = s.Common().CommonState4
	return state, nil
}
