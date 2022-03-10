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

	cmd.AddCommand(workflowRegisterCmd())
	cmd.AddCommand(workflowRmCmd())
	cmd.AddCommand(workflowExecCmd())

	return cmd
}

func workflowRegisterCmd() *cobra.Command {
	o := cli.RegisterWorkflowOpt{}

	cmd := &cobra.Command{
		Use:   "register",
		Short: "register a workflow",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			if err := o.RegisterWorkflow(ctx, nil); err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.Flags().StringVar(&o.Logfile, "log", "", "path of log files")
	cmd.Flags().StringVar(&o.ASL, "asl", "", "path of a ASL file")
	cmd.Flags().StringVar(&o.WorkflowName, "name", "", "workflow name")

	return cmd
}

func workflowRmCmd() *cobra.Command {
	o := cli.RmWorkflowOpt{}

	cmd := &cobra.Command{
		Use:   "rm",
		Short: "remove a workflow",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			if err := o.RemoveWorkflow(ctx, nil); err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.Flags().StringVar(&o.Logfile, "log", "", "path of log files")
	cmd.Flags().StringVar(&o.WorkflowName, "name", "", "workflow name")

	return cmd
}

func workflowExecCmd() *cobra.Command {
	o := cli.ExecWorkflowOneceOpt{}
	logfile := ""

	cmd := &cobra.Command{
		Use:   "exec",
		Short: "exec a workflow",
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
			if o.ExecWorkflowOpt.WorkflowName != "" {
				result, err = o.ExecWorkflow(ctx, nil)
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
	cmd.Flags().StringVar(&o.ExecWorkflowOpt.WorkflowName, "name", "", "workflow name")

	return cmd
}
