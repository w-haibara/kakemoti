package compiler

type RawWaitState struct {
	CommonState3
	Seconds       *int64  `json:"Seconds"`
	Timestamp     *string `json:"Timestamp"`
	SecondsPath   *string `json:"SecondsPath"`
	TimestampPath *string `json:"TimestampPath"`
}

func (state RawWaitState) decode(name string) (State, error) {
	s, err := state.CommonState3.decode(name)
	if err != nil {
		return nil, err
	}

	res := WaitState{
		CommonState3: s.Common().CommonState3,
	}

	if state.Timestamp != nil {
		v, err := NewTimestamp(*state.Timestamp)
		if err != nil {
			return nil, err
		}
		res.Timestamp = &v
	}

	if state.SecondsPath != nil {
		v, err := NewPath(*state.SecondsPath)
		if err != nil {
			return nil, err
		}
		res.SecondsPath = &v
	}

	if state.TimestampPath != nil {
		v, err := NewPath(*state.TimestampPath)
		if err != nil {
			return nil, err
		}
		res.TimestampPath = &v
	}

	return res, nil
}

type WaitState struct {
	CommonState3
	Seconds       *int64
	Timestamp     *Timestamp
	SecondsPath   *Path
	TimestampPath *Path
}
