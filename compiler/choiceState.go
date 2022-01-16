package compiler

import (
	"context"
	"errors"
)

type RawChoiceState struct {
	CommonState2
	Choices []map[string]interface{} `json:"Choices"`
	Default string                   `json:"Default"`
}

var (
	ErrNotFound    = errors.New("key not found")
	ErrInvalidType = errors.New("invalid type")
)

func (raw RawChoiceState) decode() (*ChoiceState, error) {
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

		cond, err := decodeBoolExpr(raw)
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

func isExistKey(m map[string]interface{}, k string) bool {
	_, ok := m[k]
	return ok
}

func decodeDataTestExpr(m map[string]interface{}) (Condition, error) {
	if !isExistKey(m, "Variable") {
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
	case isExistKey(m, "BooleanEquals"):
		v2, ok := m["BooleanEquals"].(bool)
		if !ok {
			return nil, ErrInvalidType
		}
		return BooleanEqualsRule{v1, v2}, nil
	case isExistKey(m, "BooleanEqualsPath"):
		v, ok := m["BooleanEqualsPath"].(string)
		if !ok {
			return nil, ErrInvalidType
		}
		v2, err := NewPath(v)
		if err != nil {
			return nil, err
		}
		return BooleanEqualsPathRule{v1, v2}, nil
	case isExistKey(m, "IsBoolean"):
		panic("Not Implemented")
	case isExistKey(m, "IsNull"):
		panic("Not Implemented")
	case isExistKey(m, "IsNumeric"):
		panic("Not Implemented")
	case isExistKey(m, "IsPresent"):
		panic("Not Implemented")
	case isExistKey(m, "IsString"):
		panic("Not Implemented")
	case isExistKey(m, "IsTimestamp"):
		panic("Not Implemented")
	case isExistKey(m, "NumericEquals"):
		panic("Not Implemented")
	case isExistKey(m, "NumericEqualsPath"):
		panic("Not Implemented")
	case isExistKey(m, "NumericGreaterThan"):
		panic("Not Implemented")
	case isExistKey(m, "NumericGreaterThanPath"):
		panic("Not Implemented")
	case isExistKey(m, "NumericGreaterThanEquals"):
		panic("Not Implemented")
	case isExistKey(m, "NumericGreaterThanEqualsPath"):
		panic("Not Implemented")
	case isExistKey(m, "NumericLessThan"):
		panic("Not Implemented")
	case isExistKey(m, "NumericLessThanPath"):
		panic("Not Implemented")
	case isExistKey(m, "NumericLessThanEquals"):
		panic("Not Implemented")
	case isExistKey(m, "NumericLessThanEqualsPath"):
		panic("Not Implemented")
	case isExistKey(m, "StringEquals"):
		panic("Not Implemented")
	case isExistKey(m, "StringEqualsPath"):
		panic("Not Implemented")
	case isExistKey(m, "StringGreaterThan"):
		panic("Not Implemented")
	case isExistKey(m, "StringGreaterThanPath"):
		panic("Not Implemented")
	case isExistKey(m, "StringGreaterThanEquals"):
		panic("Not Implemented")
	case isExistKey(m, "StringGreaterThanEqualsPath"):
		panic("Not Implemented")
	case isExistKey(m, "StringLessThanStringLessThanPath"):
		panic("Not Implemented")
	case isExistKey(m, "StringLessThanEqualsStringLessThanEqualsPath"):
		panic("Not Implemented")
	case isExistKey(m, "StringMatches"):
		panic("Not Implemented")
	case isExistKey(m, "TimestampEquals"):
		panic("Not Implemented")
	case isExistKey(m, "TimestampEqualsPath"):
		panic("Not Implemented")
	case isExistKey(m, "TimestampGreaterThan"):
		panic("Not Implemented")
	case isExistKey(m, "TimestampGreaterThanPath"):
		panic("Not Implemented")
	case isExistKey(m, "TimestampGreaterThanEquals"):
		panic("Not Implemented")
	case isExistKey(m, "TimestampGreaterThanEqualsPath"):
		panic("Not Implemented")
	case isExistKey(m, "TimestampLessThan"):
		panic("Not Implemented")
	case isExistKey(m, "TimestampLessThanPath"):
		panic("Not Implemented")
	case isExistKey(m, "TimestampLessThanEquals"):
		panic("Not Implemented")
	case isExistKey(m, "TimestampLessThanEqualsPath"):
		panic("Not Implemented")
	default:
		panic("Unknown Operator")
	}
}

func decodeConds(m map[string]interface{}, key string) ([]Condition, error) {
	maps, ok := m[key].([]interface{})
	if !ok {
		return nil, ErrInvalidType
	}

	conds := make([]Condition, len(maps))
	for i, m := range maps {
		v, ok := m.(map[string]interface{})
		if !ok {
			return nil, ErrInvalidType
		}

		c, err := decodeBoolExpr(v)
		if err != nil {
			return nil, err
		}
		conds[i] = c
	}

	return conds, nil
}

func decodeBoolExpr(m map[string]interface{}) (Condition, error) {
	switch {
	case isExistKey(m, "And"):
		conds, err := decodeConds(m, "And")
		if err != nil {
			return nil, err
		}
		return AndRule{conds}, nil
	case isExistKey(m, "Or"):
		conds, err := decodeConds(m, "Or")
		if err != nil {
			return nil, err
		}
		return OrRule{conds}, nil
	case isExistKey(m, "Not"):
		v, ok := m["Not"]
		if !ok {
			return nil, ErrNotFound
		}
		v1, ok := v.(map[string]interface{})
		if !ok {
			return nil, ErrInvalidType
		}

		c, err := decodeBoolExpr(v1)
		if err != nil {
			return nil, err
		}
		return NotRule{c}, nil
	default:
		return decodeDataTestExpr(m)
	}
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

type AndRule struct {
	V []Condition
}

func (r AndRule) Eval(ctx context.Context, input interface{}) (bool, error) {
	res := true
	for _, v := range r.V {
		b, err := v.Eval(ctx, input)
		if err != nil {
			return false, err
		}
		res = res && b
	}

	return res, nil
}

type OrRule struct {
	V []Condition
}

func (r OrRule) Eval(ctx context.Context, input interface{}) (bool, error) {
	res := false
	for _, v := range r.V {
		b, err := v.Eval(ctx, input)
		if err != nil {
			return false, err
		}
		res = res || b
	}

	return res, nil
}

type NotRule struct {
	V Condition
}

func (r NotRule) Eval(ctx context.Context, input interface{}) (bool, error) {
	b, err := r.V.Eval(ctx, input)
	if err != nil {
		return false, err
	}

	return !b, nil
}

type BooleanEqualsRule struct {
	V1 Path
	V2 bool
}

func (r BooleanEqualsRule) Eval(ctx context.Context, input interface{}) (bool, error) {
	v1, err := GetBool(ctx, input, r.V1)
	if err != nil {
		return false, err
	}

	return v1 == r.V2, nil
}

type BooleanEqualsPathRule struct {
	V1 Path
	V2 Path
}

func (r BooleanEqualsPathRule) Eval(ctx context.Context, input interface{}) (bool, error) {
	v1, err := GetBool(ctx, input, r.V1)
	if err != nil {
		return false, err
	}

	v2, err := GetBool(ctx, input, r.V2)
	if err != nil {
		return false, err
	}

	return v1 == v2, nil
}
