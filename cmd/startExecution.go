package cmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"github.com/w-haibara/kuirejo/cli"
)

func NewStartExecutionCmd() *cobra.Command {
	o := cli.Options{}

	cmd := &cobra.Command{
		Use:   "start-execution",
		Short: "Starts a statemachine execution",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			if _, err := cli.StartExecution(ctx, o); err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.Flags().StringVar(&o.Logfile, "log", "", "path of log files")
	cmd.Flags().StringVar(&o.Input, "input", "", "path of a input json file")
	cmd.Flags().StringVar(&o.ASL, "asl", "", "path of a ASL file")
	cmd.Flags().Int64Var(&o.Timeout, "timeout", 0, "timeout of a statemachine")

	return cmd
}
