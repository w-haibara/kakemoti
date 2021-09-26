package statemachine

type MapState struct {
	CommonState
	Iterator       StateMachine `json:"Iterator"`
	ItemsPath      string       `json:"ItemsPath"`
	MaxConcurrency int64        `json:"MaxConcurrency"`
	ResultPath     string       `json:"ResultPath"`
	ResultSelector string       `json:"ResultSelector"`
	Retry          string       `json:"Retry"`
	Catch          string       `json:"Catch"`
}
