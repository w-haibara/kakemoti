package compiler

import (
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
		v2, ok := m["StringGreaterThanEquals"].(string)
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
		v, ok := m["NumericGreaterThanEqualsPath"].(string)
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
