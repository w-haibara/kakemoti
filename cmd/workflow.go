package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/w-haibara/kakemoti/cli"
)

func init() {
	rootCmd.AddCommand(workflowCmd())
}

func workflowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workflow",
		Short: "",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Unknown command")
			}
			fmt.Println("workflow called")
			return nil
		},
	}

	cmd.AddCommand(workflowExecCmd())

	return cmd
}

func workflowExecCmd() *cobra.Command {
	var name *string
	o := cli.ExecWorkflowOneceOpt{}
	logfile := ""

	cmd := &cobra.Command{
		Use:   "exec",
		Short: "exec workflow",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

			id, err := uuid.NewRandom()
			if err != nil {
				log.Fatal(err)
			}

			var (
				result []byte
			)
			if name != nil {
				opt := o.ExecWorkflowOpt
				opt.WorkflowName = *name
				result, err = opt.ExecWorkflow(ctx, nil)
			} else {
				result, err = o.ExecWorkflowOnce(ctx, nil, logfile, id.String())
			}
			if err != nil {
				log.Fatal(err)
			}

			fmt.Fprintln(os.Stdout, string(result))
		},
	}

	cmd.Flags().StringVar(&logfile, "log", "", "path of log files")
	cmd.Flags().StringVar(&o.Input, "input", "", "path of a input json file")
	cmd.Flags().IntVar(&o.Timeout, "timeout", 0, "timeout of a statemachine")

	cmd.Flags().StringVar(&o.ASL, "asl", "", "path of a ASL file")
	cmd.Flags().StringVar(name, "name", "", "workflow name")

	return cmd
}
