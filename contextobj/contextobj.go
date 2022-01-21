package contextobj

import (
	"context"
)

type Key string

var contextObjectKey Key = Key("context object key")

type Obj struct {
	obj map[string]interface{}
}

func New(ctx context.Context) context.Context {
	v := ctx.Value(contextObjectKey)
	if v != nil {
		return ctx
	}

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

func Get(ctx context.Context, key string) (interface{}, bool) {
	v := ctx.Value(contextObjectKey)
	if v == nil {
		return nil, false
	}

	obj, ok := v.(Obj)
	if !ok {
		return nil, false
	}

	res, ok := obj.obj[key]
	if !ok {
		return nil, false
	}
	return res, true
}

func GetAll(ctx context.Context) map[string]interface{} {
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

func Del(ctx context.Context, key string) context.Context {
	v := GetAll(ctx)
	delete(v, key)
	return context.WithValue(ctx, contextObjectKey, v)
}
