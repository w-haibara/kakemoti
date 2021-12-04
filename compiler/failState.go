package compiler

type FailState struct {
}

func (state FailState) GetNext() string {
	return ""
}
