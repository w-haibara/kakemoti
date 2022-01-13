package compiler

type WaitState struct {
	CommonState3
	Seconds          *int64  `json:"Seconds"`
	Timestamp        *string `json:"Timestamp"`
	RawSecondsPath   *string `json:"SecondsPath"`
	SecondsPath      *Path
	RawTimestampPath *string `json:"TimestampPath"`
	TimestampPath    *Path
}

func (state *WaitState) DecodePath() error {
	if err := state.CommonState3.DecodePath(); err != nil {
		return err
	}

	if state.RawSecondsPath != nil {
		v1, err := NewPath(*state.RawSecondsPath)
		if err != nil {
			return err
		}
		*state.SecondsPath = Path(v1)
	}

	if state.RawTimestampPath != nil {
		v2, err := NewPath(*state.RawTimestampPath)
		if err != nil {
			return err
		}
		*state.TimestampPath = Path(v2)
	}

	return nil
}
