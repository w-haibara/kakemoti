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
	"github.com/w-haibara/kakemoti/compiler"
	"github.com/w-haibara/kakemoti/log"
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
	output := input
	branch := w.States[0]
	for {
		out, b, err := w.evalBranch(ctx, branch, output)
		if err != nil {
			return nil, err
		}

		output = out

		if b == nil {
			break
		}

		branch = b
	}

	return output, nil
}

func (w Workflow) evalBranch(ctx context.Context, branch []compiler.State, input interface{}) (interface{}, []compiler.State, error) {
	output := input
	for _, state := range branch {
		out, next, err := w.evalStateWithFilter(ctx, state, output)
		if err != nil {
			return nil, nil, err
		}

		output = out

		if next == "" {
			continue
		}

		b, err := w.nextBranchFromString(next)
		if err != nil {
			return nil, nil, err
		}
		if b != nil {
			return out, b, nil
		}
	}

	branch, err := w.nextBranch(branch[len(branch)-1])
	if err != nil {
		return nil, nil, err
	}

	return output, branch, nil
}

func (w Workflow) evalStateWithFilter(ctx context.Context, state compiler.State, rawinput interface{}) (interface{}, string, error) {
	w.loggerWithStateInfo(state).Println("eval state:", state.Name)

	effectiveInput, err := GenerateEffectiveInput(state, rawinput)
	if err != nil {
		w.errorLog(err)
		return nil, "", err
	}

	result, next, err := w.evalStateWithRetry(ctx, state, effectiveInput)
	if errors.Is(err, ErrStateMachineTerminated) {
		return result, "", err
	}
	if err != nil {
		return nil, "", err
	}

	effectiveResult, err := GenerateEffectiveResult(state, rawinput, result)
	if err != nil {
		w.errorLog(err)
		return nil, "", err
	}

	effectiveOutput, err := FilterByOutputPath(state, effectiveResult)
	if err != nil {
		w.errorLog(err)
		return nil, "", err
	}

	return effectiveOutput, next, nil
}

func (w Workflow) evalStateWithRetry(ctx context.Context, state compiler.State, input interface{}) (interface{}, string, error) {
	result, next, stateserr, err := w.evalState(ctx, state, input)
	if errors.Is(err, ErrStateMachineTerminated) {
		return result, "", err
	}
	if err != nil {
		return nil, "", err
	}

	if stateserr != "" {
		// TODO: implement retry & catch
		return result, "", ErrStateMachineTerminated
	}

	return result, next, nil
}

func (w Workflow) evalState(ctx context.Context, state compiler.State, input interface{}) (interface{}, string, statesError, error) {
	var (
		next   string
		output interface{}
		err    error
	)

	switch body := state.Body.(type) {
	case *compiler.PassState:
		output, err = w.evalPass(ctx, body, input)
	case *compiler.TaskState:
		output, err = w.evalTask(ctx, body, input)
	case *compiler.ChoiceState:
		next, output, err = w.evalChoice(ctx, body, input)
	case *compiler.WaitState:
		output, err = w.evalWait(ctx, body, input)
	case *compiler.SucceedState:
		output, err = w.evalSucceed(ctx, body, input)
	case *compiler.FailState:
		output, err = w.evalFail(ctx, body, input)
	case *compiler.ParallelState:
		output, err = w.evalParallel(ctx, body, input)
	case *compiler.MapState:
		output, err = w.evalMap(ctx, body, input)
	}

	if errors.Is(err, ErrStateMachineTerminated) {
		return output, next, "", nil
	}
	if err != nil {
		return nil, "", "", err
	}

	return output, next, "", nil
}

func (w Workflow) nextBranch(state compiler.State) ([]compiler.State, error) {
	if state.Next == "" {
		return nil, nil
	}

	return w.nextBranchFromString(state.Next)
}

func (w Workflow) nextBranchFromString(next string) ([]compiler.State, error) {
	index, ok := w.StatesIndexMap[next]
	if !ok {
		return nil, fmt.Errorf("the state name is not in the Workflow.StatesIndexMap: %s", next)
	}

	return w.States[index[0]][index[1]:], nil
}
