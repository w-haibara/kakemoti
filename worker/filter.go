package worker

import (
	"context"
	"encoding/json"
	"fmt"
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

func ResolvePayloaByJsonPath(ctx context.Context, payload map[string]interface{}, input interface{}, key, path string) (map[string]interface{}, error) {
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

		if strings.HasPrefix(path, "$") {
			v, err := ResolvePayloaByJsonPath(ctx, out, input, key, path)
			if err != nil {
				return nil, err
			}
			out = v
			continue
		}

		return nil, fmt.Errorf("invalid value of payload template: %v", path)
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

	return payload2, err
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
