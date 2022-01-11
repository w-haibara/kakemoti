package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/sen"
	"github.com/w-haibara/kakemoti/compiler"
	"github.com/w-haibara/kakemoti/contextobj"
)

func JoinByJsonPath(ctx context.Context, v1, v2 interface{}, path string) (interface{}, error) {
	if strings.HasPrefix(path, "$$") {
		return JoinByJsonPath(ctx, v1, contextobj.Get(ctx), strings.TrimPrefix(path, "$"))
	}

	p, err := jp.ParseString(path)
	if err != nil {
		return nil, fmt.Errorf("jp.ParseString(v.ResultPath) failed: %v", err)
	}

	if err := p.Set(v1, v2); err != nil {
		return nil, fmt.Errorf("path.Set(rawinput, result) failed: %v", err)
	}

	return v1, nil
}

func UnjoinByJsonPath(ctx context.Context, v interface{}, path string) (interface{}, error) {
	if strings.HasPrefix(path, "$$") {
		return UnjoinByJsonPath(ctx, contextobj.Get(ctx), strings.TrimPrefix(path, "$"))
	}

	p, err := jp.ParseString(path)
	if err != nil {
		return nil, fmt.Errorf("jp.ParseString(v.InputPath) failed: %v", err)
	}

	nodes := p.Get(v)
	if len(nodes) != 1 {
		return nil, fmt.Errorf("invalid length of path.Get(input) result")
	}

	return nodes[0], nil
}

func FilterByInputPath(ctx context.Context, state compiler.State, input interface{}) (interface{}, error) {
	if state.Body.FieldsType() < compiler.FieldsType2 {
		return input, nil
	}

	v := state.Body.Common().CommonState2
	if v.InputPath == "" {
		return input, nil
	}

	return UnjoinByJsonPath(ctx, input, v.InputPath)
}

func FilterByResultPath(ctx context.Context, state compiler.State, rawinput, result interface{}) (interface{}, error) {
	if state.Body.FieldsType() < compiler.FieldsType4 {
		return result, nil
	}

	v := state.Body.Common().CommonState4
	if v.ResultPath == "" {
		return result, nil
	}

	return JoinByJsonPath(ctx, rawinput, result, v.ResultPath)
}

func FilterByOutputPath(ctx context.Context, state compiler.State, output interface{}) (interface{}, error) {
	if state.Body.FieldsType() < compiler.FieldsType2 {
		return output, nil
	}

	v := state.Body.Common().CommonState2
	if v.OutputPath == "" {
		return output, nil
	}

	return UnjoinByJsonPath(ctx, output, v.OutputPath)
}

func SetObjectByKey(v1, v2 interface{}, key string) (interface{}, error) {
	if err := jp.C(key).Set(v1, v2); err != nil {
		return nil, fmt.Errorf("jp.N(0).C(key).Set(v1, v2) failed: %v", err)
	}

	return v1, nil
}

func resolvePayload(ctx context.Context, input interface{}, payload map[string]interface{}) (map[string]interface{}, error) {
	out := make(map[string]interface{})
	for key, val := range payload {
		temp, ok := val.(map[string]interface{})
		if !ok {
			out[key] = val
			continue
		}

		v, err := ResolvePayload(ctx, input, temp)
		if err != nil {
			return nil, err
		}
		out[key] = v
	}
	return out, nil
}

