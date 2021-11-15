package statemachine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spyzhov/ajson"
	"github.com/w-haibara/kuirejo/log"
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

type StateMachine struct {
	ID             *string                    `json:"-"`
	Comment        *string                    `json:"Comment"`
	StartAt        *string                    `json:"StartAt"`
	TimeoutSeconds *int64                     `json:"TimeoutSeconds"`
	Version        *string                    `json:"Version"`
	RawStates      map[string]json.RawMessage `json:"States"`
	States         States                     `json:"-"`
	Logger         *log.Logger                `json:"-"`
}

type States map[string]State

func NewStateMachine(asl *bytes.Buffer, logger *log.Logger) (*StateMachine, error) {
	dec := json.NewDecoder(asl)

	sm := new(StateMachine)
	if err := sm.setID(); err != nil {
		return nil, err
	}
	sm.Logger = logger

	if err := dec.Decode(sm); err != nil {
		return nil, err
	}

	states, err := sm.decodeStates()
	if err != nil {
		return nil, err
	}
	sm.States = states

	if err := sm.check(); err != nil {
		return nil, err
	}

	return sm, nil
}

func Start(ctx context.Context, asl, input *bytes.Buffer, timeout int64, logger *log.Logger) ([]byte, error) {
	ctx, cancel := context.WithCancel(ctx)
	if timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Second*time.Duration(timeout))
	}
	defer cancel()

	sm, err := NewStateMachine(asl, logger)
	if err != nil {
		logger.Fatalln(err)
	}

	sm.Logger = logger

	b, err := sm.Start(ctx, input)
	if err != nil {
		logger.Fatalln(err)
	}

	return b, nil
}

func (sm *StateMachine) check() error {
	if sm.States == nil {
		return fmt.Errorf("Top-level fields: 'States' is needed")
	}

	if sm.StartAt == nil {
		return fmt.Errorf("Top-level fields: 'StartAt' is needed")
	}

	if sm.Version == nil {
		sm.Version = new(string)
		*sm.Version = "1.0"
	}

	if sm.TimeoutSeconds == nil {
		sm.TimeoutSeconds = new(int64)
		*sm.TimeoutSeconds = 0
	}

	return nil
}

func (sm *StateMachine) Start(ctx context.Context, input *bytes.Buffer) ([]byte, error) {
	if input == nil || strings.TrimSpace(input.String()) == "" {
		input = bytes.NewBuffer(EmptyJSON)
	}

	in, err := ajson.Unmarshal(input.Bytes())
	if err != nil {
		sm.loggerWithSMInfo().Fatalln(err)
	}

	out, err := sm.start(ctx, in)
	if err != nil {
		sm.loggerWithSMInfo().Fatalln(err)
	}

	b, err := ajson.Marshal(out)
	if err != nil {
		sm.loggerWithSMInfo().Fatalln(err)
	}

	return b, nil
}

func (sm *StateMachine) start(ctx context.Context, input *ajson.Node) (*ajson.Node, error) {
	if sm == nil {
		return nil, ErrRecieverIsNil
	}

	if _, ok := sm.States[*sm.StartAt]; !ok {
		return nil, ErrInvalidStartAtValue
	}

	if sm.ID == nil {
		if err := sm.setID(); err != nil {
			return nil, err
		}
	}

	for i := range sm.States {
		sm.States[i].SetID(*sm.ID)
	}

	cur := *sm.StartAt
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

	if node, err := s.FilterInput(ctx, input); err != nil {
		return "", nil, fmt.Errorf("failed to FilterInput(): %v", err)
	} else {
		input = node
	}

	next, output, err := s.Transition(ctx, input)

	if node, err := s.FilterOutput(ctx, output); err != nil {
		return "", nil, fmt.Errorf("failed to FilterOutput(): %v", err)
	} else {
		output = node
	}

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

func (sm *StateMachine) loggerWithSMInfo() *logrus.Entry {
	return sm.Logger.WithFields(logrus.Fields{
		"id":      sm.ID,
		"startat": sm.StartAt,
		"timeout": sm.TimeoutSeconds,
		"line":    log.Line(),
	})
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
	default:
		return nil, ErrUnknownStateName
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

	str := id.String()
	sm.ID = &str

	return nil
}

func NewStates() map[string]State {
	return map[string]State{}
}
