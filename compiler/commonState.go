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

func (state *CommonState1) DecodePath() error {
	return nil
}

type CommonState2 struct {
	CommonState1
	RawInputPath  *string `json:"InputPath"`
	InputPath     *Path
	RawOutputPath *string `json:"OutputPath"`
	OutputPath    *Path
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

func (state *CommonState2) DecodePath() error {
	if err := state.CommonState1.DecodePath(); err != nil {
		return err
	}

	if state.RawInputPath != nil {
		v1, err := NewPath(*state.RawInputPath)
		if err != nil {
			return err
		}
		state.InputPath = &v1
	}

	if state.RawOutputPath != nil {
		v2, err := NewPath(*state.RawOutputPath)
		if err != nil {
			return err
		}
		state.OutputPath = &v2
	}

	return nil
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

func (state *CommonState3) DecodePath() error {
	return state.CommonState2.DecodePath()
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

func (state *CommonState4) DecodePath() error {
	if err := state.CommonState3.DecodePath(); err != nil {
		return err
	}

	if state.RawResultPath != nil {
		v, err := NewReferencePath(*state.RawResultPath)
		if err != nil {
			return err
		}
		state.ResultPath = &v
	}

	return nil
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

func (state *CommonState5) DecodePath() error {
	return state.CommonState4.DecodePath()
}
