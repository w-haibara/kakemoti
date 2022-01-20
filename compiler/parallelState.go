package compiler

import (
	"log"
)

type RawParallelState struct {
	CommonState5
	Branches []ASL `json:"Branches"`
}

func (raw RawParallelState) decode(name string) (State, error) {
	s, err := raw.CommonState5.decode(name)
	if err != nil {
		return nil, err
	}

	branches := make([]Workflow, len(raw.Branches))

	for i, branch := range raw.Branches {
		workflow, err := branch.compile()
		if err != nil {
			log.Println(err)
			return nil, err
		}

		branches[i] = *workflow
	}

	return ParallelState{
		CommonState5: s.Common(),
		Branches:     branches,
	}, nil
}

type ParallelState struct {
	CommonState5
	Branches []Workflow
}
