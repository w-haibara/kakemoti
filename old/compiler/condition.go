package compiler

import (
	"errors"
)

type Condition interface {
	Eval(coj *CtxObj, input interface{}) (bool, error)
}

type AndRule struct {
	V []Condition
}

func (r AndRule) Eval(coj *CtxObj, input interface{}) (bool, error) {
	for _, v := range r.V {
		b, err := v.Eval(coj, input)
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

func (r OrRule) Eval(coj *CtxObj, input interface{}) (bool, error) {
	res := false
	for _, v := range r.V {
		b, err := v.Eval(coj, input)
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

func (r NotRule) Eval(coj *CtxObj, input interface{}) (bool, error) {
	b, err := r.V.Eval(coj, input)
	if err != nil {
		return false, err
	}

	return !b, nil
}

type StringEqualsRule struct {
	V1 Path
	V2 string
}

func (r StringEqualsRule) Eval(coj *CtxObj, input interface{}) (bool, error) {
	v1, err := GetString(coj, input, r.V1)
	if err != nil {
		return false, err
	}

	return v1 == r.V2, nil
}

type StringEqualsPathRule struct {
	V1 Path
	V2 Path
}

func (r StringEqualsPathRule) Eval(coj *CtxObj, input interface{}) (bool, error) {
	v1, err := GetString(coj, input, r.V1)
	if err != nil {
		return false, err
	}

	v2, err := GetString(coj, input, r.V2)
	if err != nil {
		return false, err
	}

	return v1 == v2, nil
}

type StringLessThanRule struct {
	V1 Path
	V2 string
}

func (r StringLessThanRule) Eval(coj *CtxObj, input interface{}) (bool, error) {
	v1, err := GetString(coj, input, r.V1)
	if err != nil {
		return false, err
	}

	return v1 < r.V2, nil
}

type StringLessThanPathRule struct {
	V1 Path
	V2 Path
}

func (r StringLessThanPathRule) Eval(coj *CtxObj, input interface{}) (bool, error) {
	v1, err := GetString(coj, input, r.V1)
	if err != nil {
		return false, err
	}

	v2, err := GetString(coj, input, r.V2)
	if err != nil {
		return false, err
	}

	return v1 < v2, nil
}

type StringLessThanEqualsRule struct {
	V1 Path
	V2 string
}

func (r StringLessThanEqualsRule) Eval(coj *CtxObj, input interface{}) (bool, error) {
	v1, err := GetString(coj, input, r.V1)
	if err != nil {
		return false, err
	}

	return v1 <= r.V2, nil
}

type StringLessThanEqualsPathRule struct {
	V1 Path
	V2 Path
}

func (r StringLessThanEqualsPathRule) Eval(coj *CtxObj, input interface{}) (bool, error) {
	v1, err := GetString(coj, input, r.V1)
	if err != nil {
		return false, err
	}

	v2, err := GetString(coj, input, r.V2)
	if err != nil {
		return false, err
	}

	return v1 <= v2, nil
}

type StringGreaterThanRule struct {
	V1 Path
	V2 string
}

func (r StringGreaterThanRule) Eval(coj *CtxObj, input interface{}) (bool, error) {
	v1, err := GetString(coj, input, r.V1)
	if err != nil {
		return false, err
	}

	return v1 > r.V2, nil
}

type StringGreaterThanPathRule struct {
	V1 Path
	V2 Path
}

func (r StringGreaterThanPathRule) Eval(coj *CtxObj, input interface{}) (bool, error) {
	v1, err := GetString(coj, input, r.V1)
	if err != nil {
		return false, err
	}

	v2, err := GetString(coj, input, r.V2)
	if err != nil {
		return false, err
	}

	return v1 > v2, nil
}

type StringGreaterThanEqualsRule struct {
	V1 Path
	V2 string
}

func (r StringGreaterThanEqualsRule) Eval(coj *CtxObj, input interface{}) (bool, error) {
	v1, err := GetString(coj, input, r.V1)
	if err != nil {
		return false, err
	}

	return v1 >= r.V2, nil
}

type StringGreaterThanEqualsPathRule struct {
	V1 Path
	V2 Path
}

func (r StringGreaterThanEqualsPathRule) Eval(coj *CtxObj, input interface{}) (bool, error) {
	v1, err := GetString(coj, input, r.V1)
	if err != nil {
		return false, err
	}

	v2, err := GetString(coj, input, r.V2)
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

func (r StringMatchesRule) Eval(coj *CtxObj, input interface{}) (bool, error) {
	v1, err := GetString(coj, input, r.V1)
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

func (r NumericEqualsRule) Eval(coj *CtxObj, input interface{}) (bool, error) {
	v1, err := GetNumeric(coj, input, r.V1)
	if err != nil {
		return false, err
	}

	return v1 == r.V2, nil
}

type NumericEqualsPathRule struct {
	V1 Path
	V2 Path
}

func (r NumericEqualsPathRule) Eval(coj *CtxObj, input interface{}) (bool, error) {
	v1, err := GetNumeric(coj, input, r.V1)
	if err != nil {
		return false, err
	}

	v2, err := GetNumeric(coj, input, r.V2)
	if err != nil {
		return false, err
	}

	return v1 == v2, nil
}

type NumericLessThanRule struct {
	V1 Path
	V2 float64
}

func (r NumericLessThanRule) Eval(coj *CtxObj, input interface{}) (bool, error) {
	v1, err := GetNumeric(coj, input, r.V1)
	if err != nil {
		return false, err
	}

	return v1 < r.V2, nil
}

type NumericLessThanPathRule struct {
	V1 Path
	V2 Path
}

func (r NumericLessThanPathRule) Eval(coj *CtxObj, input interface{}) (bool, error) {
	v1, err := GetNumeric(coj, input, r.V1)
	if err != nil {
		return false, err
	}

	v2, err := GetNumeric(coj, input, r.V2)
	if err != nil {
		return false, err
	}

	return v1 < v2, nil
}

type NumericLessThanEqualsRule struct {
	V1 Path
	V2 float64
}

func (r NumericLessThanEqualsRule) Eval(coj *CtxObj, input interface{}) (bool, error) {
	v1, err := GetNumeric(coj, input, r.V1)
	if err != nil {
		return false, err
	}

	return v1 <= r.V2, nil
}

type NumericLessThanEqualsPathRule struct {
	V1 Path
	V2 Path
}

func (r NumericLessThanEqualsPathRule) Eval(coj *CtxObj, input interface{}) (bool, error) {
	v1, err := GetNumeric(coj, input, r.V1)
	if err != nil {
		return false, err
	}

	v2, err := GetNumeric(coj, input, r.V2)
	if err != nil {
		return false, err
	}

	return v1 <= v2, nil
}

type NumericGreaterThanRule struct {
	V1 Path
	V2 float64
}

func (r NumericGreaterThanRule) Eval(coj *CtxObj, input interface{}) (bool, error) {
	v1, err := GetNumeric(coj, input, r.V1)
	if err != nil {
		return false, err
	}

	return v1 > r.V2, nil
}

type NumericGreaterThanPathRule struct {
	V1 Path
	V2 Path
}

func (r NumericGreaterThanPathRule) Eval(coj *CtxObj, input interface{}) (bool, error) {
	v1, err := GetNumeric(coj, input, r.V1)
	if err != nil {
		return false, err
	}

	v2, err := GetNumeric(coj, input, r.V2)
	if err != nil {
		return false, err
	}

	return v1 > v2, nil
}

type NumericGreaterThanEqualsRule struct {
	V1 Path
	V2 float64
}

func (r NumericGreaterThanEqualsRule) Eval(coj *CtxObj, input interface{}) (bool, error) {
	v1, err := GetNumeric(coj, input, r.V1)
	if err != nil {
		return false, err
	}

	return v1 >= r.V2, nil
}

type NumericGreaterThanEqualsPathRule struct {
	V1 Path
	V2 Path
}

func (r NumericGreaterThanEqualsPathRule) Eval(coj *CtxObj, input interface{}) (bool, error) {
	v1, err := GetNumeric(coj, input, r.V1)
	if err != nil {
		return false, err
	}

	v2, err := GetNumeric(coj, input, r.V2)
	if err != nil {
		return false, err
	}

	return v1 >= v2, nil
}

type BooleanEqualsRule struct {
	V1 Path
	V2 bool
}

func (r BooleanEqualsRule) Eval(coj *CtxObj, input interface{}) (bool, error) {
	v1, err := GetBool(coj, input, r.V1)
	if err != nil {
		return false, err
	}

	return v1 == r.V2, nil
}

type BooleanEqualsPathRule struct {
	V1 Path
	V2 Path
}

func (r BooleanEqualsPathRule) Eval(coj *CtxObj, input interface{}) (bool, error) {
	v1, err := GetBool(coj, input, r.V1)
	if err != nil {
		return false, err
	}

	v2, err := GetBool(coj, input, r.V2)
	if err != nil {
		return false, err
	}

	return v1 == v2, nil
}

type TimestampEqualsRule struct {
	V1 Path
	V2 Timestamp
}

func (r TimestampEqualsRule) Eval(coj *CtxObj, input interface{}) (bool, error) {
	v1, err := GetTimestamp(coj, input, r.V1)
	if err != nil {
		return false, err
	}

	return v1.Equals(r.V2), nil
}

type TimestampEqualsPathRule struct {
	V1 Path
	V2 Path
}

func (r TimestampEqualsPathRule) Eval(coj *CtxObj, input interface{}) (bool, error) {
	v1, err := GetTimestamp(coj, input, r.V1)
	if err != nil {
		return false, err
	}

	v2, err := GetTimestamp(coj, input, r.V2)
	if err != nil {
		return false, err
	}

	return v1.Equals(v2), nil
}

type TimestampLessThanRule struct {
	V1 Path
	V2 Timestamp
}

func (r TimestampLessThanRule) Eval(coj *CtxObj, input interface{}) (bool, error) {
	v1, err := GetTimestamp(coj, input, r.V1)
	if err != nil {
		return false, err
	}

	return v1.LessThan(r.V2), nil
}

type TimestampLessThanPathRule struct {
	V1 Path
	V2 Path
}

func (r TimestampLessThanPathRule) Eval(coj *CtxObj, input interface{}) (bool, error) {
	v1, err := GetTimestamp(coj, input, r.V1)
	if err != nil {
		return false, err
	}

	v2, err := GetTimestamp(coj, input, r.V2)
	if err != nil {
		return false, err
	}

	return v1.LessThan(v2), nil
}

type TimestampLessThanEqualsRule struct {
	V1 Path
	V2 Timestamp
}

func (r TimestampLessThanEqualsRule) Eval(coj *CtxObj, input interface{}) (bool, error) {
	v1, err := GetTimestamp(coj, input, r.V1)
	if err != nil {
		return false, err
	}

	return v1.LessThanEquals(r.V2), nil
}

type TimestampLessThanEqualsPathRule struct {
	V1 Path
	V2 Path
}

func (r TimestampLessThanEqualsPathRule) Eval(coj *CtxObj, input interface{}) (bool, error) {
	v1, err := GetTimestamp(coj, input, r.V1)
	if err != nil {
		return false, err
	}

	v2, err := GetTimestamp(coj, input, r.V2)
	if err != nil {
		return false, err
	}

	return v1.LessThanEquals(v2), nil
}

type TimestampGreaterThanRule struct {
	V1 Path
	V2 Timestamp
}

func (r TimestampGreaterThanRule) Eval(coj *CtxObj, input interface{}) (bool, error) {
	v1, err := GetTimestamp(coj, input, r.V1)
	if err != nil {
		return false, err
	}

	return v1.GreaterThan(r.V2), nil
}

type TimestampGreaterThanPathRule struct {
	V1 Path
	V2 Path
}

func (r TimestampGreaterThanPathRule) Eval(coj *CtxObj, input interface{}) (bool, error) {
	v1, err := GetTimestamp(coj, input, r.V1)
	if err != nil {
		return false, err
	}

	v2, err := GetTimestamp(coj, input, r.V2)
	if err != nil {
		return false, err
	}

	return v1.GreaterThan(v2), nil
}

type TimestampGreaterThanEqualsRule struct {
	V1 Path
	V2 Timestamp
}

func (r TimestampGreaterThanEqualsRule) Eval(coj *CtxObj, input interface{}) (bool, error) {
	v1, err := GetTimestamp(coj, input, r.V1)
	if err != nil {
		return false, err
	}

	return v1.GreaterThanEquals(r.V2), nil
}

type TimestampGreaterThanEqualsPathRule struct {
	V1 Path
	V2 Path
}

func (r TimestampGreaterThanEqualsPathRule) Eval(coj *CtxObj, input interface{}) (bool, error) {
	v1, err := GetTimestamp(coj, input, r.V1)
	if err != nil {
		return false, err
	}

	v2, err := GetTimestamp(coj, input, r.V2)
	if err != nil {
		return false, err
	}

	return v1.GreaterThanEquals(v2), nil
}

type IsNullRule struct {
	V1 Path
}

func (r IsNullRule) Eval(coj *CtxObj, input interface{}) (bool, error) {
	v, err := UnjoinByPath(coj, input, &r.V1)
	if err != nil {
		return false, err
	}

	return v == nil, nil
}

type IsPresentRule struct {
	V1 Path
}

func (r IsPresentRule) Eval(coj *CtxObj, input interface{}) (bool, error) {
	_, err := UnjoinByPath(coj, input, &r.V1)
	if err != nil {
		return false, nil
	}

	return true, nil
}

type IsNumericRule struct {
	V1 Path
}

func (r IsNumericRule) Eval(coj *CtxObj, input interface{}) (bool, error) {
	v, err := UnjoinByPath(coj, input, &r.V1)
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

func (r IsStringRule) Eval(coj *CtxObj, input interface{}) (bool, error) {
	v, err := UnjoinByPath(coj, input, &r.V1)
	if err != nil {
		return false, err
	}

	_, ok := v.(string)
	return ok, nil
}

type IsBooleanRule struct {
	V1 Path
}

func (r IsBooleanRule) Eval(coj *CtxObj, input interface{}) (bool, error) {
	v, err := UnjoinByPath(coj, input, &r.V1)
	if err != nil {
		return false, err
	}

	_, ok := v.(bool)
	return ok, nil
}

type IsTimestampRule struct {
	V1 Path
}

func (r IsTimestampRule) Eval(coj *CtxObj, input interface{}) (bool, error) {
	v, err := UnjoinByPath(coj, input, &r.V1)
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
