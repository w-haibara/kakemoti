package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/w-haibara/kuirejo/compiler"
	"github.com/w-haibara/kuirejo/log"
)

var (
	ErrStateMachineTerminated = errors.New("state machine terminated")
	ErrUnknownStateType       = errors.New("unknown state type")
)

var (
	EmptyJSON = []byte("{}")
)

func Exec(ctx context.Context, w compiler.Workflow, input *bytes.Buffer, logger *log.Logger) ([]byte, error) {
	workflow, err := NewWorkflow(&w, logger)
	if err != nil {
		logger.Println("Error:", err)
	}

	if input == nil || strings.TrimSpace(input.String()) == "" {
		input = bytes.NewBuffer(EmptyJSON)
	}

	var in interface{}
	if err := json.Unmarshal(input.Bytes(), &in); err != nil {
		workflow.errorLog(err)
		return nil, err
	}

	out, err := workflow.exec(ctx, in)
	if err != nil {
		workflow.errorLog(err)
		return nil, err
	}

	b, err := json.Marshal(out)
	if err != nil {
		workflow.errorLog(err)
		return nil, err
	}

	return b, nil
}

type Workflow struct {
	*compiler.Workflow
	ID     string
	Logger *log.Logger
}

func NewWorkflow(w *compiler.Workflow, logger *log.Logger) (*Workflow, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return &Workflow{w, id.String(), logger}, nil
}

func (w Workflow) loggerWithInfo() *logrus.Entry {
	return w.Logger.WithFields(logrus.Fields{
		"id":      w.ID,
		"startat": w.StartAt,
		"timeout": w.TimeoutSeconds,
		"line":    log.Line(),
	})
}

func (w Workflow) errorLog(err error) {
	w.loggerWithInfo().WithField("line", log.Line()).Fatalln("Error:", err)
}

func (w Workflow) loggerWithStateInfo(s compiler.State) *logrus.Entry {
	return w.loggerWithInfo().WithFields(logrus.Fields{
		"Type": s.Type,
		"Name": s.Name,
		"Next": s.Next,
	})
}

func (w Workflow) exec(ctx context.Context, input interface{}) (interface{}, error) {
	o, err := w.execStates(ctx, &w.States, input)
	if err != nil {
		w.errorLog(err)
		return nil, err
	}

	return o, nil
}

func (w Workflow) execStates(ctx context.Context, states *compiler.States, input interface{}) (output interface{}, err error) {
	for i := range *states {
		var branch *compiler.States
		output, branch, err = w.execState(ctx, (*states)[i], input)
		if err != nil {
			return nil, err
		}
		input = output

		if branch != nil {
			return w.execStates(ctx, branch, input)
		}
	}
	return output, nil
}

func (w Workflow) execState(ctx context.Context, state compiler.State, input interface{}) (interface{}, *compiler.States, error) {
	w.loggerWithStateInfo(state).Println("eval state:", state.Name)

	var output interface{}
	if choice, ok := state.Body.(*compiler.ChoiceState); ok {
		next := ""
		next, out, err := w.evalChoice(ctx, choice, input)
		if err != nil {
			w.errorLog(err)
			return nil, nil, err
		}
		s, ok := state.Choices[next]
		if !ok {
			err = fmt.Errorf("'next' key is invalid: %s", next)
			w.errorLog(err)
			return nil, nil, err
		}
		return out, s, nil
	} else {
		out, err := w.eval(ctx, &state, input)
		if err != nil {
			w.errorLog(err)
			return nil, nil, err
		}
		output = out
	}
	return output, nil, nil
}

func (w Workflow) eval(ctx context.Context, state *compiler.State, input interface{}) (interface{}, error) {
	switch body := state.Body.(type) {
	case *compiler.FailState:
		output, err := w.evalFail(ctx, body, input)
		if err != nil {
			w.errorLog(err)
			return nil, err
		}
		return output, nil
	case *compiler.MapState:
		output, err := w.evalMap(ctx, body, input)
		if err != nil {
			w.errorLog(err)
			return nil, err
		}
		return output, nil
	case *compiler.ParallelState:
		output, err := w.evalParallel(ctx, body, input)
		if err != nil {
			w.errorLog(err)
			return nil, err
		}
		return output, nil
	case *compiler.PassState:
		output, err := w.evalPass(ctx, body, input)
		if err != nil {
			w.errorLog(err)
			return nil, err
		}
		return output, nil
	case *compiler.SucceedState:
		output, err := w.evalSucceed(ctx, body, input)
		if err != nil {
			w.errorLog(err)
			return nil, err
		}
		return output, nil
	case *compiler.TaskState:
		output, err := w.evalTask(ctx, body, input)
		if err != nil {
			w.errorLog(err)
			return nil, err
		}
		return output, nil
	case *compiler.WaitState:
		output, err := w.evalWait(ctx, body, input)
		if err != nil {
			w.errorLog(err)
			return nil, err
		}
		return output, nil
	}

	w.errorLog(ErrUnknownStateType)
	return nil, ErrUnknownStateType
}
