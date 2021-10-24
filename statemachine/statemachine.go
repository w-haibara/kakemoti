package statemachine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spyzhov/ajson"
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
	ErrInvalidJsonPath       = fmt.Errorf("invalid JsonPath")
	ErrInvalidInputPath      = fmt.Errorf("invalid InputPath")
	ErrInvalidRawJSON        = fmt.Errorf("invalid raw json")
)

var (
	EmptyJSON = []byte("{}")
)

type Options struct {
	Input   string
	ASL     string
	Timeout int64
}

type StateMachine struct {
	ID             string                     `json:"-"`
	Comment        string                     `json:"Comment"`
	StartAt        string                     `json:"StartAt"`
	TimeoutSeconds int64                      `json:"TimeoutSeconds"`
	Version        int64                      `json:"Version"`
	RawStates      map[string]json.RawMessage `json:"States"`
	States         States                     `json:"-"`
	Logger         *logrus.Entry              `json:"-"`
}

type States map[string]State

func NewStateMachine(asl *bytes.Buffer) (*StateMachine, error) {
	dec := json.NewDecoder(asl)
	sm := new(StateMachine)
	if err := sm.setID(); err != nil {
		return nil, err
	}
	sm.Logger = NewLogger()

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

func Start(ctx context.Context, l *logrus.Entry, o *Options) ([]byte, error) {
	ctx, cancel := context.WithCancel(ctx)
	if o.Timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Second*time.Duration(o.Timeout))
	}
	defer cancel()

	if strings.TrimSpace(o.Input) == "" {
		l.Fatalln("input option value is empty")
	}

	if strings.TrimSpace(o.ASL) == "" {
		l.Fatalln("ASL option value is empty")
	}

	f1, input, err := readFile(o.Input)
	if err != nil {
		l.Fatalln(err)
	}
	defer func() {
		if err := f1.Close(); err != nil {
			l.Fatalln(err)
		}
	}()

	f2, asl, err := readFile(o.ASL)
	if err != nil {
		l.Fatalln(err)
	}
	defer func() {
		if err := f2.Close(); err != nil {
			l.Fatalln(err)
		}
	}()

	sm, err := NewStateMachine(asl)
	if err != nil {
		l.Fatalln(err)
	}

	sm.Logger = l

	b, err := sm.Start(ctx, input)
	if err != nil {
		l.Fatalln(err)
	}

	return b, nil
}

func (sm *StateMachine) Start(ctx context.Context, input *bytes.Buffer) ([]byte, error) {
	in, err := ajson.Unmarshal(input.Bytes())
	if err != nil {
		sm.logger(nil).Fatalln(err)
	}

	out, err := sm.start(ctx, in)
	if err != nil {
		sm.logger(nil).Fatalln(err)
	}

	b, err := ajson.Marshal(out)
	if err != nil {
		sm.logger(nil).Fatalln(err)
	}

	return b, nil
}

func (sm *StateMachine) start(ctx context.Context, input *ajson.Node) (*ajson.Node, error) {
	if sm == nil {
		return nil, ErrRecieverIsNil
	}

	if _, ok := sm.States[sm.StartAt]; !ok {
		return nil, ErrInvalidStartAtValue
	}

	if sm.ID == "" {
		if err := sm.setID(); err != nil {
			return nil, err
		}
	}

	for i := range sm.States {
		sm.States[i].SetID(sm.ID)
	}

	cur := sm.StartAt
	for {
		var err error
		cur, input, err = sm.transition(ctx, cur, input)
		switch err {
		case nil:
		case ErrSucceededStateMachine, ErrFailedStateMachine, ErrEndStateMachine:
			return input, nil
		default:
			return nil, err
		}
	}
}

func (sm *StateMachine) transition(ctx context.Context, next string, input *ajson.Node) (string, *ajson.Node, error) {
	select {
	case <-ctx.Done():
		return "", nil, ErrStoppedStateMachine
	default:
	}

	s, ok := sm.States[next]
	if !ok {
		return "", nil, ErrUnknownStateName
	}

	if v, err := input.Unpack(); err != nil {
		return "", nil, fmt.Errorf("invalid ajson.Node error: %v", err)
	} else {
		s.Logger(logrus.Fields{
			"input": v,
		}).Info("state start")
	}

	next, output, err := s.Transition(ctx, input)

	if output == nil {
		output = &ajson.Node{}
	}
	if v, err := output.Unpack(); err != nil {
		return "", nil, fmt.Errorf("invalid ajson.Node error: %v", err)
	} else {
		s.Logger(logrus.Fields{
			"output": v,
		}).Info("state end")
	}

	if err != nil {
		return "", output, err
	}

	return next, output, nil
}

func (sm *StateMachine) logger(v logrus.Fields) *logrus.Entry {
	return sm.Logger.WithFields(logrus.Fields{
		"id":      sm.ID,
		"startat": sm.StartAt,
		"timeout": sm.TimeoutSeconds,
	}).WithFields(v)
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

func (sm *StateMachine) setID() error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	sm.ID = id.String()

	return nil
}

func NewStates() map[string]State {
	return map[string]State{}
}

func readFile(path string) (*os.File, *bytes.Buffer, error) {
	f, err := os.Open(path) // #nosec G304
	if err != nil {
		return nil, nil, err
	}

	b := new(bytes.Buffer)
	if _, err := b.ReadFrom(f); err != nil {
		return nil, nil, err
	}

	return f, b, nil
}
