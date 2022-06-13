package worker

import (
	"bytes"
	"context"

	"github.com/w-haibara/kakemoti/controller/compiler"
	"github.com/w-haibara/kakemoti/worker/workflow"
)

func ExecWorkflow(ctx context.Context, coj *compiler.CtxObj, w compiler.Workflow, input *bytes.Buffer) ([]byte, error) {
	return workflow.Exec(ctx, coj, w, input)
}
