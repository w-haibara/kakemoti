package compiler

type RawPassState struct {
	CommonState4
	Result interface{} `json:"Result"`
}

func (raw RawPassState) decode(name string) (State, error) {
	s, err := raw.CommonState4.decode(name)
	if err != nil {
		return nil, err
	}
	raw.CommonState4 = s.Common().CommonState4
	return PassState(raw), nil
}

type PassState struct {
	CommonState4
	Result interface{}
}
