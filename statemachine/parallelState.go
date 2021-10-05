package statemachine

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"
)

type ParallelState struct {
	CommonState
	Branches       []StateMachine `json:"Branches"`
	ResultPath     string         `json:"ResultPath"`
	ResultSelector string         `json:"ResultSelector"`
	Retry          string         `json:"Retry"`
	Catch          string         `json:"Catch"`
}

type outputs struct {
	mu sync.Mutex
	v  []*bytes.Buffer
}

func (s *ParallelState) Transition(ctx context.Context, r, w *bytes.Buffer) (next string, err error) {
	if s == nil {
		return "", nil
	}

	select {
	case <-ctx.Done():
		return "", ErrStoppedStateMachine
	default:
	}

	var eg errgroup.Group
	var outputs outputs
	for _, sm := range s.Branches {
		sm := sm
		sm.Logger = s.GetLogger()

		eg.Go(func() error {
			r2 := new(bytes.Buffer)
			if _, err := r2.Write(r.Bytes()); err != nil {
				return err
			}
			if ok := ValidateJSON(r2); !ok {
				return err
			}

			w2 := new(bytes.Buffer)
			if err := sm.start(ctx, r2, w2); err != nil {
				return err
			}

			outputs.mu.Lock()
			outputs.v = append(outputs.v, w2)
			outputs.mu.Unlock()

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return "", err
	}

	if err := func() error {
		buf := bufio.NewWriter(w)
		if _, err := buf.WriteRune('['); err != nil {
			return err
		}

		for i, output := range outputs.v {
			if _, err := output.WriteTo(buf); err != nil {
				return err
			}

			if i < len(outputs.v)-1 {
				if _, err := buf.WriteRune(','); err != nil {
					return err
				}
			}
		}

		if _, err := buf.WriteRune(']'); err != nil {
			return err
		}

		if err := buf.Flush(); err != nil {
			return err
		}

		w1 := new(bytes.Buffer)
		if err := json.Indent(w1, w.Bytes(), "", "\t"); err != nil {
			return err
		}

		w = w1

		return nil
	}(); err != nil {
		return "", err
	}

	if s.End {
		return "", ErrEndStateMachine
	}

	if strings.TrimSpace(s.Next) == "" {
		return "", ErrNextStateIsBrank
	}

	return s.Next, nil
}
