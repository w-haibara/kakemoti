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
	/*
	 * String
	 */
	case isExistKey(m, "StringEquals"):
		v2, ok := m["StringEquals"].(string)
		if !ok {
			return nil, ErrInvalidType
		}
		return StringEqualsRule{v1, v2}, nil
	case isExistKey(m, "StringEqualsPath"):
		v, ok := m["StringEqualsPath"].(string)
		if !ok {
			return nil, ErrInvalidType
		}
		v2, err := NewPath(v)
		if err != nil {
			return nil, err
		}
		return StringEqualsPathRule{v1, v2}, nil
	case isExistKey(m, "StringGreaterThan"):
		v2, ok := m["StringGreaterThan"].(string)
		if !ok {
			return nil, ErrInvalidType
		}
		return StringGreaterThanRule{v1, v2}, nil
	case isExistKey(m, "StringGreaterThanPath"):
		v, ok := m["StringGreaterThanPath"].(string)
		if !ok {
			return nil, ErrInvalidType
		}
		v2, err := NewPath(v)
		if err != nil {
			return nil, err
		}
		return StringGreaterThanPathRule{v1, v2}, nil
	case isExistKey(m, "StringGreaterThanEquals"):
		v2, ok := m["StringGreaterEquals"].(string)
		if !ok {
			return nil, ErrInvalidType
		}
		return StringGreaterThanEqualsRule{v1, v2}, nil
	case isExistKey(m, "StringGreaterThanEqualsPath"):
		v, ok := m["StringGreaterThanEqualsPath"].(string)
		if !ok {
			return nil, ErrInvalidType
		}
		v2, err := NewPath(v)
		if err != nil {
			return nil, err
		}
		return StringGreaterThanEqualsPathRule{v1, v2}, nil
	case isExistKey(m, "StringLessThan"):
		v2, ok := m["StringLessThan"].(string)
		if !ok {
			return nil, ErrInvalidType
		}
		return StringLessThanRule{v1, v2}, nil
	case isExistKey(m, "StringLessThanPath"):
		v, ok := m["StringLessThanPath"].(string)
		if !ok {
			return nil, ErrInvalidType
		}
		v2, err := NewPath(v)
		if err != nil {
			return nil, err
		}
		return StringLessThanPathRule{v1, v2}, nil
	case isExistKey(m, "StringLessThanEquals"):
		v2, ok := m["StringLessThanEquals"].(string)
		if !ok {
			return nil, ErrInvalidType
		}
		return StringLessThanEqualsRule{v1, v2}, nil
	case isExistKey(m, "StringLessThanEqualsPath"):
		v, ok := m["StringLessThanEqualsPath"].(string)
		if !ok {
			return nil, ErrInvalidType
		}
		v2, err := NewPath(v)
		if err != nil {
			return nil, err
		}
		return StringLessThanEqualsPathRule{v1, v2}, nil
	case isExistKey(m, "StringMatches"):
		panic("Not Implemented")
	/*
	 * Numeric
	 */
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
	/*
	 * Boolean
	 */
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
	/*
	 * Timestamp
	 */
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
	/*
	 * Check type of value
	 */
	case isExistKey(m, "IsBoolean"):
		return IsBooleanRule{v1}, nil
	case isExistKey(m, "IsNull"):
		return IsNullRule{v1}, nil
	case isExistKey(m, "IsNumeric"):
		return IsNumericRule{v1}, nil
	case isExistKey(m, "IsPresent"):
		return IsPresentRule{v1}, nil
	case isExistKey(m, "IsString"):
		return IsStringRule{v1}, nil
	case isExistKey(m, "IsTimestamp"):
		return IsTimestampRule{v1}, nil
	/*
	 * Unknown Operator
	 */
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

type IsBooleanRule struct {
	V1 Path
}

func (r IsBooleanRule) Eval(ctx context.Context, input interface{}) (bool, error) {
	v, err := UnjoinByPath(ctx, input, &r.V1)
	if err != nil {
		return false, err
	}

	_, ok := v.(bool)
	return ok, nil
}

type IsNullRule struct {
	V1 Path
}

func (r IsNullRule) Eval(ctx context.Context, input interface{}) (bool, error) {
	v, err := UnjoinByPath(ctx, input, &r.V1)
	if err != nil {
		return false, err
	}

	return v == nil, nil
}

type IsNumericRule struct {
	V1 Path
}

func (r IsNumericRule) Eval(ctx context.Context, input interface{}) (bool, error) {
	v, err := UnjoinByPath(ctx, input, &r.V1)
	if err != nil {
		return false, err
	}

	switch v.(type) {
	case int, float64:
		return true, nil
	default:
		return false, nil
	}
}

