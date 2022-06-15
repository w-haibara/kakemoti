package workflow

import (
	"bytes"
	"context"
	"fmt"

	"github.com/w-haibara/kakemoti/db"
	"github.com/w-haibara/kakemoti/worker"
)

type Server struct {
}

func (s *Server) Start(ctx context.Context, in *StartRequest) (*StartResponce, error) {
	fmt.Println("========= workflow.Exec =========")
	fmt.Println("WorkflowName", in.WorkflowName, "Input:", in.Input)

	workflow, err := db.GetWorkflow(in.WorkflowName)
	if err != nil {
		return nil, err
	}

	wf, err := workflow.DecodeWorkflow()
	if err != nil {
		return nil, err
	}

	res, err := worker.ExecWorkflow(context.TODO(), nil, wf, bytes.NewBufferString(in.Input))
	if err != nil {
		return nil, err
	}

	return &StartResponce{Output: string(res)}, nil
}
