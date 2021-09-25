package main

import (
	"bytes"
	"log"
	"strings"
	"sync"

	"github.com/k0kubun/pp"
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

func (s ParallelState) Transition(r, w *bytes.Buffer) (next string, err error) {
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

	w.WriteString("[\n")
	for i, output := range outputs.v {
		w.WriteString("	")
		output.WriteTo(w)
		if i < len(outputs.v)-1 {
			w.WriteString(",\n")
		}
	}
	w.WriteString("\n]")

	if s.End {
		return "", ErrEndStateMachine
	}

	if strings.TrimSpace(s.Next) == "" {
		return "", ErrNextStateIsBrank
	}

	return s.Next, nil
}

func (s ParallelState) Print() {
	pp.Println(s)
}
