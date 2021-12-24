package worker

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ohler55/ojg/jp"
	"github.com/w-haibara/kuirejo/compiler"
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

func GenerateEffectiveInput(state compiler.State, rawinput interface{}) (interface{}, error) {
	input, err := FilterByInputPath(state, rawinput)
	if err != nil {
		return nil, fmt.Errorf("FilterByInputPath(state, rawinput) failed: %v", err)
	}
	return input, nil
}

func FilterByResultSelector(state compiler.State, result interface{}) (interface{}, error) {
	if state.Body.FieldsType() >= compiler.FieldsType5 {
		v := state.Body.Common()
		if v.ResultSelector != nil {
			selector := make(map[string]interface{})
			if err := json.Unmarshal(*v.ResultSelector, &selector); err != nil {
				return nil, fmt.Errorf("json.Unmarshal(*v.ResultSelector, &selector) failed: %v", err)
			}

			for key, val := range selector {
				if strings.HasSuffix(key, ".$") {
					s, ok := val.(string)
					if !ok {
						continue
					}
					p, err := jp.ParseString(s)
					if err != nil {
						return nil, fmt.Errorf("jp.ParseString(s) failed: %v", err)
					}
					if err := jp.N(0).C(strings.TrimSuffix(key, ".$")).Set(selector, p.Get(result)); err != nil {
						return nil, fmt.Errorf("jp.N(0).C(strings.TrimSuffix(key, \".$\")).Set(selector, p.Get(result)) failed: %v", err)
					}
				}
			}
		}
	}
	return result, nil
}

func FilterByResultPath(state compiler.State, rawinput, result interface{}) (interface{}, error) {
	if state.Body.FieldsType() >= compiler.FieldsType4 {
		v := state.Body.Common().CommonState4
		if v.OutputPath != "" {
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
