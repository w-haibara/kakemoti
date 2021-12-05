package compiler

type CommonState struct {
	Type       string `json:"Type"`
	Next       string `json:"Next"`
	End        bool   `json:"End"`
	Comment    string `json:"Comment"`
	InputPath  string `json:"InputPath"`
	OutputPath string `json:"OutputPath"`
}

func (state CommonState) GetNext() string {
	return state.Next
}
