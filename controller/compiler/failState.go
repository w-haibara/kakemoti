package compiler

type RawFailState struct {
	CommonState1
	Cause string `json:"Cause"`
	Error string `json:"Error"`
}

func (raw RawFailState) decode(name string) (State, error) {
	s, err := raw.CommonState1.decode(name)
	if err != nil {
		return nil, err
	}
	raw.CommonState1 = s.Common().CommonState1
	return FailState(raw), nil
}

type FailState struct {
	CommonState1
	Cause string
	Error string
}
