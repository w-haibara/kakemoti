package db

import "sync"

var m = sync.Map{}

func Save(id string, v interface{}) {
	m.Store(id, v)
}

func Find(id string) interface{} {
	v, ok := m.Load(id)
	if !ok {
		panic("m.Load failed")
	}

	return v
}
