package compiler

type MapState struct {
	CommonState5
	Iterator       ASL    `json:"Iterator"`
	ItemsPath      string `json:"ItemsPath"`
	MaxConcurrency int64  `json:"MaxConcurrency"`
}
