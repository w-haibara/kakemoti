package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/k0kubun/pp"
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

	out, err := workflow.Exec(ctx, in)
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

func (w Workflow) Exec(ctx context.Context, input interface{}) (interface{}, error) {
	branch := w.States[0]
	for {
		index, err := w.evalBranch(ctx, branch)
		if err != nil {
			return nil, err
		}
		if index == [2]int{} {
			break
		}

		branch = w.States[index[0]][index[1]:]
	}

	return input, nil
}

func (w Workflow) evalBranch(ctx context.Context, branch []compiler.State) ([2]int, error) {
	for _, state := range branch {
		if err := w.evalState(ctx, state); err != nil {
			return [2]int{}, err
		}
	}

	return w.StatesIndexMap[branch[len(branch)-1].Next], nil
}

func (w Workflow) evalState(ctx context.Context, state compiler.State) error {
	_, _ = pp.Println(state.Name, state.Type)

	return nil
}
