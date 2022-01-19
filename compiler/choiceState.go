package compiler

import (
	"context"
	"errors"
	"fmt"
	"runtime"
)

type RawChoiceState struct {
	CommonState2
	Choices []map[string]interface{} `json:"Choices"`
	Default string                   `json:"Default"`
}

var (
	ErrNotFound = errors.New("key not found")
)

func invalidTypeError() error {
	_, _, line, _ := runtime.Caller(1)
	return fmt.Errorf("invalid type, line: %v", line)
}

func (raw RawChoiceState) decode() (*ChoiceState, error) {
	choices := make([]Choice, len(raw.Choices))
	for i, raw := range raw.Choices {
		n, ok := raw["Next"]
		if !ok {
			return nil, ErrNotFound
		}
		next, ok := n.(string)
		if !ok {
			return nil, invalidTypeError()
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
		return nil, invalidTypeError()
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
			return nil, invalidTypeError()
		}
		return StringEqualsRule{v1, v2}, nil
	case isExistKey(m, "StringEqualsPath"):
		v, ok := m["StringEqualsPath"].(string)
		if !ok {
			return nil, invalidTypeError()
		}
		v2, err := NewPath(v)
		if err != nil {
			return nil, err
		}
		return StringEqualsPathRule{v1, v2}, nil
	case isExistKey(m, "StringGreaterThan"):
		v2, ok := m["StringGreaterThan"].(string)
		if !ok {
			return nil, invalidTypeError()
		}
		return StringGreaterThanRule{v1, v2}, nil
	case isExistKey(m, "StringGreaterThanPath"):
		v, ok := m["StringGreaterThanPath"].(string)
		if !ok {
			return nil, invalidTypeError()
		}
		v2, err := NewPath(v)
		if err != nil {
			return nil, err
		}
		return StringGreaterThanPathRule{v1, v2}, nil
	case isExistKey(m, "StringGreaterThanEquals"):
		v2, ok := m["StringGreaterEquals"].(string)
		if !ok {
			return nil, invalidTypeError()
		}
		return StringGreaterThanEqualsRule{v1, v2}, nil
	case isExistKey(m, "StringGreaterThanEqualsPath"):
		v, ok := m["StringGreaterThanEqualsPath"].(string)
		if !ok {
			return nil, invalidTypeError()
		}
		v2, err := NewPath(v)
		if err != nil {
			return nil, err
		}
		return StringGreaterThanEqualsPathRule{v1, v2}, nil
	case isExistKey(m, "StringLessThan"):
		v2, ok := m["StringLessThan"].(string)
		if !ok {
			return nil, invalidTypeError()
		}
		return StringLessThanRule{v1, v2}, nil
	case isExistKey(m, "StringLessThanPath"):
		v, ok := m["StringLessThanPath"].(string)
		if !ok {
			return nil, invalidTypeError()
		}
		v2, err := NewPath(v)
		if err != nil {
			return nil, err
		}
		return StringLessThanPathRule{v1, v2}, nil
	case isExistKey(m, "StringLessThanEquals"):
		v2, ok := m["StringLessThanEquals"].(string)
		if !ok {
			return nil, invalidTypeError()
		}
		return StringLessThanEqualsRule{v1, v2}, nil
	case isExistKey(m, "StringLessThanEqualsPath"):
		v, ok := m["StringLessThanEqualsPath"].(string)
		if !ok {
			return nil, invalidTypeError()
		}
		v2, err := NewPath(v)
		if err != nil {
			return nil, err
		}
		return StringLessThanEqualsPathRule{v1, v2}, nil
	case isExistKey(m, "StringMatches"):
		v2, ok := m["StringMatches"].(string)
		if !ok {
			return nil, invalidTypeError()
		}
		return StringMatchesRule{v1, v2}, nil
	/*
	 * Numeric
	 */
	case isExistKey(m, "NumericEquals"):
		v2, ok := m["NumericEquals"].(float64)
		if !ok {
			return nil, invalidTypeError()
		}
		return NumericEqualsRule{v1, v2}, nil
	case isExistKey(m, "NumericEqualsPath"):
		v, ok := m["NumericEqualsPath"].(string)
		if !ok {
			return nil, invalidTypeError()
		}
		v2, err := NewPath(v)
		if err != nil {
			return nil, err
		}
		return NumericEqualsPathRule{v1, v2}, nil
	case isExistKey(m, "NumericGreaterThan"):
		v2, ok := m["NumericGreaterThan"].(float64)
		if !ok {
			return nil, invalidTypeError()
		}
		return NumericGreaterThanRule{v1, v2}, nil
	case isExistKey(m, "NumericGreaterThanPath"):
		v, ok := m["NumericGreaterThanPath"].(string)
		if !ok {
			return nil, invalidTypeError()
		}
		v2, err := NewPath(v)
		if err != nil {
			return nil, err
		}
		return NumericGreaterThanPathRule{v1, v2}, nil
	case isExistKey(m, "NumericGreaterThanEquals"):
		v2, ok := m["NumericGreaterThanEquals"].(float64)
		if !ok {
			return nil, invalidTypeError()
		}
		return NumericGreaterThanEqualsRule{v1, v2}, nil
	case isExistKey(m, "NumericGreaterThanEqualsPath"):
		v, ok := m["NumericGreaterThanEqualsPathRule"].(string)
		if !ok {
			return nil, invalidTypeError()
		}
		v2, err := NewPath(v)
		if err != nil {
			return nil, err
		}
		return NumericGreaterThanEqualsPathRule{v1, v2}, nil
	case isExistKey(m, "NumericLessThan"):
		v2, ok := m["NumericLessThan"].(float64)
		if !ok {
			return nil, invalidTypeError()
		}
		return NumericLessThanRule{v1, v2}, nil
	case isExistKey(m, "NumericLessThanPath"):
		v, ok := m["NumericLessThanPath"].(string)
		if !ok {
			return nil, invalidTypeError()
		}
		v2, err := NewPath(v)
		if err != nil {
			return nil, err
		}
		return NumericLessThanPathRule{v1, v2}, nil
	case isExistKey(m, "NumericLessThanEquals"):
		v2, ok := m["NumericLessThanEquals"].(float64)
		if !ok {
			return nil, invalidTypeError()
		}
		return NumericLessThanEqualsRule{v1, v2}, nil
	case isExistKey(m, "NumericLessThanEqualsPath"):
		v, ok := m["NumericLessThanEqualsPath"].(string)
		if !ok {
			return nil, invalidTypeError()
		}
		v2, err := NewPath(v)
		if err != nil {
			return nil, err
		}
		return NumericLessThanEqualsPathRule{v1, v2}, nil
	/*
	 * Boolean
	 */
	case isExistKey(m, "BooleanEquals"):
		v2, ok := m["BooleanEquals"].(bool)
		if !ok {
			return nil, invalidTypeError()
		}
		return BooleanEqualsRule{v1, v2}, nil
	case isExistKey(m, "BooleanEqualsPath"):
		v, ok := m["BooleanEqualsPath"].(string)
		if !ok {
			return nil, invalidTypeError()
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
		v, ok := m["TimestampEquals"].(string)
		if !ok {
			return nil, invalidTypeError()
		}
		v2, err := NewTimestamp(v)
		if err != nil {
			return nil, err
		}
		return TimestampEqualsRule{v1, v2}, nil
	case isExistKey(m, "TimestampEqualsPath"):
		v, ok := m["TimestampEqualsPath"].(string)
		if !ok {
			return nil, invalidTypeError()
		}
		v2, err := NewPath(v)
		if err != nil {
			return nil, err
		}
		return TimestampEqualsPathRule{v1, v2}, nil
	case isExistKey(m, "TimestampGreaterThan"):
		v, ok := m["TimestampGreaterThan"].(string)
		if !ok {
			return nil, invalidTypeError()
		}
		v2, err := NewTimestamp(v)
		if err != nil {
			return nil, err
		}
		return TimestampGreaterThanRule{v1, v2}, nil
	case isExistKey(m, "TimestampGreaterThanPath"):
		v, ok := m["TimestampGreaterThanPath"].(string)
		if !ok {
			return nil, invalidTypeError()
		}
		v2, err := NewPath(v)
		if err != nil {
			return nil, err
		}
		return TimestampGreaterThanPathRule{v1, v2}, nil
	case isExistKey(m, "TimestampGreaterThanEquals"):
		v, ok := m["TimestampGreaterThanEquals"].(string)
		if !ok {
			return nil, invalidTypeError()
		}
		v2, err := NewTimestamp(v)
		if err != nil {
			return nil, err
		}
		return TimestampGreaterThanEqualsRule{v1, v2}, nil
	case isExistKey(m, "TimestampGreaterThanEqualsPath"):
		v, ok := m["TimestampGreaterThanEqualsPath"].(string)
		if !ok {
			return nil, invalidTypeError()
		}
		v2, err := NewPath(v)
		if err != nil {
			return nil, err
		}
		return TimestampGreaterThanEqualsPathRule{v1, v2}, nil
	case isExistKey(m, "TimestampLessThan"):
		v, ok := m["TimestampLessThan"].(string)
		if !ok {
			return nil, invalidTypeError()
		}
		v2, err := NewTimestamp(v)
		if err != nil {
			return nil, err
		}
		return TimestampLessThanRule{v1, v2}, nil
	case isExistKey(m, "TimestampLessThanPath"):
		v, ok := m["TimestampLessThanPath"].(string)
		if !ok {
			return nil, invalidTypeError()
		}
		v2, err := NewPath(v)
		if err != nil {
			return nil, err
		}
		return TimestampLessThanPathRule{v1, v2}, nil
	case isExistKey(m, "TimestampLessThanEquals"):
		v, ok := m["TimestampLessThanEquals"].(string)
		if !ok {
			return nil, invalidTypeError()
		}
		v2, err := NewTimestamp(v)
		if err != nil {
			return nil, err
		}
		return TimestampLessThanEqualsRule{v1, v2}, nil
	case isExistKey(m, "TimestampLessThanEqualsPath"):
		v, ok := m["TimestampLessThanEqualsPath"].(string)
		if !ok {
			return nil, invalidTypeError()
		}
		v2, err := NewPath(v)
		if err != nil {
			return nil, err
		}
		return TimestampLessThanEqualsPathRule{v1, v2}, nil

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
		return nil, invalidTypeError()
	}

	conds := make([]Condition, len(maps))
	for i, m := range maps {
		v, ok := m.(map[string]interface{})
		if !ok {
			return nil, invalidTypeError()
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
			return nil, invalidTypeError()
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
	for _, v := range r.V {
		b, err := v.Eval(ctx, input)
		if err != nil {
			return false, err
		}
		if !b {
			return false, nil
		}
	}

	return true, nil
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

type StringMatchesRule struct {
	V1 Path
	V2 string
}

var ErrOpenBackslashFound = errors.New("open backslash found")

func (r StringMatchesRule) Eval(ctx context.Context, input interface{}) (bool, error) {
	v1, err := GetString(ctx, input, r.V1)
	if err != nil {
		return false, err
	}

	pos := 0
	for i := 0; i < len(r.V2); i++ {
		switch r.V2[i] {
		case '\\':
			if i == len(r.V2)-1 {
				return false, ErrOpenBackslashFound
			}
			switch r.V2[i+1] {
			case '*':
				i++
				if '*' != v1[pos] {
					return false, nil
				}
			case '\\':
				i++
				if '\\' != v1[pos] {
					return false, nil
				}
			default:
				return false, ErrOpenBackslashFound
			}
		case '*':
			if i == len(r.V2)-1 {
				return true, nil
			}
			if r.V2[i+1] == v1[pos] {
				pos--
				break
			}
			for {
				pos++
				if pos >= len(v1) {
					return false, nil
				}
				if r.V2[i+1] == v1[pos] {
					i++
					break
				}
			}
		default:
			if r.V2[i] != v1[pos] {
				return false, nil
			}
		}
		pos++
	}

	return true, nil
}

type NumericEqualsRule struct {
	V1 Path
	V2 float64
}

func (r NumericEqualsRule) Eval(ctx context.Context, input interface{}) (bool, error) {
	v1, err := GetNumeric(ctx, input, r.V1)
	if err != nil {
		return false, err
	}

	return v1 == r.V2, nil
}

type NumericEqualsPathRule struct {
	V1 Path
	V2 Path
}

func (r NumericEqualsPathRule) Eval(ctx context.Context, input interface{}) (bool, error) {
	v1, err := GetNumeric(ctx, input, r.V1)
	if err != nil {
		return false, err
	}

	v2, err := GetNumeric(ctx, input, r.V2)
	if err != nil {
		return false, err
	}

	return v1 == v2, nil
}

type NumericLessThanRule struct {
	V1 Path
	V2 float64
}

func (r NumericLessThanRule) Eval(ctx context.Context, input interface{}) (bool, error) {
	v1, err := GetNumeric(ctx, input, r.V1)
	if err != nil {
		return false, err
	}

	return v1 < r.V2, nil
}

type NumericLessThanPathRule struct {
	V1 Path
	V2 Path
}

func (r NumericLessThanPathRule) Eval(ctx context.Context, input interface{}) (bool, error) {
	v1, err := GetNumeric(ctx, input, r.V1)
	if err != nil {
		return false, err
	}

	v2, err := GetNumeric(ctx, input, r.V2)
	if err != nil {
		return false, err
	}

	return v1 < v2, nil
}

type NumericLessThanEqualsRule struct {
	V1 Path
	V2 float64
}

func (r NumericLessThanEqualsRule) Eval(ctx context.Context, input interface{}) (bool, error) {
	v1, err := GetNumeric(ctx, input, r.V1)
	if err != nil {
		return false, err
	}

	return v1 <= r.V2, nil
}

type NumericLessThanEqualsPathRule struct {
	V1 Path
	V2 Path
}

func (r NumericLessThanEqualsPathRule) Eval(ctx context.Context, input interface{}) (bool, error) {
	v1, err := GetNumeric(ctx, input, r.V1)
	if err != nil {
		return false, err
	}

	v2, err := GetNumeric(ctx, input, r.V2)
	if err != nil {
		return false, err
	}

	return v1 <= v2, nil
}

type NumericGreaterThanRule struct {
	V1 Path
	V2 float64
}

func (r NumericGreaterThanRule) Eval(ctx context.Context, input interface{}) (bool, error) {
	v1, err := GetNumeric(ctx, input, r.V1)
	if err != nil {
		return false, err
	}

	return v1 > r.V2, nil
}

type NumericGreaterThanPathRule struct {
	V1 Path
	V2 Path
}

func (r NumericGreaterThanPathRule) Eval(ctx context.Context, input interface{}) (bool, error) {
	v1, err := GetNumeric(ctx, input, r.V1)
	if err != nil {
		return false, err
	}

	v2, err := GetNumeric(ctx, input, r.V2)
	if err != nil {
		return false, err
	}

	return v1 > v2, nil
}

type NumericGreaterThanEqualsRule struct {
	V1 Path
	V2 float64
}

func (r NumericGreaterThanEqualsRule) Eval(ctx context.Context, input interface{}) (bool, error) {
	v1, err := GetNumeric(ctx, input, r.V1)
	if err != nil {
		return false, err
	}

	return v1 >= r.V2, nil
}

type NumericGreaterThanEqualsPathRule struct {
	V1 Path
	V2 Path
}

func (r NumericGreaterThanEqualsPathRule) Eval(ctx context.Context, input interface{}) (bool, error) {
	v1, err := GetNumeric(ctx, input, r.V1)
	if err != nil {
		return false, err
	}

	v2, err := GetNumeric(ctx, input, r.V2)
	if err != nil {
		return false, err
	}

	return v1 >= v2, nil
}

type TimestampEqualsRule struct {
	V1 Path
	V2 Timestamp
}

func (r TimestampEqualsRule) Eval(ctx context.Context, input interface{}) (bool, error) {
	v1, err := GetTimestamp(ctx, input, r.V1)
	if err != nil {
		return false, err
	}

	return v1.Equals(r.V2), nil
}

type TimestampEqualsPathRule struct {
	V1 Path
	V2 Path
}

func (r TimestampEqualsPathRule) Eval(ctx context.Context, input interface{}) (bool, error) {
	v1, err := GetTimestamp(ctx, input, r.V1)
	if err != nil {
		return false, err
	}

	v2, err := GetTimestamp(ctx, input, r.V2)
	if err != nil {
		return false, err
	}

	return v1.Equals(v2), nil
}

type TimestampLessThanRule struct {
	V1 Path
	V2 Timestamp
}

func (r TimestampLessThanRule) Eval(ctx context.Context, input interface{}) (bool, error) {
	v1, err := GetTimestamp(ctx, input, r.V1)
	if err != nil {
		return false, err
	}

	return v1.LessThan(r.V2), nil
}

type TimestampLessThanPathRule struct {
	V1 Path
	V2 Path
}

func (r TimestampLessThanPathRule) Eval(ctx context.Context, input interface{}) (bool, error) {
	v1, err := GetTimestamp(ctx, input, r.V1)
	if err != nil {
		return false, err
	}

	v2, err := GetTimestamp(ctx, input, r.V2)
	if err != nil {
		return false, err
	}

	return v1.LessThan(v2), nil
}

type TimestampLessThanEqualsRule struct {
	V1 Path
	V2 Timestamp
}

func (r TimestampLessThanEqualsRule) Eval(ctx context.Context, input interface{}) (bool, error) {
	v1, err := GetTimestamp(ctx, input, r.V1)
	if err != nil {
		return false, err
	}

	return v1.LessThanEquals(r.V2), nil
}

type TimestampLessThanEqualsPathRule struct {
	V1 Path
	V2 Path
}

func (r TimestampLessThanEqualsPathRule) Eval(ctx context.Context, input interface{}) (bool, error) {
	v1, err := GetTimestamp(ctx, input, r.V1)
	if err != nil {
		return false, err
	}

	v2, err := GetTimestamp(ctx, input, r.V2)
	if err != nil {
		return false, err
	}

	return v1.LessThanEquals(v2), nil
}

type TimestampGreaterThanRule struct {
	V1 Path
	V2 Timestamp
}

func (r TimestampGreaterThanRule) Eval(ctx context.Context, input interface{}) (bool, error) {
	v1, err := GetTimestamp(ctx, input, r.V1)
	if err != nil {
		return false, err
	}

	return v1.GreaterThan(r.V2), nil
}

type TimestampGreaterThanPathRule struct {
	V1 Path
	V2 Path
}

func (r TimestampGreaterThanPathRule) Eval(ctx context.Context, input interface{}) (bool, error) {
	v1, err := GetTimestamp(ctx, input, r.V1)
	if err != nil {
		return false, err
	}

	v2, err := GetTimestamp(ctx, input, r.V2)
	if err != nil {
		return false, err
	}

	return v1.GreaterThan(v2), nil
}

type TimestampGreaterThanEqualsRule struct {
	V1 Path
	V2 Timestamp
}

func (r TimestampGreaterThanEqualsRule) Eval(ctx context.Context, input interface{}) (bool, error) {
	v1, err := GetTimestamp(ctx, input, r.V1)
	if err != nil {
		return false, err
	}

	return v1.GreaterThanEquals(r.V2), nil
}

type TimestampGreaterThanEqualsPathRule struct {
	V1 Path
	V2 Path
}

func (r TimestampGreaterThanEqualsPathRule) Eval(ctx context.Context, input interface{}) (bool, error) {
	v1, err := GetTimestamp(ctx, input, r.V1)
	if err != nil {
		return false, err
	}

	v2, err := GetTimestamp(ctx, input, r.V2)
	if err != nil {
		return false, err
	}

	return v1.GreaterThanEquals(v2), nil
}
