package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

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

func Exec(ctx context.Context, coj *compiler.CtxObj, w compiler.Workflow, input *bytes.Buffer, logger *log.Logger) ([]byte, error) {
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

	out, err := workflow.Exec(ctx, coj, in)
	if !errors.Is(err, ErrStateMachineTerminated) && err != nil {
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
	return w.loggerWithInfo().WithField("line", log.Line()).WithFields(logrus.Fields{
		"Type": s.Common().Type,
		"Name": s.Name(),
		"Next": s.Next(),
	})
}

func (w Workflow) Exec(ctx context.Context, coj *compiler.CtxObj, input interface{}) (interface{}, error) {
	ctx, cancel := context.WithCancel(ctx)
	if w.TimeoutSeconds > 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Second*time.Duration(w.TimeoutSeconds))
	}
	defer cancel()

	output := input
	branch := w.States[0]
	for {
		out, b, err := w.evalBranch(ctx, coj, branch, output)
		if errors.Is(err, ErrStateMachineTerminated) {
			return out, err
		}
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

func (w Workflow) evalBranch(ctx context.Context, coj *compiler.CtxObj, branch []compiler.State, input interface{}) (interface{}, []compiler.State, error) {
	output := input
	for _, state := range branch {
		out, next, err := w.evalStateWithRetryAndCatch(ctx, coj, state, output)
		w.loggerWithStateInfo(state).WithFields(logrus.Fields{
			"_input":  input,
			"_output": out,
		}).Println()
		if errors.Is(err, ErrStateMachineTerminated) {
			return out, nil, err
		}
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

func (w Workflow) evalStateWithRetryAndCatch(ctx context.Context, coj *compiler.CtxObj, state compiler.State, input interface{}) (interface{}, string, error) {
	origresult, next, origerr := w.evalStateWithFilter(ctx, coj, state, input)
	if origerr.IsEmpty() {
		return origresult, next, nil
	}

	w.loggerWithStateInfo(state).Printf("%s failed: %s", state.Name(), origerr.Error())

	if state.FieldsType() < compiler.FieldsType5 {
		return origresult, next, origerr
	}

	result, next, stateserr := w.retry(ctx, coj, state, input, state.Common().Retry, origerr)
	if stateserr.IsEmpty() {
		return result, next, nil
	}

	return w.catch(ctx, coj, state, input, origresult, origerr)
}

func (w Workflow) retry(ctx context.Context, coj *compiler.CtxObj, state compiler.State, input interface{}, retry []compiler.Retry, stateserr statesError) (interface{}, string, statesError) {
	for _, retry := range retry {
		maxAttempts := 3
		if retry.MaxAttempts != nil {
			maxAttempts = *retry.MaxAttempts
		}

		backoffRate := 2.0
		if retry.BackoffRate != nil {
			backoffRate = *retry.BackoffRate
		}

		intervalSeconds := 1
		if retry.IntervalSeconds != nil {
			intervalSeconds = *retry.IntervalSeconds
		}

		for count := 0; count < maxAttempts; count++ {
			if !func() bool {
				for _, target := range retry.ErrorEquals {
					switch target {
					case StatesErrorALL, stateserr.statesErr, "":
						return true
					}
				}
				return false
			}() {
				break
			}

			ind := float64(intervalSeconds)
			if count > 0 {
				ind += math.Pow(backoffRate, float64(count))
			}

			w.loggerWithStateInfo(state).WithFields(
				logrus.Fields{
					"retry-interval": ind,
					"retry-count":    count,
				}).Println("retry:", state.Name())
			r, n, err := w.retryWithInterval(ctx, coj, state, input, ind)
			if err.IsEmpty() {
				return r, n, err
			}

			w.loggerWithStateInfo(state).Printf("%s failed: %v", state.Name(), err)

			if count == maxAttempts-1 {
				return r, n, err
			}
		}
	}

	err := errors.New("retry() failed")
	return nil, "", NewStatesError(err.Error(), err)
}

func (w Workflow) retryWithInterval(ctx context.Context, coj *compiler.CtxObj, state compiler.State, input interface{}, interval float64) (interface{}, string, statesError) {
	time.Sleep(time.Duration(interval) * time.Second)
	return w.evalState(ctx, coj, state, input)
}

func (w Workflow) catch(ctx context.Context, coj *compiler.CtxObj, state compiler.State, input, result interface{}, stateserr statesError) (interface{}, string, error) {
	if state.FieldsType() < compiler.FieldsType5 {
		return result, "", stateserr
	}

	common := state.Common()
	for _, catch := range common.Catch {
		for _, target := range catch.ErrorEquals {
			if target != StatesErrorALL && target != stateserr.statesErr {
				continue
			}

			if catch.ResultPath == nil {
				return input, catch.Next, nil
			}

			v, err := compiler.JoinByPath(coj, input, result, catch.ResultPath)
			if err != nil {
				return nil, "", err
			}

			return v, catch.Next, nil
		}

	}

	return result, "", stateserr
}

func (w Workflow) evalStateWithFilter(ctx context.Context, coj *compiler.CtxObj, state compiler.State, rawinput interface{}) (interface{}, string, statesError) {
	w.loggerWithStateInfo(state).Println("eval state:", state.Name())

	effectiveInput, stateerr := func() (interface{}, statesError) {
		v1, err := compiler.FilterByInputPath(coj, state, rawinput)
		if err != nil {
			return nil, NewStatesError("", fmt.Errorf("FilterByInputPath(state, rawinput) failed: %v", err))
		}

		v2, err := compiler.FilterByParameters(ctx, coj, state, v1)
		if err != nil {
			if errors.Is(err, compiler.ErrIntrinsicFunctionFailed) {
				return nil, NewStatesError(StatesErrorIntrinsicFailure, err)
			}
			return nil, NewStatesError(StatesErrorParameterPathFailure, fmt.Errorf("FilterByParameters(state, input) failed: %v", err))
		}

		return v2, NewStatesError("", nil)
	}()
	if !stateerr.IsEmpty() {
		return nil, "", stateerr
	}

	result, next, stateerr := w.evalState(ctx, coj, state, effectiveInput)
	if errors.Is(stateerr, ErrStateMachineTerminated) {
		return result, "", stateerr
	}
	if !stateerr.IsEmpty() {
		return nil, "", stateerr
	}

	effectiveResult, stateerr := func() (interface{}, statesError) {
		v1, err := compiler.FilterByResultSelector(ctx, coj, state, result)
		if err != nil {
			if errors.Is(err, compiler.ErrIntrinsicFunctionFailed) {
				return nil, NewStatesError(StatesErrorIntrinsicFailure, err)
			}
			return nil, NewStatesError("", fmt.Errorf("FilterByResultSelector(state, result) failed: %v", err))
		}

		v2, err := compiler.FilterByResultPath(coj, state, rawinput, v1)
		if err != nil {
			return nil, NewStatesError(StatesErrorResultPathMatchFailure, fmt.Errorf("FilterByResultPath(state, rawinput, result) failed: %v", err))
		}

		return v2, NewStatesError("", nil)
	}()
	if !stateerr.IsEmpty() {
		return nil, "", stateerr
	}

	effectiveOutput, err := compiler.FilterByOutputPath(coj, state, effectiveResult)
	if err != nil {
		return nil, "", NewStatesError("", err)
	}

	return effectiveOutput, next, NewStatesError("", nil)
}

func (w Workflow) evalState(ctx context.Context, coj *compiler.CtxObj, state compiler.State, input interface{}) (interface{}, string, statesError) {
	var wg sync.WaitGroup
	wg.Add(1)

	canceled := make(chan bool, 1)
	go func() {
		<-ctx.Done()
		defer wg.Done()

		canceled <- true
		w.loggerWithStateInfo(state).Println("w.evalState() canceled.")
	}()

	var (
		next   string
		output interface{}
		err    statesError
	)

	go func() {
		defer wg.Done()

		switch v := state.(type) {
		case compiler.PassState:
			output, err = w.evalPass(ctx, v, input)
		case compiler.TaskState:
			output, err = w.evalTask(ctx, v, input)
		case compiler.ChoiceState:
			next, output, err = w.evalChoice(ctx, coj, v, input)
		case compiler.WaitState:
			output, err = w.evalWait(ctx, coj, v, input)
		case compiler.SucceedState:
			output, err = w.evalSucceed(ctx, v, input)
		case compiler.FailState:
			output, err = w.evalFail(ctx, v, input)
		case compiler.ParallelState:
			output, err = w.evalParallel(ctx, coj, v, input)
		case compiler.MapState:
			output, err = w.evalMap(ctx, coj, v, input)
		default:
			panic(fmt.Sprintf("unknow state type: %#v", v))
		}
	}()

	wg.Wait()

	select {
	case <-canceled:
		return nil, "", NewStatesError(StatesErrorTimeout, nil)
	default:
	}

	return output, next, err
}

func (w Workflow) nextBranch(state compiler.State) ([]compiler.State, error) {
	if state.Next() == "" {
		return nil, nil
	}

	return w.nextBranchFromString(state.Next())
}

func (w Workflow) nextBranchFromString(next string) ([]compiler.State, error) {
	index, ok := w.StatesIndexMap[next]
	if !ok {
		return nil, fmt.Errorf("the state name is not in the Workflow.StatesIndexMap: %s", next)
	}

	return w.States[index[0]][index[1]:], nil
}
