package main

import ()

type CommonState struct {
	Type       string `json:"Type"`
	Next       string `json:"Next"`
	End        bool   `json:"End"`
	Comment    string `json:"Comment"`
	InputPath  string `json:"InputPath"`
	OutputPath string `json:"OutputPath"`
}

func (s CommonState) Transition() (next string, err error) {
	if s.End {
		return "", EndStateMachine
	}

	return s.Next, nil
}
