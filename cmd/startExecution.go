package cmd

import (
	"context"
	"karage/statemachine"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewStartExecutionCmd() *cobra.Command {
	o := new(statemachine.Options)

	cmd := &cobra.Command{
		Use:   "start-execution",
		Short: "Starts a statemachine execution",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			if _, err := statemachine.Start(ctx, o); err != nil {
				logrus.Fatal(err)
			}
		},
	}

	cmd.Flags().StringVar(&o.Input, "input", "", "path of a input json file")
	cmd.Flags().StringVar(&o.ASL, "asl", "", "path of a ASL file")
	cmd.Flags().Int64Var(&o.Timeout, "timeout", 0, "timeout of a statemachine")

	return cmd
}
