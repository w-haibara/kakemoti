package compiler

import (
	"context"
	"errors"
	"fmt"
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

	decodeCond := func(m map[string]interface{}) (Condition, error) {
		if !exist(m, "Variable") {
			return nil, ErrNotFound
		}
		v, ok := m["Variable"].(string)
		if !ok {
			return nil, ErrInvalidType
		}
		v1, err := NewPath(v)
		if err != nil {
			return nil, err
		}

		switch {
		case exist(m, "BooleanEquals"):
			v2, ok := m["BooleanEquals"].(bool)
			if !ok {
				return nil, ErrInvalidType
			}

			return BooleanEqualsRule{v1, v2}, nil
		case exist(m, "BooleanEqualsPath"):
			panic("Not Implemented")
		case exist(m, "IsBoolean"):
			panic("Not Implemented")
		case exist(m, "IsNull"):
			panic("Not Implemented")
		case exist(m, "IsNumeric"):
			panic("Not Implemented")
		case exist(m, "IsPresent"):
			panic("Not Implemented")
		case exist(m, "IsString"):
			panic("Not Implemented")
		case exist(m, "IsTimestamp"):
			panic("Not Implemented")
		case exist(m, "NumericEquals"):
			panic("Not Implemented")
		case exist(m, "NumericEqualsPath"):
			panic("Not Implemented")
		case exist(m, "NumericGreaterThan"):
			panic("Not Implemented")
		case exist(m, "NumericGreaterThanPath"):
			panic("Not Implemented")
		case exist(m, "NumericGreaterThanEquals"):
			panic("Not Implemented")
		case exist(m, "NumericGreaterThanEqualsPath"):
			panic("Not Implemented")
		case exist(m, "NumericLessThan"):
			panic("Not Implemented")
		case exist(m, "NumericLessThanPath"):
			panic("Not Implemented")
		case exist(m, "NumericLessThanEquals"):
			panic("Not Implemented")
		case exist(m, "NumericLessThanEqualsPath"):
			panic("Not Implemented")
		case exist(m, "StringEquals"):
			panic("Not Implemented")
		case exist(m, "StringEqualsPath"):
			panic("Not Implemented")
		case exist(m, "StringGreaterThan"):
			panic("Not Implemented")
		case exist(m, "StringGreaterThanPath"):
			panic("Not Implemented")
		case exist(m, "StringGreaterThanEquals"):
			panic("Not Implemented")
		case exist(m, "StringGreaterThanEqualsPath"):
			panic("Not Implemented")
		case exist(m, "StringLessThanStringLessThanPath"):
			panic("Not Implemented")
		case exist(m, "StringLessThanEqualsStringLessThanEqualsPath"):
			panic("Not Implemented")
		case exist(m, "StringMatches"):
			panic("Not Implemented")
		case exist(m, "TimestampEquals"):
			panic("Not Implemented")
		case exist(m, "TimestampEqualsPath"):
			panic("Not Implemented")
		case exist(m, "TimestampGreaterThan"):
			panic("Not Implemented")
		case exist(m, "TimestampGreaterThanPath"):
			panic("Not Implemented")
		case exist(m, "TimestampGreaterThanEquals"):
			panic("Not Implemented")
		case exist(m, "TimestampGreaterThanEqualsPath"):
			panic("Not Implemented")
		case exist(m, "TimestampLessThan"):
			panic("Not Implemented")
		case exist(m, "TimestampLessThanPath"):
			panic("Not Implemented")
		case exist(m, "TimestampLessThanEquals"):
			panic("Not Implemented")
		case exist(m, "TimestampLessThanEqualsPath"):
			panic("Not Implemented")
		}

		return nil, nil
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

		cond, err := decodeCond(raw)
		if err != nil {
			return nil, err
		}

		choices[i] = Choice{
			Condition: cond,
			Next:      next,
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
	Choices []Choice
	Default string
}

func (state ChoiceState) GetNext() string {
	return ""
}

func (state ChoiceState) GetNexts() []string {
	nexts := make([]string, len(state.Choices)+1)
	for i, choice := range state.Choices {
		nexts[i] = choice.Next
	}
	nexts[len(nexts)-1] = state.Default
	return nexts
}

type Choice struct {
	Condition Condition
	Next      string
}

type Condition interface {
	Eval(ctx context.Context, input interface{}) (bool, error)
}

type BooleanEqualsRule struct {
	V1 Path
	V2 bool
}

func (r BooleanEqualsRule) Eval(ctx context.Context, input interface{}) (bool, error) {
	v, err := UnjoinByPath(ctx, input, &r.V1)
	if err != nil {
		return false, err
	}

	v1, ok := v.(bool)
	if !ok {
		return false, fmt.Errorf("invalid field value (must be boolean) : [%s]=[%v]", r.V1, v)
	}

	return v1 == r.V2, nil
}
