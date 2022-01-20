package compiler

type RawSucceedState struct {
	CommonState2
}

func (raw RawSucceedState) decode(name string) (State, error) {
	s, err := raw.CommonState2.decode(name)
	if err != nil {
		return nil, err
	}
	raw.CommonState2 = s.Common().CommonState2
	return SucceedState(raw), nil
}

type SucceedState struct {
	CommonState2
}
