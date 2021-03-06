package compiler

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/sen"
	"github.com/w-haibara/kakemoti/intrinsic"
)

func JoinByPath(coj *CtxObj, v1, v2 interface{}, path *Path) (interface{}, error) {
	if path.IsContextPath {
		path.IsContextPath = false
		return JoinByPath(coj, v1, coj.GetAll(), path)
	}

	if err := path.Expr.Set(v1, v2); err != nil {
		return nil, fmt.Errorf("path.Set(rawinput, result) failed (path=[%s]) : %v", path, err)
	}

	return v1, nil
}

func UnjoinByPath(coj *CtxObj, v interface{}, path *Path) (interface{}, error) {
	if path.IsContextPath {
		path.IsContextPath = false
		return UnjoinByPath(coj, coj.GetAll(), path)
	}

	nodes := path.Expr.Get(v)
	if len(nodes) != 1 {
		return nil, fmt.Errorf("invalid length of path.Get(input) result (path=[%s])", path)
	}

	return nodes[0], nil
}

func GetString(coj *CtxObj, input interface{}, path Path) (string, error) {
	v, err := UnjoinByPath(coj, input, &path)
	if err != nil {
		return "", err
	}

	v1, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("invalid field value (must be string) : [%s]=[%v]", path.String(), v1)
	}

	return v1, nil
}

func GetNumeric(coj *CtxObj, input interface{}, path Path) (float64, error) {
	v, err := UnjoinByPath(coj, input, &path)
	if err != nil {
		return 0, err
	}

	v1, ok := v.(float64)
	if !ok {
		return 0, fmt.Errorf("invalid field value (must be float64) : [%s]=[%v]", path.String(), v1)
	}

	return v1, nil
}

func GetBool(coj *CtxObj, input interface{}, path Path) (bool, error) {
	v, err := UnjoinByPath(coj, input, &path)
	if err != nil {
		return false, err
	}

	v1, ok := v.(bool)
	if !ok {
		return false, fmt.Errorf("invalid field value (must be boolean) : [%s]=[%v]", path.String(), v1)
	}

	return v1, nil
}

func GetTimestamp(coj *CtxObj, input interface{}, path Path) (Timestamp, error) {
	v, err := GetString(coj, input, path)
	if err != nil {
		return Timestamp{}, err
	}

	v1, err := NewTimestamp(v)
	if err != nil {
		return Timestamp{}, err
	}

	return v1, nil
}

func FilterByInputPath(coj *CtxObj, state State, input interface{}) (interface{}, error) {
	if state.FieldsType() < FieldsType2 {
		return input, nil
	}

	v := state.Common().CommonState2
	if v.InputPath == nil {
		return input, nil
	}

	return UnjoinByPath(coj, input, v.InputPath)
}

func FilterByResultPath(coj *CtxObj, state State, rawinput, result interface{}) (interface{}, error) {
	if state.FieldsType() < FieldsType4 {
		return result, nil
	}

	v := state.Common().CommonState4
	if v.ResultPath == nil {
		return result, nil
	}

	return JoinByPath(coj, rawinput, result, &v.ResultPath.Path)
}

func FilterByOutputPath(coj *CtxObj, state State, output interface{}) (interface{}, error) {
	if state.FieldsType() < FieldsType2 {
		return output, nil
	}

	v := state.Common().CommonState2
	if v.OutputPath == nil {
		return output, nil
	}

	return UnjoinByPath(coj, output, v.OutputPath)
}

func SetObjectByKey(v1, v2 interface{}, key string) (interface{}, error) {
	if err := jp.C(key).Set(v1, v2); err != nil {
		return nil, fmt.Errorf("jp.N(0).C(key).Set(v1, v2) failed: %v", err)
	}

	return v1, nil
}

func resolvePayload(ctx context.Context, coj *CtxObj, input interface{}, payload map[string]interface{}) (map[string]interface{}, error) {
	out := make(map[string]interface{})
	for key, val := range payload {
		temp, ok := val.(map[string]interface{})
		if !ok {
			out[key] = val
			continue
		}

		v, err := ResolvePayload(ctx, coj, input, temp)
		if err != nil {
			return nil, err
		}
		out[key] = v
	}
	return out, nil
}

func resolvePayloadByPath(coj *CtxObj, input interface{}, payload map[string]interface{}) (map[string]interface{}, error) {
	out := make(map[string]interface{})
	for key, val := range payload {
		if !strings.HasSuffix(key, ".$") {
			out[key] = val
			continue
		}

		path, ok := val.(string)
		if !ok {
			return nil, fmt.Errorf("value of payload template is not string: %v", path)
		}

		if !strings.HasPrefix(path, "$") {
			out[key] = path
			continue
		}

		p, err := NewPath(path)
		if err != nil {
			return nil, err
		}

		got, err := UnjoinByPath(coj, input, &p)
		if err != nil {
			return nil, err
		}

		v, err := SetObjectByKey(payload, got, strings.TrimSuffix(key, ".$"))
		if err != nil {
			return nil, err
		}

		v1, ok := v.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("result of SetObjectByKey() is invarid: %s", sen.String(v, &ojg.Options{Sort: true}))
		}

		delete(v1, key)
		out = v1
	}

	return out, nil
}

var ErrIntrinsicFunctionFailed = errors.New("intrinsic function failed")

