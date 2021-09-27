package cmd

import (
	"bytes"
	"log"
	"os"
	"strings"

	"karage/statemachine"

	"github.com/spf13/cobra"
)

func NewStartExecutionCmd() *cobra.Command {
	type Options struct {
		Input string
		ASL   string
	}

	o := new(Options)

	cmd := &cobra.Command{
		Use:   "start-execution",
		Short: "Starts a statemachine execution",
		Run: func(cmd *cobra.Command, args []string) {
			if strings.TrimSpace(o.Input) == "" {
				log.Panic("input option value is empty")
			}

			if strings.TrimSpace(o.ASL) == "" {
				log.Panic("ASL option value is empty")
			}

			r, err := readFile(o.Input)
			if err != nil {
				log.Panic(err.Error())
			}

			asl, err := readFile(o.ASL)
			if err != nil {
				log.Panic(err.Error())
			}

			sm, err := statemachine.NewStateMachine(asl, log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile))
			if err != nil {
				log.Panic(err.Error())
			}

			//sm.PrintInfo()
			//sm.PrintStates()

			log.Println("===  First input  ===", "\n"+r.String())

			w := new(bytes.Buffer)
			if err := sm.Start(r, w); err != nil {
				log.Panic(err.Error())
			}

			log.Println("=== Finaly output ===", "\n"+w.String())

		},
	}

	cmd.Flags().StringVar(&o.Input, "input", "", "path of a input json file")
	cmd.Flags().StringVar(&o.ASL, "asl", "", "path of a ASL file")

	return cmd
}

func readFile(path string) (*bytes.Buffer, error) {
	f, err := os.Open(path) // #nosec G304
	if err != nil {
		return nil, err
	}

	b := new(bytes.Buffer)
	if _, err := b.ReadFrom(f); err != nil {
		return nil, err
	}

	return b, nil
}
