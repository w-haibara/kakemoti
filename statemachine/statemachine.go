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

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint: true,
	})
}

func NewStateMachine(asl *bytes.Buffer) (*StateMachine, error) {
	dec := json.NewDecoder(asl)
	sm := new(StateMachine)
	if err := sm.setID(); err != nil {
		return nil, err
	}
	sm.Logger = sm.logger()

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

func Start(ctx context.Context, o *Options) ([]byte, error) {
	close := setLogWriter()
	defer close()

	ctx, cancel := context.WithCancel(ctx)
	if o.Timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Second*time.Duration(o.Timeout))
	}
	defer cancel()

	if strings.TrimSpace(o.Input) == "" {
		logrus.Fatalln("input option value is empty")
	}

	if strings.TrimSpace(o.ASL) == "" {
		logrus.Fatalln("ASL option value is empty")
	}

	f1, input, err := readFile(o.Input)
	if err != nil {
		logrus.Fatalln(err.Error())
	}
	defer func() {
		if err := f1.Close(); err != nil {
			logrus.Fatalln(err.Error())
		}
	}()

	f2, asl, err := readFile(o.ASL)
	if err != nil {
		logrus.Fatalln(err.Error())
	}
	defer func() {
		if err := f2.Close(); err != nil {
			logrus.Fatalln(err.Error())
		}
	}()

	sm, err := NewStateMachine(asl)
	if err != nil {
		logrus.Fatalln(err.Error())
	}

	return sm.Start(ctx, input)
}

func (sm *StateMachine) Start(ctx context.Context, input *bytes.Buffer) ([]byte, error) {
	r, err := ajson.Unmarshal(input.Bytes())
	if err != nil {
		return nil, err
	}

	w, err := sm.start(ctx, r)
	if err != nil {
		return nil, err
	}

	v, err := w.Unpack()
	if err != nil {
		return nil, err
	}

	output, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func (sm *StateMachine) start(ctx context.Context, r *ajson.Node) (*ajson.Node, error) {
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

	l := sm.logger()

	if v, err := r.Unpack(); err != nil {
		return nil, fmt.Errorf("invalid ajson.Node error: %v", err)
	} else {
		l.WithFields(logrus.Fields{
			"output": v,
		}).Info("statemachine start")
	}

	var w *ajson.Node
	cur := sm.StartAt
	for {
		select {
		case <-ctx.Done():
			break
		default:
		}

		s, ok := sm.States[cur]
		if !ok {
			return nil, ErrUnknownStateName
		}

		if v, err := r.Unpack(); err != nil {
			return nil, fmt.Errorf("invalid ajson.Node error: %v", err)
		} else {
			s.Logger().WithFields(logrus.Fields{
				"input": v,
			}).Info("state start")
		}

		var (
			next string
			err  error
		)
		next, w, err = s.Transition(ctx, r)

		if w == nil {
			w = &ajson.Node{}
		}
		if v, err := w.Unpack(); err != nil {
			return nil, fmt.Errorf("invalid ajson.Node error: %v", err)
		} else {
			s.Logger().WithFields(logrus.Fields{
				"output": v,
			}).Info("state end")
		}

		switch {
		case err == ErrSucceededStateMachine:
			goto End
		case err == ErrFailedStateMachine:
			goto End
		case err == ErrEndStateMachine:
			goto End
		case err != nil:
			return nil, err
		}

		if _, ok := sm.States[next]; !ok {
			return nil, ErrUnknownStateName
		}

		r = w
		cur = next
	}

End:
	if v, err := w.Unpack(); err != nil {
		return nil, fmt.Errorf("invalid ajson.Node error: %v", err)
	} else {
		l.WithFields(logrus.Fields{
			"output": v,
		}).Info("statemachine end")
	}

	return w, nil
}

func (sm *StateMachine) logger() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"id":      sm.ID,
		"startat": sm.StartAt,
		"timeout": sm.TimeoutSeconds,
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
