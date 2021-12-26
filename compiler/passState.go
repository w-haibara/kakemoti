package compiler

type PassState struct {
	CommonState4
	Result interface{} `json:"Result"`
}