type IsStringRule struct {
	V1 Path
}

func (r IsStringRule) Eval(ctx context.Context, input interface{}) (bool, error) {
	v, err := UnjoinByPath(ctx, input, &r.V1)
	if err != nil {
		return false, err
	}

	_, ok := v.(string)
	return ok, nil
}

type IsTimestampRule struct {
	V1 Path
}

func (r IsTimestampRule) Eval(ctx context.Context, input interface{}) (bool, error) {
	v, err := UnjoinByPath(ctx, input, &r.V1)
	if err != nil {
		return false, err
	}

	str, ok := v.(string)
	if !ok {
		return false, nil
	}

	_, err = NewTimestamp(str)
	if err != nil {
		return false, nil
	}

	return true, nil
}

type IsPresentRule struct {
	V1 Path
}

func (r IsPresentRule) Eval(ctx context.Context, input interface{}) (bool, error) {
	_, err := UnjoinByPath(ctx, input, &r.V1)
	if err != nil {
		return false, nil
	}

	return true, nil
}

type StringEqualsRule struct {
	V1 Path
	V2 string
}

func (r StringEqualsRule) Eval(ctx context.Context, input interface{}) (bool, error) {
	v1, err := GetString(ctx, input, r.V1)
	if err != nil {
		return false, err
	}

	return v1 == r.V2, nil
}

type StringEqualsPathRule struct {
	V1 Path
	V2 Path
}

func (r StringEqualsPathRule) Eval(ctx context.Context, input interface{}) (bool, error) {
	v1, err := GetString(ctx, input, r.V1)
	if err != nil {
		return false, err
	}

	v2, err := GetString(ctx, input, r.V2)
	if err != nil {
		return false, err
	}

	return v1 == v2, nil
}

type StringLessThanRule struct {
	V1 Path
	V2 string
}

func (r StringLessThanRule) Eval(ctx context.Context, input interface{}) (bool, error) {
	v1, err := GetString(ctx, input, r.V1)
	if err != nil {
		return false, err
	}

	return v1 < r.V2, nil
}

type StringLessThanPathRule struct {
	V1 Path
	V2 Path
}

func (r StringLessThanPathRule) Eval(ctx context.Context, input interface{}) (bool, error) {
	v1, err := GetString(ctx, input, r.V1)
	if err != nil {
		return false, err
	}

	v2, err := GetString(ctx, input, r.V2)
	if err != nil {
		return false, err
	}

	return v1 < v2, nil
}

type StringLessThanEqualsRule struct {
	V1 Path
	V2 string
}

func (r StringLessThanEqualsRule) Eval(ctx context.Context, input interface{}) (bool, error) {
	v1, err := GetString(ctx, input, r.V1)
	if err != nil {
		return false, err
	}

	return v1 <= r.V2, nil
}

type StringLessThanEqualsPathRule struct {
	V1 Path
	V2 Path
}

func (r StringLessThanEqualsPathRule) Eval(ctx context.Context, input interface{}) (bool, error) {
	v1, err := GetString(ctx, input, r.V1)
	if err != nil {
		return false, err
	}

	v2, err := GetString(ctx, input, r.V2)
	if err != nil {
		return false, err
	}

	return v1 <= v2, nil
}

type StringGreaterThanRule struct {
	V1 Path
	V2 string
}

func (r StringGreaterThanRule) Eval(ctx context.Context, input interface{}) (bool, error) {
	v1, err := GetString(ctx, input, r.V1)
	if err != nil {
		return false, err
	}

	return v1 > r.V2, nil
}

type StringGreaterThanPathRule struct {
	V1 Path
	V2 Path
}

func (r StringGreaterThanPathRule) Eval(ctx context.Context, input interface{}) (bool, error) {
	v1, err := GetString(ctx, input, r.V1)
	if err != nil {
		return false, err
	}

	v2, err := GetString(ctx, input, r.V2)
	if err != nil {
		return false, err
	}

	return v1 > v2, nil
}

type StringGreaterThanEqualsRule struct {
	V1 Path
	V2 string
}

func (r StringGreaterThanEqualsRule) Eval(ctx context.Context, input interface{}) (bool, error) {
	v1, err := GetString(ctx, input, r.V1)
	if err != nil {
		return false, err
	}

	return v1 >= r.V2, nil
}

type StringGreaterThanEqualsPathRule struct {
	V1 Path
	V2 Path
}

func (r StringGreaterThanEqualsPathRule) Eval(ctx context.Context, input interface{}) (bool, error) {
	v1, err := GetString(ctx, input, r.V1)
	if err != nil {
		return false, err
	}

	v2, err := GetString(ctx, input, r.V2)
	if err != nil {
		return false, err
	}

	return v1 >= v2, nil
}
