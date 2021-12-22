package task

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/w-haibara/kuirejo/task/fn"
)

type (
	Fn    func(context.Context, string, fn.Obj) (fn.Obj, error)
	FnMap map[string]Fn
)

var fnMap FnMap

func init() {
	fnMap = make(FnMap)
	RegisterDefault()
}

func RegisterDefault() {
	Register("script", fn.DoScriptTask)
}

func Register(name string, fn Fn) {
	fnMap[name] = fn
}

func Do(ctx context.Context, resouceType, resoucePath string, input interface{}) (interface{}, error) {
	fn, ok := fnMap[resouceType]
	if !ok {
		return nil, fmt.Errorf("invalid resouce type: %s", resouceType)
	}

	in, err := unmarshal(input)
	if err != nil {
		return nil, fmt.Errorf("unmarshal() failed: %v", err)
	}

	out, err := fn(ctx, resoucePath, in)
	if err != nil {
		return nil, fmt.Errorf("fn() failed: %v", err)
	}

	output, err := marshal(out)
	if err != nil {
		return nil, fmt.Errorf("marshal() failed: %v", err)
	}

	return output, nil
}

func unmarshal(node interface{}) (fn.Obj, error) {
	b, err := json.Marshal(node)
	if err != nil {
		return nil, err
	}

	var in fn.Obj
	if err := json.Unmarshal(b, &in); err != nil {
		return nil, err
	}

	return in, nil
}

func marshal(obj fn.Obj) (interface{}, error) {
	b, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	var node interface{}
	if err := json.Unmarshal(b, &node); err != nil {
		return nil, err
	}

	return node, nil
}
