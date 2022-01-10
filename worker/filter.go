package worker

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ohler55/ojg/jp"
	"github.com/w-haibara/kakemoti/compiler"
)

func FilterByInputPath(state compiler.State, input interface{}) (interface{}, error) {
	if state.Body.FieldsType() >= compiler.FieldsType2 {
		v := state.Body.Common().CommonState2
		if v.InputPath != "" {
			path, err := jp.ParseString(v.InputPath)
			if err != nil {
				return nil, fmt.Errorf("jp.ParseString(v.InputPath) failed: %v", err)
			}
			nodes := path.Get(input)
			if len(nodes) != 1 {
				return nil, fmt.Errorf("invalid length of path.Get(input) result")
			}
			return nodes[0], nil
		}
	}
	return input, nil
}

func FilterByPayloadTemplate(state compiler.State, input interface{}, template map[string]interface{}) (interface{}, error) {
	out := make([]interface{}, 1)
	out[0] = make(map[string]interface{})
	for key, val := range template {
		if !strings.HasSuffix(key, ".$") {
			out[0].(map[string]interface{})[key] = val
			continue
		}

		s, ok := val.(string)
		if !ok {
			continue
		}

		p, err := jp.ParseString(s)
		if err != nil {
			return nil, fmt.Errorf("jp.ParseString(s) failed: %v", err)
		}
		got := p.Get(input)
		if len(got) < 1 {
			return nil, fmt.Errorf("p.Get(input) failed")
		}

		if err := jp.N(0).C(strings.TrimSuffix(key, ".$")).Set(out, got[0]); err != nil {
			return nil, fmt.Errorf("jp.N(0).C(strings.TrimSuffix(key, \".$\")).Set(selector, p.Get(result)) failed: %v", err)
		}
	}

	return out[0], nil
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

func FilterByResultPath(state compiler.State, rawinput, result interface{}) (interface{}, error) {
	if state.Body.FieldsType() >= compiler.FieldsType4 {
		v := state.Body.Common().CommonState4
		if v.ResultPath != "" {
			path, err := jp.ParseString(v.ResultPath)
			if err != nil {
				return nil, fmt.Errorf("jp.ParseString(v.ResultPath) failed: %v", err)
			}
			if err := path.Set(rawinput, result); err != nil {
				return nil, fmt.Errorf("path.Set(rawinput, result) failed: %v", err)
			}
			return rawinput, nil
		}
	}
	return result, nil
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

func FilterByOutputPath(state compiler.State, output interface{}) (interface{}, error) {
	if state.Body.FieldsType() >= compiler.FieldsType2 {
		v := state.Body.Common().CommonState2
		if v.OutputPath != "" {
			path, err := jp.ParseString(v.OutputPath)
			if err != nil {
				return nil, fmt.Errorf("jp.ParseString(v.OutputPath) failed: %v", err)
			}
			nodes := path.Get(output)
			if len(nodes) != 1 {
				return nil, fmt.Errorf("invalid length of path.Get(output) result")
			}
			return nodes[0], nil
		}
	}
	return output, nil
}
