package compiler

type RawState interface {
	decode(name string) (State, error)
}

type State interface {
	Name() string
	Next() string
	FieldsType() int
	Common() CommonState5
}
