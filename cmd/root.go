package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "karage",
		Short: "Orchestration tool for scripts",
	}

	cmd.AddCommand(NewStartExecutionCmd())
	return cmd
}

func Execute() {
	cmd := NewRootCmd()
	cmd.SetOutput(os.Stdout)
	if err := cmd.Execute(); err != nil {
		cmd.SetOutput(os.Stderr)
		cmd.Println(err)
		os.Exit(1)
	}
}
