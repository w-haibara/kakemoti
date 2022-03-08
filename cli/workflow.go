package cli

import (
	"context"
	"strings"

	"github.com/w-haibara/kakemoti/compiler"
	"github.com/w-haibara/kakemoti/log"
	"github.com/w-haibara/kakemoti/worker"
)

var tmpWorkflowMap = make(map[string]compiler.Workflow)

type ExecWorkflowOneceOpt struct {
	*RegisterWorkflowOpt
	*ExecWorkflowOpt
}

func (opt ExecWorkflowOneceOpt) ExecWorkflowOnce(ctx context.Context, coj *compiler.CtxObj, logfile string, workflowName string) ([]byte, error) {
	opt.RegisterWorkflowOpt.WorkflowName = workflowName
	opt.ExecWorkflowOpt.WorkflowName = workflowName

	if opt.RegisterWorkflowOpt.Logfile == "" {
		opt.RegisterWorkflowOpt.Logfile = logfile
	}
	if opt.ExecWorkflowOpt.Logfile == "" {
		opt.ExecWorkflowOpt.Logfile = logfile
	}

	var workflow compiler.Workflow

	rfn := func(name string, w compiler.Workflow) error {
		workflow = w
		return nil
	}
	if err := opt.registerWorkflow(ctx, nil, rfn); err != nil {
		return nil, err
	}

	ffn := func(name string) (compiler.Workflow, error) {
		return workflow, nil
	}
	result, err := opt.execWorkflow(ctx, coj, ffn)
	if err != nil {
		return nil, err
	}

	return result, nil
}

type RegisterWorkflowOpt struct {
	Logfile      string
	ASL          string
	WorkflowName string
}

func (opt RegisterWorkflowOpt) RegisterWorkflow(ctx context.Context, coj *compiler.CtxObj) error {
	return opt.registerWorkflow(ctx, coj,
		func(name string, w compiler.Workflow) error {
			tmpWorkflowMap[name] = w
			return nil
		},
	)
}

type registerWorkflowFunc func(name string, w compiler.Workflow) error

func (opt RegisterWorkflowOpt) registerWorkflow(ctx context.Context, coj *compiler.CtxObj, fn registerWorkflowFunc) error {
	if strings.TrimSpace(opt.Logfile) == "" {
		opt.Logfile = "logs"
	}

	logger := log.NewLogger()
	close := setLogOutput(logger, opt.Logfile)
	defer func() {
		if err := close(); err != nil {
			panic(err.Error())
		}
	}()

	if strings.TrimSpace(opt.ASL) == "" {
		logger.Fatalln("ASL option value is empty")
	}

	f1, asl, err := readFile(opt.ASL)
	if err != nil {
		logger.Fatalln(err)
	}
	defer func() {
		if err := f1.Close(); err != nil {
			logger.Fatalln(err)
		}
	}()

	workflow, err := compiler.Compile(ctx, asl)
	if err != nil {
		logger.Fatalln(err)
	}

	if err := fn(opt.WorkflowName, *workflow); err != nil {
		return err
	}

	return nil
}

type ExecWorkflowOpt struct {
	Logfile      string
	WorkflowName string
	Input        string
	Timeout      int
}

func (opt ExecWorkflowOpt) ExecWorkflow(ctx context.Context, coj *compiler.CtxObj) ([]byte, error) {
	return opt.execWorkflow(ctx, coj,
		func(name string) (compiler.Workflow, error) {
			return tmpWorkflowMap[name], nil
		},
	)
}

type fetchWorkflowFunc func(name string) (compiler.Workflow, error)

func (opt ExecWorkflowOpt) execWorkflow(ctx context.Context, coj *compiler.CtxObj, fn fetchWorkflowFunc) ([]byte, error) {
	if strings.TrimSpace(opt.Logfile) == "" {
		opt.Logfile = "logs"
	}

	logger := log.NewLogger()
	close := setLogOutput(logger, opt.Logfile)
	defer func() {
		if err := close(); err != nil {
			panic(err.Error())
		}
	}()

	if strings.TrimSpace(opt.Input) == "" {
		logger.Fatalln("input option value is empty")
	}

	f2, input, err := readFile(opt.Input)
	if err != nil {
		logger.Fatalln(err)
	}
	defer func() {
		if err := f2.Close(); err != nil {
			logger.Fatalln(err)
		}
	}()

	if coj == nil {
		coj = &compiler.CtxObj{}
	}

	w, err := fn(opt.WorkflowName)
	if err != nil {
		return nil, err
	}

	return worker.Exec(ctx, coj, w, input, logger)
}
