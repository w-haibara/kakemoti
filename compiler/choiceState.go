package compiler

import (
	"errors"
)

type RawChoiceState struct {
	CommonState2
	Choices []map[string]interface{} `json:"Choices"`
	Default string                   `json:"Default"`
}

func (raw RawChoiceState) decode() (*ChoiceState, error) {
	ErrNotFound := errors.New("key not found")
	ErrInvalidType := errors.New("invalid type")

	exist := func(m map[string]interface{}, k string) bool {
		_, ok := m[k]
		return ok
	}

	choices := make([]Choice, len(raw.Choices))
	for i, raw := range raw.Choices {
		n, ok := raw["Next"]
		if !ok {
			return nil, ErrNotFound
		}
		next, ok := n.(string)
		if !ok {
			return nil, ErrInvalidType
		}

		var1, ok := raw["Variable"]
		if !ok {
			return nil, ErrNotFound
		}

		var (
			op   string
			var2 interface{}
		)
		switch {
		case exist(raw, "BooleanEquals"):
			op = "BooleanEquals"
			var2 = raw[op]
		}

		choices[i] = Choice{
			Rule: &Rule{
				Variable1: var1,
				Variable2: var2,
				Operator:  op,
			},
			Next: next,
		}
	}

	return &ChoiceState{
		CommonState2: raw.CommonState2,
		Choices:      choices,
		Default:      raw.Default,
	}, nil
}

type ChoiceState struct {
	CommonState2
	Choices []Choice `json:"Choices"`
	Default string   `json:"Default"`
}

type Choice struct {
	Rule     *Rule
	BoolExpr BoolExpr
	Next     string
}

type BoolExpr map[string][]Choice

type Rule struct {
	Variable1 interface{}
	Variable2 interface{}
	Operator  string
}

func (state ChoiceState) GetNext() string {
	return ""
}
