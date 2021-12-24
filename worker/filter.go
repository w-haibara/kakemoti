package worker

import (
	"fmt"

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
