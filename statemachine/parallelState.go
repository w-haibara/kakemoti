package statemachine

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log"
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
	Logger         *log.Logger    `json:"-"`
}

type outputs struct {
	mu sync.Mutex
	v  []*bytes.Buffer
}

func (s *ParallelState) Transition(r, w *bytes.Buffer) (next string, err error) {
	if s == nil {
		return "", nil
	}

	var eg errgroup.Group
	var outputs outputs
	for _, sm := range s.Branches {
		sm := sm
		eg.Go(func() error {
			var buf bytes.Buffer
			defer s.Logger.Println(buf.String())

			logger := log.New(&buf, "	", log.Ldate|log.Ltime|log.Lshortfile)
			logger.Println("=== parallel workflow ===")

			sm.CompleteStateMachine(logger)
			r2 := new(bytes.Buffer)
			w2 := new(bytes.Buffer)
			if _, err := r2.Write(r.Bytes()); err != nil {
				return err
			}

			if ok := ValidateJSON(r2); !ok {
				logger.Println("=== invalid json input ===", "\n"+r.String())
				return err
			}
			logger.Println("===  First input  ===", "\n"+r2.String())

			if err := sm.Start(r2, w2); err != nil {
				return err
			}

			logger.Println("=== Finaly output ===", "\n"+w2.String())

			outputs.mu.Lock()
			outputs.v = append(outputs.v, w2)
			outputs.mu.Unlock()

			s.Logger.Println("\n" + buf.String())

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
