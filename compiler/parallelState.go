package compiler

import (
	"bytes"
	"encoding/json"
	"log"
)

type RawParallelState struct {
	CommonState5
	Branches []json.RawMessage `json:"Branches"`
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
		CommonState5: raw.CommonState5,
		Branches:     branches,
	}, nil
}

type ParallelState struct {
	CommonState5
	Branches []Workflow `json:"Branches"`
}
