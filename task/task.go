package task

import (
	"context"
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

func Do(ctx context.Context, resourceType, resoucePath string, input interface{}) (interface{}, error) {
	f, ok := fnMap[resourceType]
	if !ok {
		return nil, fmt.Errorf("invalid resouce type: %s", resourceType)
	}

	inMap, ok := input.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("can not cast 'input' to map[string]fn.Obj")
	}

	in, ok := inMap[resourceType]
	if !ok {
		return nil, fmt.Errorf("'inObj' not contains the key: %s", resourceType)
	}

	obj, ok := in.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("can not cast 'in' to fn.Obj")
	}

	out, err := f(ctx, resoucePath, obj)
	if err != nil {
		return nil, fmt.Errorf("fn() failed: %v", err)
	}

	return out, nil
}
