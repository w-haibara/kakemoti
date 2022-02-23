package workflow

import (
	"github.com/spf13/cobra"
)

type Workflow struct {
}

func (w Workflow) AddCmd(cmd *cobra.Command) *cobra.Command {
	cmd.AddCommand(NewStartExecutionCmd())
	return cmd
}
