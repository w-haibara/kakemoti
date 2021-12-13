package task

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/spyzhov/ajson"
)

type (
	Fn    func(context.Context, string, Obj) (Obj, error)
	FnMap map[string]Fn
	Obj   map[string]interface{}
)

var fnMap FnMap

func init() {
	fnMap = FnMap{
		"script": func(ctx context.Context, path string, in Obj) (Obj, error) {
			return in, nil
		},
	}
}

func Do(ctx context.Context, resouceType, resoucePath string, input *ajson.Node) (*ajson.Node, error) {
	fn, ok := fnMap[resouceType]
	if !ok {
		return nil, errors.New("")
	}

	in, err := unmarshal(input)
	if err != nil {
		return nil, errors.New("")
	}

	out, err := fn(ctx, resoucePath, in)
	if err != nil {
		return nil, errors.New("")
	}

	output, err := marshal(out)
	if err != nil {
		return nil, errors.New("")
	}

	return output, nil
}

func unmarshal(node *ajson.Node) (Obj, error) {
	b, err := ajson.Marshal(node)
	if err != nil {
		return nil, errors.New("")
	}

	in := new(Obj)
	if err := json.Unmarshal(b, in); err != nil {
		return nil, errors.New("")
	}

	return *in, nil
}

func marshal(obj Obj) (*ajson.Node, error) {
	b, err := json.Marshal(obj)
	if err != nil {
		return nil, errors.New("")
	}

	node, err := ajson.Unmarshal(b)
	if err != nil {
		return nil, errors.New("")
	}

	return node, nil
}
