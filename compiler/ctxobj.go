package compiler

type CtxObj struct {
	v interface{}
}

func (c *CtxObj) SetByString(path string, val interface{}) (*CtxObj, error) {
	p, err := NewPath(path)
	if err != nil {
		return nil, err
	}

	return c.Set(&p, val)
}

func (c *CtxObj) Set(path *Path, val interface{}) (*CtxObj, error) {
	if c.v == nil {
		c.v = make(map[string]interface{})
	}

	v, err := JoinByPath(c, c.v, val, path)
	if err != nil {
		return nil, err
	}

	return &CtxObj{v}, nil
}

func (c *CtxObj) GetByString(path string) (interface{}, bool) {
	p, err := NewPath(path)
	if err != nil {
		return nil, false
	}
	return c.Get(&p)
}

func (c *CtxObj) Get(path *Path) (interface{}, bool) {
	v, err := UnjoinByPath(c, c.v, path)
	if err != nil {
		return nil, false
	}

	return v, true
}

func (c *CtxObj) GetAll() interface{} {
	return c.v
}

func (c *CtxObj) Del(key string) {
	/*
		c.mu.Lock()
		c.mu.Unlock()
	*/
}
