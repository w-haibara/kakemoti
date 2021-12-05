package compiler

type MapState struct {
	CommonState
	Iterator       ASL    `json:"Iterator"`
	ItemsPath      string `json:"ItemsPath"`
	MaxConcurrency int64  `json:"MaxConcurrency"`
	ResultPath     string `json:"ResultPath"`
	ResultSelector string `json:"ResultSelector"`
	Retry          string `json:"Retry"`
	Catch          string `json:"Catch"`
}
