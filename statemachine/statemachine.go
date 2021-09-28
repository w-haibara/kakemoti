package statemachine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/k0kubun/pp"
)

var (
	ErrInvalidStartAtValue   = fmt.Errorf("invalid StateAt value")
	ErrInvalidJSONInput      = fmt.Errorf("invalid json input")
	ErrInvalidJSONOutput     = fmt.Errorf("invalid json output")
	ErrUnknownStateName      = fmt.Errorf("unknown state name")
	ErrUnknownStateType      = fmt.Errorf("unknown state type")
	ErrNextStateIsBrank      = fmt.Errorf("next state is brank")
	ErrSucceededStateMachine = fmt.Errorf("state machine stopped successfully")
	ErrFailedStateMachine    = fmt.Errorf("state machine stopped unsuccessfully")
	ErrEndStateMachine       = fmt.Errorf("end state machine")
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
	Logger         *log.Logger            `json:"-"`
}

func NewStateMachine(asl *bytes.Buffer, logger *log.Logger) (*StateMachine, error) {
	dec := json.NewDecoder(asl)

	sm := new(StateMachine)
	if err := dec.Decode(sm); err != nil {
		return nil, err
	}

	sm.CompleteStateMachine(logger)

	return sm, nil
}

func (sm *StateMachine) CompleteStateMachine(logger *log.Logger) {
	sm.Logger = logger
	sm.SetStates()
}

func (sm *StateMachine) SetStates() {
	states := map[string]State{}
	for name, state := range sm.RawStates {
		s, ok := state.(map[string]interface{})
		if !ok {
			sm.Logger.Println("invalid state definition:", name)
			continue
		}

		t, ok := s["Type"].(string)
		if !ok {
			sm.Logger.Println("invalid type value:", s["Type"])
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
				sm.Logger.Println("error:", err)
				continue
			}
		case "Task":
			states[name] = new(TaskState)
			if err := convert(s, states[name]); err != nil {
				sm.Logger.Println("error:", err)
				continue
			}
		case "Choice":
			states[name] = new(ChoiceState)
			if err := convert(s, states[name]); err != nil {
				sm.Logger.Println("error:", err)
				continue
			}
		case "Wait":
			states[name] = new(WaitState)
			if err := convert(s, states[name]); err != nil {
				sm.Logger.Println("error:", err)
				continue
			}
		case "Succeed":
			states[name] = new(SucceedState)
			if err := convert(s, states[name]); err != nil {
				sm.Logger.Println("error:", err)
				continue
			}
		case "Fail":
			states[name] = new(FailState)
			if err := convert(s, states[name]); err != nil {
				sm.Logger.Println("error:", err)
				continue
			}
		case "Parallel":
			v := new(ParallelState)
			if err := convert(s, v); err != nil {
				sm.Logger.Println("error:", err)
				continue
			}
			v.Logger = sm.Logger
			states[name] = v
		case "Map":
			states[name] = new(MapState)
			if err := convert(s, states[name]); err != nil {
				sm.Logger.Println("error:", err)
				continue
			}
		}
	}

	sm.States = states
}

func (sm StateMachine) PrintInfo() {
	fmt.Println("====== StateMachine Info ======")
	_, _ = pp.Println("Comment", sm.Comment)
	_, _ = pp.Println("StartAt", sm.StartAt)
	_, _ = pp.Println("TimeoutSeconds", sm.TimeoutSeconds)
	_, _ = pp.Println("Version", sm.Version)
	fmt.Println("===============================")
}

func (sm StateMachine) PrintStates() {
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

func (sm StateMachine) Start(r, w *bytes.Buffer) error {
	if _, ok := sm.States[sm.StartAt]; !ok {
		return ErrInvalidStartAtValue
	}

	cur := sm.StartAt
	var err error
	for {
		s, ok := sm.States[cur]
		if !ok {
			sm.Logger.Println("UnknownStateName:", cur)
			return ErrUnknownStateName
		}

		if ok := ValidateJSON(r); !ok {
			sm.Logger.Println("=== invalid json input ===", "\n"+r.String())
			return ErrInvalidJSONInput
		}
		sm.Logger.Println("State:", cur, "( Type =", s.StateType(), ")")
		sm.Logger.Println("=== input  ===", "\n"+r.String())

		cur, err = s.Transition(r, w)

		if ok := ValidateJSON(w); !ok {
			sm.Logger.Println("=== invalid json output ===\n", "\n"+w.String())
			return ErrInvalidJSONOutput
		}
		sm.Logger.Println("=== output ===", "\n"+w.String())

		switch {
		case err == ErrUnknownStateType:
			sm.Logger.Println("UnknownStateType:", cur)
			return err
		case err == ErrSucceededStateMachine:
			sm.Logger.Println(err)
			goto End
		case err == ErrFailedStateMachine:
			sm.Logger.Println(err)
			goto End
		case err == ErrEndStateMachine:
			sm.Logger.Println(err)
			goto End
		case err != nil:
			return err
		}

		if _, ok := sm.States[cur]; !ok {
			sm.Logger.Println("UnknownStateName: [", cur, "]")
			return ErrUnknownStateName
		}

		r.Reset()
		if _, err := w.WriteTo(r); err != nil {
			sm.Logger.Println("WriteTo error:", err)
			return err
		}
		w.Reset()
	}

End:
	return nil
}
