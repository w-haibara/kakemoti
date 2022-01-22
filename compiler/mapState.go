package compiler

import "log"

type RawMapState struct {
	CommonState5
	Iterator       ASL    `json:"Iterator"`
	ItemsPath      string `json:"ItemsPath"`
	MaxConcurrency int    `json:"MaxConcurrency"`
}

func (raw RawMapState) decode(name string) (State, error) {
	s, err := raw.CommonState5.decode(name)
	if err != nil {
		return nil, err
	}

	workflow, err := raw.Iterator.compile()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	path, err := NewReferencePath(raw.ItemsPath)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return MapState{
		CommonState5:   s.Common(),
		Iterator:       *workflow,
		ItemsPath:      path,
		MaxConcurrency: raw.MaxConcurrency,
	}, nil
}

type MapState struct {
	CommonState5
	Iterator       Workflow
	ItemsPath      ReferencePath
	MaxConcurrency int
}
