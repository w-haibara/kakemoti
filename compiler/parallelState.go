package compiler

type ParallelState struct {
	CommonState
	Branches       []ASL  `json:"Branches"`
	ResultPath     string `json:"ResultPath"`
	ResultSelector string `json:"ResultSelector"`
	Retry          string `json:"Retry"`
	Catch          string `json:"Catch"`
}
