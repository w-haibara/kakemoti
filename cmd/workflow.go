package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/w-haibara/kakemoti/cli"
)

func NewWorkflowCmd() *cobra.Command {
	return &cobra.Command{
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
}

func NewWorkflowExecCmd() *cobra.Command {
	o := cli.Options{}

	cmd := &cobra.Command{
		Use:   "exec",
		Short: "exec workflow",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			result, err := cli.StartExecution(ctx, nil, o)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Fprintln(os.Stdout, string(result))
		},
	}

	cmd.Flags().StringVar(&o.Logfile, "log", "", "path of log files")
	cmd.Flags().StringVar(&o.Input, "input", "", "path of a input json file")
	cmd.Flags().StringVar(&o.ASL, "asl", "", "path of a ASL file")
	cmd.Flags().IntVar(&o.Timeout, "timeout", 0, "timeout of a statemachine")

	return cmd
}

func init() {
	cmd := NewWorkflowCmd()
	cmd.AddCommand(NewWorkflowExecCmd())
	rootCmd.AddCommand(cmd)
}
