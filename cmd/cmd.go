package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/w-haibara/kakemoti/cmd/workflow"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kakemoti",
		Short: "Orchestration tool for scripts",
	}

	w := workflow.Workflow{}
	cmd = w.AddCmd(cmd)

	return cmd
}

func Execute() {
	cmd := NewCmd()
	cmd.SetOutput(os.Stdout)
	if err := cmd.Execute(); err != nil {
		cmd.SetOutput(os.Stderr)
		cmd.Println(err)
		os.Exit(1)
	}
}
