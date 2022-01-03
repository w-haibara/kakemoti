package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/w-haibara/kakemoti/cli"
)

func NewStartExecutionCmd() *cobra.Command {
	o := cli.Options{}

	cmd := &cobra.Command{
		Use:   "start-execution",
		Short: "Starts a statemachine execution",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			result, err := cli.StartExecution(ctx, o)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Fprintln(os.Stdout, string(result))
		},
	}

	cmd.Flags().StringVar(&o.Logfile, "log", "", "path of log files")
	cmd.Flags().StringVar(&o.Input, "input", "", "path of a input json file")
	cmd.Flags().StringVar(&o.ASL, "asl", "", "path of a ASL file")
	cmd.Flags().Int64Var(&o.Timeout, "timeout", 0, "timeout of a statemachine")

	return cmd
}
