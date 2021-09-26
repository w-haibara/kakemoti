package cmd

import (
	"bytes"
	"log"
	"os"

	"karage/statemachine"

	"github.com/spf13/cobra"
)

func NewCmdRun() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "* * *",
		Run: func(cmd *cobra.Command, args []string) {
			const (
				input1 = `{
				"IsHelloWorldExample": true
			}`
				input2 = `{
				"IsHelloWorldExample": false
			}`
				input3 = `{
				"IsHelloWorldExample": true,
				"Seconds": 5
			}`

				input4 = `{
					"IsHelloWorldExample": true,
					"Timestamp": "2021-09-25T21:14:10Z"
				}`
			)

			sm, err := statemachine.NewStateMachine("./workflow.asl.json", log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile))
			if err != nil {
				log.Panic("error:", err)
			}

			//sm.PrintInfo()
			//sm.PrintStates()

			r := new(bytes.Buffer)
			w := new(bytes.Buffer)
			if _, err := r.WriteString(input1); err != nil {
				log.Panic("error:", err)
			}

			log.Println("===  First input  ===", "\n"+r.String())

			if err := sm.Start(r, w); err != nil {
				log.Panic("error:", err)
			}

			log.Println("=== Finaly output ===", "\n"+w.String())

		},
	}

	return cmd
}
