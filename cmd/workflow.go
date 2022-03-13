package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

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
	cmd.AddCommand(workflowListCmd())
	cmd.AddCommand(workflowExecCmd())
	cmd.AddCommand(workflowRmCmd())
	cmd.AddCommand(workflowDropCmd())

	return cmd
}

var MsgMustSpecifyWorkflowName = "You must specify a workflow name."

func workflowRegisterCmd() *cobra.Command {
	o := cli.RegisterWorkflowOpt{}

	cmd := &cobra.Command{
		Use:   "register [WORKFLOW NAME]",
		Short: "register a workflow",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				log.Fatal(MsgMustSpecifyWorkflowName)
			}

			o.WorkflowName = args[0]

			ctx := context.Background()
			if err := o.RegisterWorkflow(ctx, nil); err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.Flags().StringVar(&o.Logfile, "log", "", "path of log files")
	cmd.Flags().StringVar(&o.ASL, "asl", "", "path of a ASL file")
	cmd.Flags().BoolVarP(&o.Force, "force", "f", false, "if the name isn't exists, will update it")

	return cmd
}

func workflowListCmd() *cobra.Command {
	o := cli.ListWorkflowOpt{}
	logfile := ""

	cmd := &cobra.Command{
		Use:   "list",
		Short: "list workflows",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			names, err := o.ListWorkflow()
			if err != nil {
				log.Fatal(err)
			}

			str := ""
			for i, v := range names {
				str += fmt.Sprintln(strconv.Itoa(i)+".", v)
			}

			fmt.Fprintln(os.Stdout, str)
		},
	}

	cmd.Flags().StringVar(&logfile, "log", "", "path of log files")

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

			if len(args) < 1 {
				log.Fatal(MsgMustSpecifyWorkflowName)
			}

			o.ExecWorkflowOpt.WorkflowName = args[0]

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

	return cmd
}

func workflowRmCmd() *cobra.Command {
	o := cli.RemoveWorkflowOpt{}

	cmd := &cobra.Command{
		Use:   "rm",
		Short: "remove a workflow",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

			if len(args) < 1 {
				log.Fatal(MsgMustSpecifyWorkflowName)
			}

			o.WorkflowName = args[0]

			if err := o.RemoveWorkflow(ctx, nil); err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.Flags().StringVar(&o.Logfile, "log", "", "path of log files")
	cmd.Flags().BoolVarP(&o.Force, "force", "f", false, "if the name isn't exists, will update it")

	return cmd
}

func workflowDropCmd() *cobra.Command {
	o := cli.DropWorkflowOpt{}

	cmd := &cobra.Command{
		Use:   "drop",
		Short: "drop the workflows table",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			if err := o.DropWorkflow(ctx, nil); err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.Flags().StringVar(&o.Logfile, "log", "", "path of log files")
	cmd.Flags().BoolVarP(&o.Force, "force", "f", false, "if the name isn't exists, will update it")

	return cmd
}
