package compiler

type StateBody interface {
	GetNext() string
	FieldsType() int
	Common() CommonState5
	DecodePath() error
}
