package cmd

import (
	"bytes"
	"context"
	"log"
	"os"
	"strings"
	"time"

	"karage/statemachine"

	"github.com/spf13/cobra"
)

func NewStartExecutionCmd() *cobra.Command {
	type Options struct {
		Input   string
		ASL     string
		Timeout int64
	}

	o := new(Options)

	cmd := &cobra.Command{
		Use:   "start-execution",
		Short: "Starts a statemachine execution",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(context.Background())
			if o.Timeout > 0 {
				ctx, cancel = context.WithTimeout(ctx, time.Second*time.Duration(o.Timeout))
			}
			defer cancel()

			if strings.TrimSpace(o.Input) == "" {
				log.Panic("input option value is empty")
			}

			if strings.TrimSpace(o.ASL) == "" {
				log.Panic("ASL option value is empty")
			}

			f1, r, err := readFile(o.Input)
			if err != nil {
				log.Panic(err.Error())
			}
			defer func() {
				if err := f1.Close(); err != nil {
					log.Panic(err.Error())
				}
			}()

			f2, asl, err := readFile(o.ASL)
			if err != nil {
				log.Panic(err.Error())
			}
			defer func() {
				if err := f2.Close(); err != nil {
					log.Panic(err.Error())
				}
			}()

			sm, err := statemachine.NewStateMachine(asl, log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile))
			if err != nil {
				log.Panic(err.Error())
			}

			//sm.PrintInfo()
			//sm.PrintStates()

			log.Println("===  First input  ===", "\n"+r.String())

			w := new(bytes.Buffer)
			if err := sm.Start(ctx, r, w); err != nil {
				log.Panic(err.Error())
			}

			log.Println("=== Finaly output ===", "\n"+w.String())

		},
	}

	cmd.Flags().StringVar(&o.Input, "input", "", "path of a input json file")
	cmd.Flags().StringVar(&o.ASL, "asl", "", "path of a ASL file")
	cmd.Flags().Int64Var(&o.Timeout, "timeout", 0, "timeout of a statemachine")

	return cmd
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
