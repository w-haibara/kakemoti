package main

import (
	"io"
	"strings"

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

func (s CommonState) Transition(r io.Reader, w io.Writer) (next string, err error) {
	if s.End {
		return "", EndStateMachine
	}

	if strings.TrimSpace(s.Next) == "" {
		return "", NextStateIsBrank
	}

	return s.Next, nil
}

func (s CommonState) Print() {
	pp.Println(s)
}
