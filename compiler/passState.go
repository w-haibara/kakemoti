package compiler

type PassState struct {
	CommonState
	Result     string `json:"Result"`
	ResultPath string `json:"ResultPath"`
	Parameters string `json:"Parameters"`
}
