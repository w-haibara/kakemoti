package compiler

type WaitState struct {
}

func (state WaitState) GetNext() string {
	return ""
}
