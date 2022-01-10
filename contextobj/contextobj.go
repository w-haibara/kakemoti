package contextobj

type ContextObj struct {
	obj map[string]interface{}
}

func NewContextObj() ContextObj {
	return ContextObj{make(map[string]interface{})}
}

func (c ContextObj) Set(key string, val interface{}) {
	c.obj[key] = val
}

func (c ContextObj) Get() map[string]interface{} {
	return c.obj
}
