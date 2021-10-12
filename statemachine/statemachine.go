package statemachine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"karage/log"

	"github.com/google/uuid"
	"github.com/k0kubun/pp"
)

var (
	ErrRecieverIsNil         = fmt.Errorf("receiver is nil")
	ErrInvalidStartAtValue   = fmt.Errorf("invalid StateAt value")
	ErrInvalidJSONInput      = fmt.Errorf("invalid json input")
	ErrInvalidJSONOutput     = fmt.Errorf("invalid json output")
	ErrUnknownStateName      = fmt.Errorf("unknown state name")
	ErrUnknownStateType      = fmt.Errorf("unknown state type")
	ErrNextStateIsBrank      = fmt.Errorf("next state is brank")
	ErrSucceededStateMachine = fmt.Errorf("state machine stopped successfully")
	ErrFailedStateMachine    = fmt.Errorf("state machine stopped unsuccessfully")
	ErrEndStateMachine       = fmt.Errorf("end state machine")
	ErrStoppedStateMachine   = fmt.Errorf("stopped state machine")
)

var (
	EmptyJSON = []byte("{}")
)

type StateMachine struct {
	ID             string                     `json:"-"`
	Comment        string                     `json:"Comment"`
	StartAt        string                     `json:"StartAt"`
	TimeoutSeconds int64                      `json:"TimeoutSeconds"`
	Version        int64                      `json:"Version"`
	RawStates      map[string]json.RawMessage `json:"States"`
	States         States                     `json:"-"`
	RootLogger     *log.RootLogger            `json:"-"`
	Logger         *log.Logger                `json:"-"`
}

type States map[string]State

func NewStateMachine(asl *bytes.Buffer) (*StateMachine, error) {
	dec := json.NewDecoder(asl)
	sm := new(StateMachine)
	if err := sm.setID(); err != nil {
		return nil, err
	}
	sm.Logger = log.NewLogger(sm.ID)

	if err := dec.Decode(sm); err != nil {
		return nil, err
	}

	var err error
	sm.States, err = sm.decodeStates()
	if err != nil {
		return nil, err
	}

	return sm, nil
}

func NewStates() map[string]State {
	return map[string]State{}
}

func (sm *StateMachine) decodeStates() (States, error) {
	if sm == nil {
		return nil, ErrRecieverIsNil
	}

	states := NewStates()

	for name, raw := range sm.RawStates {
		state, err := sm.decodeState(raw)
		if err != nil {
			return nil, err
		}

		state.SetName(name)
		state.SetLogger(sm.Logger)

		states[name] = state
	}

	return states, nil
}

func (sm *StateMachine) decodeState(raw json.RawMessage) (State, error) {
	var t struct {
		Type string `json:"Type"`
	}

	if err := json.Unmarshal(raw, &t); err != nil {
		return nil, err
	}

	switch t.Type {
	case "Parallel":
		v := new(ParallelState)
		if err := json.Unmarshal(raw, v); err != nil {
			return nil, err
		}

		for k := range v.Branches {
			v.Branches[k].Logger = sm.Logger

			var err error
			v.Branches[k].States, err = v.Branches[k].decodeStates()
			if err != nil {
				return nil, err
			}
		}

		return v, nil
	}

	var state State
	switch t.Type {
	case "Pass":
		state = new(PassState)
	case "Task":
		state = new(TaskState)
	case "Choice":
		state = new(ChoiceState)
	case "Wait":
		state = new(WaitState)
	case "Succeed":
		state = new(SucceedState)
	case "Fail":
		state = new(FailState)
	case "Map":
		state = new(MapState)
	}

	if err := json.Unmarshal(raw, state); err != nil {
		return nil, err
	}

	return state, nil
}

func ValidateJSON(j *bytes.Buffer) bool {
	b := j.Bytes()

	if len(bytes.TrimSpace(b)) == 0 {
		j.Reset()
		j.Write(EmptyJSON)
		return true
	}

	if !json.Valid(b) {
		return false
	}

	return true
}

func (sm *StateMachine) setID() error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	sm.ID = id.String()

	return nil
}

func (sm *StateMachine) Start(ctx context.Context, r, w *bytes.Buffer) error {
	return sm.start(ctx, r, w)
}

func (sm *StateMachine) start(ctx context.Context, r, w *bytes.Buffer) error {
	if sm == nil {
		return ErrRecieverIsNil
	}

	if _, ok := sm.States[sm.StartAt]; !ok {
		return ErrInvalidStartAtValue
	}

	if sm.ID == "" {
		if err := sm.setID(); err != nil {
			return err
		}
	}

	defer sm.Logger.Close(sm.ID)

	for i := range sm.States {
		sm.States[i].SetID(sm.ID)
	}

	sm.StateMachineStartLog()

	cur := sm.StartAt
	for {
		s, ok := sm.States[cur]
		if !ok {
			return ErrUnknownStateName
		}

		if ok := ValidateJSON(r); !ok {
			return ErrInvalidJSONInput
		}

		s.StateStartLog()
		next, err := s.Transition(ctx, r, w)
		s.StateEndLog()

		if ok := ValidateJSON(w); !ok {
			return ErrInvalidJSONOutput
		}

		switch {
		case err == ErrUnknownStateType:
			return err
		case err == ErrSucceededStateMachine:
			goto End
		case err == ErrFailedStateMachine:
			goto End
		case err == ErrEndStateMachine:
			goto End
		case err != nil:
			return err
		}

		if _, ok := sm.States[next]; !ok {
			return ErrUnknownStateName
		}

		r.Reset()
		if _, err := w.WriteTo(r); err != nil {
			return err
		}
		w.Reset()

		cur = next
	}

End:
	sm.StateMachineEndLog()
	return nil
}

func (sm *StateMachine) Log(v ...interface{}) {
	sm.RootLogger.Println(sm.ID, "", "", fmt.Sprint(v...))
}

func (sm *StateMachine) StateMachineStartLog() {
	sm.Log("START")
}

func (sm *StateMachine) StateMachineEndLog() {
	sm.Log("END")
}

func (sm *StateMachine) PrintInfo() {
	if sm == nil {
		return
	}

	fmt.Println("====== StateMachine Info ======")
	_, _ = pp.Println("Comment", sm.Comment)
	_, _ = pp.Println("StartAt", sm.StartAt)
	_, _ = pp.Println("TimeoutSeconds", sm.TimeoutSeconds)
	_, _ = pp.Println("Version", sm.Version)
	fmt.Println("===============================")
}

func (sm *StateMachine) PrintStates() {
	if sm == nil {
		return
	}

	s := sm.States
	fmt.Println("=========== States  ===========")
	for k, v := range s {
		_, _ = pp.Println(k, "\n", v, "\n")
	}
	fmt.Println("===============================")
}
