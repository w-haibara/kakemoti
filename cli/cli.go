package cli

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/w-haibara/kakemoti/compiler"
	"github.com/w-haibara/kakemoti/log"
	"github.com/w-haibara/kakemoti/worker"
)

var tmpWorkflowMap = make(map[string]compiler.Workflow)

type ExecWorkflowOneceOpt struct {
	*RegisterWorkflowOpt
	*ExecWorkflowOpt
}

func (opt ExecWorkflowOneceOpt) ExecWorkflowOnece(ctx context.Context, coj *compiler.CtxObj, logfile string, workflowName string) ([]byte, error) {
	opt.RegisterWorkflowOpt.WorkflowName = workflowName
	opt.ExecWorkflowOpt.WorkflowName = workflowName

	if opt.RegisterWorkflowOpt.Logfile == "" {
		opt.RegisterWorkflowOpt.Logfile = logfile
	}
	if opt.ExecWorkflowOpt.Logfile == "" {
		opt.ExecWorkflowOpt.Logfile = logfile
	}

	if err := opt.RegisterWorkflow(ctx, nil); err != nil {
		return nil, err
	}

	result, err := opt.ExecWorkflow(ctx, coj)
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

func setLogOutput(l *log.Logger, path string) (close func() error) {
	id, err := uuid.NewRandom()
	if err != nil {
		l.Fatal(err)
	}

	if _, err := os.Stat(path); err != nil {
		if err := os.Mkdir(path, os.ModePerm); err != nil {
			l.Fatal(err)
		}
	}

	path = filepath.Join(path, time.Now().Format("dt=20060102"))
	if _, err := os.Stat(path); err != nil {
		if err := os.Mkdir(path, os.ModePerm); err != nil {
			l.Fatal(err)
		}
	}

	f, err := os.Create(filepath.Join(path, id.String()+".log"))
	if err != nil {
		l.Fatal(err)
	}

	w := io.MultiWriter(os.Stderr, f)
	l.Logger.SetOutput(w)

	return func() error {
		return f.Close()
	}
}

func readFile(path string) (*os.File, *bytes.Buffer, error) {
	f, err := os.Open(path) // #nosec G304
	if err != nil {
		return nil, nil, err
	}

	b := new(bytes.Buffer)
	if _, err := b.ReadFrom(f); err != nil {
		return nil, nil, err
	}

	return f, b, nil
}
