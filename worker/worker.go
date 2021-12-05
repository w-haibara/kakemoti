package worker

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/k0kubun/pp"
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
	o, err := w.execStates(ctx, &w.States, input)
	if err != nil {
		w.loggerWithInfo().Println(err)
		return nil, err
	}

	return o, nil
}

func (w Workflow) execStates(ctx context.Context, states *compiler.States, input *ajson.Node) (output *ajson.Node, err error) {
	for i := range *states {
		if (*states)[i].Type == "Choice" {
			next := ""
			next, output, err = evalChoice(ctx, &(*states)[i], input)
			if err != nil {
				w.loggerWithInfo().Println(err)
				return nil, err
			}
			s, ok := (*states)[i].Choices[next]
			if !ok {
				err = fmt.Errorf("'next' key is invalid: %s", next)
				w.loggerWithInfo().Println(err)
				return nil, err
			}
			return w.execStates(ctx, s, output)
		} else {
			output, err = eval(ctx, &(*states)[i], input)
			if err != nil {
				w.loggerWithInfo().Println(err)
				return nil, err
			}
		}

		input = output
	}

	return output, nil
}

func eval(ctx context.Context, state *compiler.State, input *ajson.Node) (*ajson.Node, error) {
	_, _ = pp.Println(state.Name, "-->", state.Next)
	return input, nil
}

func evalChoice(ctx context.Context, state *compiler.State, input *ajson.Node) (string, *ajson.Node, error) {
	next := "Yes"
	_, _ = pp.Println(state.Name, "-->", next)
	return next, input, nil
}