func resolvePayloadByJsonPath(ctx context.Context, input interface{}, payload map[string]interface{}) (map[string]interface{}, error) {
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

		got, err := UnjoinByJsonPath(ctx, input, path)
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

var ErrParseFailed = errors.New("parseIntrinsicFunction() failed")

func parseIntrinsicFunction(ctx context.Context, fnstr string, input interface{}) (string, []interface{}, error) {
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
	fn, argsstr, err := fnAndArgsStr()
	if err != nil {
		return "", nil, err
	}

	resolvePath := func(path string) (interface{}, error) {
		v, err := UnjoinByJsonPath(ctx, input, path)
		if err != nil {
			return nil, err
		}
		return v, nil
	}

	parseArg := func(str string) (interface{}, error) {
		b1 := strings.HasPrefix(str, "'")
		b2 := strings.HasSuffix(str, "'")
		if !b1 && !b2 {
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

			return nil, ErrParseFailed
		}
		if (b1 && !b2) || (!b1 && b2) {
			return nil, ErrParseFailed
		}

		return strings.TrimPrefix(strings.TrimSuffix(str, "'"), "'"), nil
	}

	parseArgs := func(str string) ([]interface{}, error) {
		args := strings.Split(str, ",")
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
	args, err := parseArgs(argsstr)
	if err != nil {
		return "", nil, err
	}

	return fn, args, nil
}

func resolveIntrinsicFunction(ctx context.Context, input interface{}, payload map[string]interface{}) (map[string]interface{}, error) {
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

		fn, args, err := parseIntrinsicFunction(ctx, fnstr, input)
		if err != nil {
			return nil, err
		}

		// TODO: implement instrinsic function
		result, err := func(fn string, args ...interface{}) (string, error) {
			return fn + ": " + fmt.Sprint(args...), nil
		}(fn, args)
		if err != nil {
			return nil, err
		}

		out[strings.TrimSuffix(key, ".$")] = result
	}

	return out, nil
}

func ResolvePayload(ctx context.Context, input interface{}, payload map[string]interface{}) (interface{}, error) {
	payload1, err := resolvePayload(ctx, input, payload)
	if err != nil {
		return nil, err
	}

	payload2, err := resolvePayloadByJsonPath(ctx, input, payload1)
	if err != nil {
		return nil, err
	}

	payload3, err := resolveIntrinsicFunction(ctx, input, payload2)
	if err != nil {
		return nil, err
	}

	return payload3, err
}

func FilterByParameters(ctx context.Context, state compiler.State, input interface{}) (interface{}, error) {
	if state.Body.FieldsType() < compiler.FieldsType4 {
		return input, nil
	}

	v := state.Body.Common().CommonState4
	if v.Parameters == nil {
		return input, nil
	}

	str := ""
	if err := json.Unmarshal(*v.Parameters, &str); err == nil {
		return input, nil
	}

	parameter := make(map[string]interface{})
	if err := json.Unmarshal(*v.Parameters, &parameter); err != nil {
		return nil, fmt.Errorf("json.Unmarshal(*v.Parameters, &selector) failed: %v", err)
	}

	return ResolvePayload(ctx, input, parameter)
}

func FilterByResultSelector(ctx context.Context, state compiler.State, result interface{}) (interface{}, error) {
	if state.Body.FieldsType() < compiler.FieldsType5 {
		return result, nil
	}

	v := state.Body.Common()
	if v.ResultSelector == nil {
		return result, nil
	}

	selector := make(map[string]interface{})
	if err := json.Unmarshal(*v.ResultSelector, &selector); err != nil {
		return nil, fmt.Errorf("json.Unmarshal(*v.ResultSelector, &selector) failed: %v", err)
	}

	return ResolvePayload(ctx, result, selector)
}

func GenerateEffectiveResult(ctx context.Context, state compiler.State, rawinput, result interface{}) (interface{}, error) {
	v1, err := FilterByResultSelector(ctx, state, result)
	if err != nil {
		return nil, fmt.Errorf("FilterByResultSelector(state, result) failed: %v", err)
	}

	v2, err := FilterByResultPath(ctx, state, rawinput, v1)
	if err != nil {
		return nil, fmt.Errorf("FilterByResultPath(state, rawinput, result) failed: %v", err)
	}

	return v2, nil
}

func GenerateEffectiveInput(ctx context.Context, state compiler.State, input interface{}) (interface{}, error) {
	v1, err := FilterByInputPath(ctx, state, input)
	if err != nil {
		return nil, fmt.Errorf("FilterByInputPath(state, rawinput) failed: %v", err)
	}

	v2, err := FilterByParameters(ctx, state, v1)
	if err != nil {
		return nil, fmt.Errorf("FilterByParameters(state, input) failed: %v", err)
	}

	return v2, nil
}