func parseIntrinsicFunction(ctx context.Context, coj *CtxObj, fnstr string, input interface{}) (string, []interface{}, error) {
	var ErrParseFailed = errors.New("parseIntrinsicFunction() failed")

	fnAndArgsStr := func() (string, string, error) {
		n1 := strings.Index(fnstr, "(")
		if n1 < 1 && n1+1 < len(fnstr) {
			return "", "", ErrParseFailed
		}

		n2 := strings.LastIndex(fnstr, ")")
		if n2 < 2 || n1 >= n2 {
			return "", "", ErrParseFailed
		}

		return fnstr[:n1], fnstr[n1+1 : n2], nil
	}

	resolvePath := func(path string) (interface{}, error) {
		p, err := NewPath(path)
		if err != nil {
			return nil, err
		}
		v, err := UnjoinByPath(coj, input, &p)
		if err != nil {
			return nil, err
		}
		return v, nil
	}

	parseArg := func(str string) (interface{}, error) {
		b1 := strings.HasPrefix(str, "'")
		b2 := strings.HasSuffix(str, "'")
		if (b1 && !b2) || (!b1 && b2) {
			return nil, ErrParseFailed
		}

		if b1 && b2 {
			return strings.TrimPrefix(strings.TrimSuffix(str, "'"), "'"), nil
		}

		if strings.HasPrefix(str, "$") {
			v, err := resolvePath(str)
			if err != nil {
				return nil, err
			}
			return v, nil
		}

		switch str {
		case "true":
			return true, nil
		case "false":
			return false, nil
		case "null":
			return "null", nil
		}

		if v, err := strconv.Atoi(str); err == nil {
			return v, nil
		}
		if v, err := strconv.ParseFloat(str, 64); err == nil {
			return v, nil
		}

		fn, args, err := parseIntrinsicFunction(ctx, coj, str, input)
		if err != nil {
			return nil, err
		}
		result, err := intrinsic.Do(ctx, fn, args)
		if err != nil {
			return nil, err
		}

		return result, nil
	}

	parseArgs := func(str string) ([]interface{}, error) {
		args := []string{""}
		parenCount := 0
		quoted := false
		for _, s := range str {
			if (parenCount > 0 && s != '(' && s != ')') || (quoted && s != '\'') {
				args[len(args)-1] += string(s)
				continue
			}

			if s == ',' {
				args = append(args, "")
				continue
			}

			args[len(args)-1] += string(s)

			if s == '(' {
				parenCount++
				continue
			}
			if s == ')' {
				parenCount--
				continue
			}

			if s == '\'' {
				quoted = !quoted
				continue
			}
		}

		result := make([]interface{}, len(args))
		for i, arg := range args {
			arg = strings.TrimSpace(arg)

			s, err := parseArg(arg)
			if err != nil {
				return nil, err
			}

			result[i] = s
		}
		return result, nil
	}

	fn, argsstr, err := fnAndArgsStr()
	if err != nil {
		return "", nil, fmt.Errorf("%s: %w", err.Error(), ErrIntrinsicFunctionFailed)
	}

	args, err := parseArgs(argsstr)
	if err != nil {
		return "", nil, fmt.Errorf("%s: %w", err.Error(), ErrIntrinsicFunctionFailed)
	}

	return fn, args, nil
}

func resolveIntrinsicFunction(ctx context.Context, coj *CtxObj, input interface{}, payload map[string]interface{}) (map[string]interface{}, error) {
	out := make(map[string]interface{})
	for key, val := range payload {
		if !strings.HasSuffix(key, ".$") {
			out[key] = val
			continue
		}

		fnstr, ok := val.(string)
		if !ok {
			return nil, fmt.Errorf("value of payload template is not string: %v", fnstr)
		}

		fn, args, err := parseIntrinsicFunction(ctx, coj, fnstr, input)
		if err != nil {
			return nil, err
		}

		result, err := intrinsic.Do(ctx, fn, args)
		if err != nil {
			return nil, err
		}

		out[strings.TrimSuffix(key, ".$")] = result
	}

	return out, nil
}

func ResolvePayload(ctx context.Context, coj *CtxObj, input interface{}, payload interface{}) (interface{}, error) {
	v, ok := payload.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid type of payload: [%#v]", payload)
	}
	payload1, err := resolvePayload(ctx, coj, input, v)
	if err != nil {
		return nil, err
	}

	payload2, err := resolvePayloadByPath(coj, input, payload1)
	if err != nil {
		return nil, err
	}

	payload3, err := resolveIntrinsicFunction(ctx, coj, input, payload2)
	if err != nil {
		return nil, err
	}

	return payload3, err
}

func FilterByParameters(ctx context.Context, coj *CtxObj, state State, input interface{}) (interface{}, error) {
	if state.FieldsType() < FieldsType4 {
		return input, nil
	}

	v := state.Common().CommonState4
	if v.Parameters == nil {
		return input, nil
	}

	return ResolvePayload(ctx, coj, input, v.Parameters)
}

func FilterByResultSelector(ctx context.Context, coj *CtxObj, state State, result interface{}) (interface{}, error) {
	if state.FieldsType() < FieldsType5 {
		return result, nil
	}

	v := state.Common()
	if v.ResultSelector == nil {
		return result, nil
	}

	return ResolvePayload(ctx, coj, result, v.ResultSelector)
}
