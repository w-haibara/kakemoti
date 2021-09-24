package main

import ()

type ParallelState struct {
	CommonState
	Branches       []StateMachine `json:"Branches"`
	ResultPath     string         `json:"ResultPath"`
	ResultSelector string         `json:"ResultSelector"`
	Retry          string         `json:"Retry"`
	Catch          string         `json:"Catch"`
}
