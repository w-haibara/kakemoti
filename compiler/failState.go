package compiler

type FailState struct {
	CommonState1
	Cause string `json:"Cause"`
	Error string `json:"Error"`
}
