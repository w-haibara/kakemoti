package intrinsic

import (
	"context"
	"fmt"

	"github.com/w-haibara/kakemoti/intrinsic/fn"
)

type (
	Fn    func(context.Context, []interface{}) (interface{}, error)
	FnMap map[string]Fn
)

var fnMap FnMap

func init() {
	fnMap = make(FnMap)
	RegisterDefault()
}

func RegisterDefault() {
	Register("States.Format", fn.DoStatesFormat)
	Register("States.StringToJson", fn.DoStatesStringToJson)
	Register("States.JsonToString", fn.DoStatesJsonToString)
	Register("States.Array", fn.DoStatesArray)
}

func Register(name string, fn Fn) {
	fnMap[name] = fn
}

func Do(ctx context.Context, fnname string, args []interface{}) (interface{}, error) {
	f, ok := fnMap[fnname]
	if !ok {
		return nil, fmt.Errorf("unknown intrinsic function: %s", fnname)
	}

	out, err := f(ctx, args)
	if err != nil {
		return nil, fmt.Errorf("fn() failed: %v", err)
	}

	return out, nil
}
