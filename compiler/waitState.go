package compiler

type WaitState struct {
	CommonState3
	Seconds          *int64  `json:"Seconds"`
	RawTimestamp     *string `json:"Timestamp"`
	Timestamp        *Timestamp
	RawSecondsPath   *string `json:"SecondsPath"`
	SecondsPath      *Path
	RawTimestampPath *string `json:"TimestampPath"`
	TimestampPath    *Path
}

func (state *WaitState) DecodePath() error {
	if err := state.CommonState3.DecodePath(); err != nil {
		return err
	}

	if state.RawTimestamp != nil {
		v, err := NewTimestamp(*state.RawTimestamp)
		if err != nil {
			return err
		}
		*state.Timestamp = v
	}

	if state.RawSecondsPath != nil {
		v, err := NewPath(*state.RawSecondsPath)
		if err != nil {
			return err
		}
		*state.SecondsPath = v
	}

	if state.RawTimestampPath != nil {
		v, err := NewPath(*state.RawTimestampPath)
		if err != nil {
			return err
		}
		*state.TimestampPath = v
	}

	return nil
}
