package worker

import (
	"bytes"
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spyzhov/ajson"
	"github.com/w-haibara/kuirejo/compiler"
	"github.com/w-haibara/kuirejo/log"
)

var ErrStateMachineTerminated = errors.New("state machine terminated")

var (
	EmptyJSON = []byte("{}")
)

func Exec(ctx context.Context, w compiler.Workflow, input *bytes.Buffer, logger *log.Logger) ([]byte, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	workflow := Workflow{&w, id.String(), logger}

	if input == nil || strings.TrimSpace(input.String()) == "" {
		input = bytes.NewBuffer(EmptyJSON)
	}

	in, err := ajson.Unmarshal(input.Bytes())
	if err != nil {
		workflow.loggerWithInfo().Println(err)
		return nil, err
	}

	out, err := workflow.exec(ctx, in)
	if err != nil {
		workflow.loggerWithInfo().Println(err)
		return nil, err
	}

	b, err := ajson.Marshal(out)
	if err != nil {
		workflow.loggerWithInfo().Println(err)
		return nil, err
	}

	return b, nil
}

type Workflow struct {
	*compiler.Workflow
	ID     string
	Logger *log.Logger
}

func (w Workflow) loggerWithInfo() *logrus.Entry {
	return w.Logger.WithFields(logrus.Fields{
		"id":      w.ID,
		"startat": w.StartAt,
		"timeout": w.TimeoutSeconds,
		"line":    log.Line(),
	})
}

func (w Workflow) exec(ctx context.Context, input *ajson.Node) (*ajson.Node, error) {
	return input, nil
}
