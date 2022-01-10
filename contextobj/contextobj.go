package contextobj

import (
	"context"
)

type Key struct {
	key string
}

var contextObjectKey Key = Key{"context object key"}

type Obj struct {
	obj map[string]interface{}
}

func New(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextObjectKey, Obj{map[string]interface{}{}})
}

func Set(ctx context.Context, key string, val interface{}) context.Context {
	v := ctx.Value(contextObjectKey)
	if v == nil {
		return ctx
	}

	obj, ok := v.(Obj)
	if !ok {
		return ctx
	}

	obj.obj[key] = val
	return context.WithValue(ctx, contextObjectKey, obj)
}

func Get(ctx context.Context) map[string]interface{} {
	v := ctx.Value(contextObjectKey)
	if v == nil {
		return nil
	}

	obj, ok := v.(Obj)
	if !ok {
		return nil
	}

	return obj.obj
}
