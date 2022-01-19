package compiler

import "log"

type RawMapState struct {
	CommonState5
	Iterator       ASL    `json:"Iterator"`
	ItemsPath      string `json:"ItemsPath"`
	MaxConcurrency int64  `json:"MaxConcurrency"`
}

func (raw RawMapState) decode() (*MapState, error) {
	workflow, err := raw.Iterator.compile()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	path, err := NewPath(raw.ItemsPath)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &MapState{
		CommonState5:   raw.CommonState5,
		Iterator:       *workflow,
		ItemsPath:      path,
		MaxConcurrency: raw.MaxConcurrency,
	}, nil
}

type MapState struct {
	CommonState5
	Iterator       Workflow
	ItemsPath      Path
	MaxConcurrency int64
}
