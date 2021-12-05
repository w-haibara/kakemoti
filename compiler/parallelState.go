package compiler

import (
	"bytes"
	"encoding/json"
	"log"
)

type RawParallelState struct {
	CommonState
	Branches       []json.RawMessage `json:"Branches"`
	ResultPath     string            `json:"ResultPath"`
	ResultSelector string            `json:"ResultSelector"`
	Retry          string            `json:"Retry"`
	Catch          string            `json:"Catch"`
}

func (raw RawParallelState) decode() (*ParallelState, error) {
	branches := make([]Workflow, len(raw.Branches))

	for i, branch := range raw.Branches {
		b, err := json.Marshal(branch)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		buf := bytes.NewBuffer(b)
		asl, err := NewASL(buf)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		workflow, err := asl.compile()
		if err != nil {
			log.Println(err)
			return nil, err
		}

		branches[i] = *workflow
	}

	return &ParallelState{
		CommonState:    raw.CommonState,
		Branches:       branches,
		ResultPath:     raw.ResultPath,
		ResultSelector: raw.ResultSelector,
		Retry:          raw.Retry,
		Catch:          raw.Catch,
	}, nil
}

type ParallelState struct {
	CommonState
	Branches       []Workflow `json:"Branches"`
	ResultPath     string     `json:"ResultPath"`
	ResultSelector string     `json:"ResultSelector"`
	Retry          string     `json:"Retry"`
	Catch          string     `json:"Catch"`
}
