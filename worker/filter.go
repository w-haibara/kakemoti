package worker

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/sen"
	"github.com/w-haibara/kakemoti/compiler"
)

func JoinByJsonPath(v1, v2 interface{}, path string) (interface{}, error) {
	p, err := jp.ParseString(path)
	if err != nil {
		return nil, fmt.Errorf("jp.ParseString(v.ResultPath) failed: %v", err)
	}
	if err := p.Set(v1, v2); err != nil {
		return nil, fmt.Errorf("path.Set(rawinput, result) failed: %v", err)
	}
	return v1, nil

}

func UnjoinByJsonPath(v interface{}, path string) (interface{}, error) {
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

func FilterByInputPath(state compiler.State, input interface{}) (interface{}, error) {
	if state.Body.FieldsType() < compiler.FieldsType2 {
		return input, nil
	}

	v := state.Body.Common().CommonState2
	if v.InputPath == "" {
		return input, nil
	}

	return UnjoinByJsonPath(input, v.InputPath)
}

func FilterByResultPath(state compiler.State, rawinput, result interface{}) (interface{}, error) {
	if state.Body.FieldsType() < compiler.FieldsType4 {
		return result, nil
	}

	v := state.Body.Common().CommonState4
	if v.ResultPath == "" {
		return result, nil
	}

	return JoinByJsonPath(rawinput, result, v.ResultPath)
}

func FilterByOutputPath(state compiler.State, output interface{}) (interface{}, error) {
	if state.Body.FieldsType() < compiler.FieldsType2 {
		return output, nil
	}

	v := state.Body.Common().CommonState2
	if v.OutputPath == "" {
		return output, nil
	}

	return UnjoinByJsonPath(output, v.OutputPath)
}

func SetObjectByKey(v1, v2 interface{}, key string) (interface{}, error) {
	if err := jp.C(key).Set(v1, v2); err != nil {
		return nil, fmt.Errorf("jp.N(0).C(key).Set(v1, v2) failed: %v", err)
	}

	return v1, nil
}

func resolveJsonPath(template map[string]interface{}, input interface{}, key, path string) (map[string]interface{}, error) {
	got, err := UnjoinByJsonPath(input, path)
	if err != nil {
		return nil, err
	}

	v, err := SetObjectByKey(template, got, strings.TrimSuffix(key, ".$"))
	if err != nil {
		return nil, err
	}

	v1, ok := v.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("result of SetObjectByKey() is invarid: %s", sen.String(v, &ojg.Options{Sort: true}))
	}

	return v1, nil
}

func FilterByPayloadTemplate(state compiler.State, input interface{}, template map[string]interface{}) (interface{}, error) {
	out := make(map[string]interface{})
	for key, val := range template {
		if !strings.HasSuffix(key, ".$") {
			out[key] = val
			continue
		}

		path, ok := val.(string)
		if !ok {
			return nil, fmt.Errorf("value of payload template is not string: %v", path)
		}

		switch {
		case strings.HasPrefix(path, "$"):
			v, err := resolveJsonPath(out, input, key, path)
			if err != nil {
				return nil, err
			}
			out = v
		default:
			return nil, fmt.Errorf("invalid value of payload template: %v", path)
		}
	}

	return out, nil
}

func FilterByParameters(state compiler.State, input interface{}) (interface{}, error) {
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

	return FilterByPayloadTemplate(state, input, parameter)
}

func FilterByResultSelector(state compiler.State, result interface{}) (interface{}, error) {
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

	return FilterByPayloadTemplate(state, result, selector)
}

func GenerateEffectiveResult(state compiler.State, rawinput, result interface{}) (interface{}, error) {
	var err error
	result, err = FilterByResultSelector(state, result)
	if err != nil {
		return nil, fmt.Errorf("FilterByResultSelector(state, result) failed: %v", err)
	}

	result, err = FilterByResultPath(state, rawinput, result)
	if err != nil {
		return nil, fmt.Errorf("FilterByResultPath(state, rawinput, result) failed: %v", err)
	}

	return result, nil
}

func GenerateEffectiveInput(state compiler.State, input interface{}) (interface{}, error) {
	var err error
	input, err = FilterByInputPath(state, input)
	if err != nil {
		return nil, fmt.Errorf("FilterByInputPath(state, rawinput) failed: %v", err)
	}

	input, err = FilterByParameters(state, input)
	if err != nil {
		return nil, fmt.Errorf("FilterByParameters(state, input) failed: %v", err)
	}

	return input, nil
}
