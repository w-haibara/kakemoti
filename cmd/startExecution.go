package cmd

import (
	"context"
	"karage/statemachine"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewStartExecutionCmd() *cobra.Command {
	o := new(statemachine.Options)
	l := statemachine.NewLogger()

	cmd := &cobra.Command{
		Use:   "start-execution",
		Short: "Starts a statemachine execution",
		Run:   executionFn(o, l),
	}

	cmd.Flags().StringVar(&o.Input, "input", "", "path of a input json file")
	cmd.Flags().StringVar(&o.ASL, "asl", "", "path of a ASL file")
	cmd.Flags().Int64Var(&o.Timeout, "timeout", 0, "timeout of a statemachine")

	return cmd
}

func executionFn(opt *statemachine.Options, logger *logrus.Entry) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		close := statemachine.SetLogWriter(logger)
		defer close()

		ctx := context.Background()
		if _, err := statemachine.Start(ctx, logger, opt); err != nil {
			logrus.Fatal(err)
		}
	}
}
