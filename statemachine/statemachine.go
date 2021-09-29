package statemachine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

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
	Comment        string                 `json:"Comment"`
	StartAt        string                 `json:"StartAt"`
	TimeoutSeconds int64                  `json:"TimeoutSeconds"`
	Version        int64                  `json:"Version"`
	RawStates      map[string]interface{} `json:"States"`
	States         map[string]State       `json:"-"`
}

func NewStateMachine(asl *bytes.Buffer) (*StateMachine, error) {
	dec := json.NewDecoder(asl)

	sm := new(StateMachine)
	if err := dec.Decode(sm); err != nil {
		return nil, err
	}

	sm.SetStates()

	return sm, nil
}

func (sm *StateMachine) SetStates() {
	if sm == nil {
		return
	}

	states := map[string]State{}
	for name, state := range sm.RawStates {
		s, ok := state.(map[string]interface{})
		if !ok {
			continue
		}

		t, ok := s["Type"].(string)
		if !ok {
			continue
		}

		convert := func(src, dst interface{}) error {
			var buf bytes.Buffer
			enc := json.NewEncoder(&buf)
			if err := enc.Encode(src); err != nil {
				return err
			}

			dec := json.NewDecoder(&buf)
			if err := dec.Decode(&dst); err != nil {
				return err
			}

			return nil
		}

		switch t {
		case "Pass":
			states[name] = new(PassState)
			if err := convert(s, states[name]); err != nil {
				continue
			}
		case "Task":
			states[name] = new(TaskState)
			if err := convert(s, states[name]); err != nil {
				continue
			}
		case "Choice":
			states[name] = new(ChoiceState)
			if err := convert(s, states[name]); err != nil {
				continue
			}
		case "Wait":
			states[name] = new(WaitState)
			if err := convert(s, states[name]); err != nil {
				continue
			}
		case "Succeed":
			states[name] = new(SucceedState)
			if err := convert(s, states[name]); err != nil {
				continue
			}
		case "Fail":
			states[name] = new(FailState)
			if err := convert(s, states[name]); err != nil {
				continue
			}
		case "Parallel":
			v := new(ParallelState)
			if err := convert(s, v); err != nil {
				continue
			}
			states[name] = v
		case "Map":
			states[name] = new(MapState)
			if err := convert(s, states[name]); err != nil {
				continue
			}
		}
	}

	sm.States = states
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

func (sm *StateMachine) Start(ctx context.Context, r, w *bytes.Buffer) error {
	if sm == nil {
		return ErrRecieverIsNil
	}

	return sm.start(ctx, r, w)
}

func (sm *StateMachine) start(ctx context.Context, r, w *bytes.Buffer) error {
	if sm == nil {
		return ErrRecieverIsNil
	}

	if _, ok := sm.States[sm.StartAt]; !ok {
		return ErrInvalidStartAtValue
	}

	cur := sm.StartAt
	for {
		s, ok := sm.States[cur]
		if !ok {
			return ErrUnknownStateName
		}

		if ok := ValidateJSON(r); !ok {
			return ErrInvalidJSONInput
		}

		s.StateStartLog(cur)
		next, err := s.Transition(ctx, r, w)
		s.StateEndLog(cur)

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
	return nil
}
