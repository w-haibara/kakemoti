package task

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spyzhov/ajson"
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

func Do(ctx context.Context, resouceType, resoucePath string, input *ajson.Node) (*ajson.Node, error) {
	fn, ok := fnMap[resouceType]
	if !ok {
		return nil, fmt.Errorf("invalid resouce type: %s", resouceType)
	}

	in, err := unmarshal(input)
	if err != nil {
		return nil, fmt.Errorf("unmarshal() failed: %w", err)
	}

	out, err := fn(ctx, resoucePath, in)
	if err != nil {
		return nil, fmt.Errorf("fn() failed: %w", err)
	}

	output, err := marshal(out)
	if err != nil {
		return nil, fmt.Errorf("marshal() failed: %w", err)
	}

	return output, nil
}

func unmarshal(node *ajson.Node) (fn.Obj, error) {
	b, err := ajson.Marshal(node)
	if err != nil {
		return nil, err
	}

	in := new(fn.Obj)
	if err := json.Unmarshal(b, in); err != nil {
		return nil, err
	}

	return *in, nil
}

func marshal(obj fn.Obj) (*ajson.Node, error) {
	b, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	node, err := ajson.Unmarshal(b)
	if err != nil {
		return nil, err
	}

	return node, nil
}
