package task

import (
	"context"
	"fmt"

	"github.com/w-haibara/kakemoti/task/fn"
)

type (
	Fn    func(context.Context, string, fn.Obj) (fn.Obj, string, error)
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

func Do(ctx context.Context, resourceType, resoucePath string, input interface{}) (interface{}, string, error) {
	f, ok := fnMap[resourceType]
	if !ok {
		return nil, "", fmt.Errorf("invalid resouce type: %s", resourceType)
	}

	var in fn.Obj
	switch v := input.(type) {
	case map[string]interface{}:
		in = v
	case fn.Obj:
		in = v
	default:
		return nil, "", fmt.Errorf("invalid input type: %T, %#v", input, input)
	}

	out, stateserr, err := f(ctx, resoucePath, in)
	if stateserr != "" {
		return nil, stateserr, fmt.Errorf("fn() failed: %s", stateserr)
	}
	if err != nil {
		return nil, "", fmt.Errorf("fn() failed: %v", err)
	}

	return out, "", nil
}
