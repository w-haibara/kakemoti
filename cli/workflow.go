package cli

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/ohler55/ojg/jp"
	"github.com/olekukonko/tablewriter"
	"github.com/w-haibara/kakemoti/compiler"
	"github.com/w-haibara/kakemoti/db"
	"github.com/w-haibara/kakemoti/log"
	"github.com/w-haibara/kakemoti/worker"
)

func init() {
	registeerTypesForGob()
}

type ExecWorkflowOneceOpt struct {
	RegisterWorkflowOpt
	ExecWorkflowOpt
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

	rfn := func(name string, w compiler.Workflow, asl []byte, force bool) error {
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
	Force        bool
}

func (opt RegisterWorkflowOpt) RegisterWorkflow(ctx context.Context, coj *compiler.CtxObj) error {
	return opt.registerWorkflow(ctx, coj, db.RegisterWorkflow)
}

type registerWorkflowFunc func(name string, w compiler.Workflow, asl []byte, force bool) error

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

	aslb := asl.Bytes()
	asl = bytes.NewBuffer(aslb)

	workflow, err := compiler.Compile(ctx, asl)
	if err != nil {
		logger.Fatalln(err)
	}

	if err := fn(opt.WorkflowName, *workflow, aslb, opt.Force); err != nil {
		return err
	}

	return nil
}

type RemoveWorkflowOpt struct {
	Logfile      string
	WorkflowName string
	Force        bool
}

func (opt RemoveWorkflowOpt) RemoveWorkflow(ctx context.Context, coj *compiler.CtxObj) error {
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

	if !opt.Force && !confirm("WARNING! This will remove the workflow: "+opt.WorkflowName) {
		return nil
	}

	return db.RemoveWorkflow(opt.WorkflowName)
}

type DropWorkflowOpt struct {
	Logfile string
	Force   bool
}

func (opt DropWorkflowOpt) DropWorkflow(ctx context.Context, coj *compiler.CtxObj) error {
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

	if !opt.Force && !confirm("WARNING! This will remove all workflows") {
		return nil
	}

	return db.DropWorkflow()
}

type ExecWorkflowOpt struct {
	Logfile      string
	WorkflowName string
	Input        string
	Timeout      int
}

func (opt ExecWorkflowOpt) ExecWorkflow(ctx context.Context, coj *compiler.CtxObj) ([]byte, error) {
	return opt.execWorkflow(ctx, coj, func(name string) (compiler.Workflow, error) {
		w, err := db.GetWorkflow(name)
		if err != nil {
			return compiler.Workflow{}, err
		}

		wf, err := w.DecodeWorkflow()
		if err != nil {
			return compiler.Workflow{}, err
		}

		return wf, nil
	})
}

type getWorkflowDataFunc func(name string) (compiler.Workflow, error)

func (opt ExecWorkflowOpt) execWorkflow(ctx context.Context, coj *compiler.CtxObj, fn getWorkflowDataFunc) ([]byte, error) {
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

type ListWorkflowOpt struct {
	Writer  io.Writer
	JSON    bool
	Logfile string
}

func (opt ListWorkflowOpt) ListWorkflow() error {
	if strings.TrimSpace(opt.Logfile) == "" {
		opt.Logfile = "logs"
	}

	logger := log.NewLogger()
	close := setLogOutput(logger, opt.Logfile)
	defer func() {
		if err := close(); err != nil {
			logger.Fatalln(err)
		}
	}()

	w, err := db.ListWorkflow()
	if err != nil {
		return err
	}

	if opt.JSON {
		b, err := json.MarshalIndent(w, "", "  ")
		if err != nil {
			return err
		}

		fmt.Println(string(b))

		return nil
	}

	table := tablewriter.NewWriter(opt.Writer)
	table.SetHeader([]string{"Name", "CreatedAt"})
	for _, v := range w {
		table.Append([]string{
			v.Name,
			v.CreatedAt.Format(time.RFC3339),
		})
	}
	table.Render()

	return nil
}

type ShowWorkflowOpt struct {
	Writer       io.Writer
	WorkflowName string
	JSON         bool
	Logfile      string
}

func (opt ShowWorkflowOpt) ShowWorkflow() error {
	if strings.TrimSpace(opt.Logfile) == "" {
		opt.Logfile = "logs"
	}

	logger := log.NewLogger()
	close := setLogOutput(logger, opt.Logfile)
	defer func() {
		if err := close(); err != nil {
			logger.Fatalln(err)
		}
	}()

	w, err := db.GetWorkflow(opt.WorkflowName)
	if err != nil {
		return err
	}

	if opt.JSON {
		b, err := json.MarshalIndent(w, "", "  ")
		if err != nil {
			return err
		}

		fmt.Fprintln(opt.Writer, string(b))

		return nil
	}

	table := tablewriter.NewWriter(opt.Writer)
	table.SetHeader([]string{"Name", "CreatedAt"})
	table.Append([]string{
		w.Name,
		w.CreatedAt.Format(time.RFC3339),
	})
	table.Render()

	asl, err := w.DecodeASL()
	if err != nil {
		return err
	}
	fmt.Fprintln(opt.Writer, asl)

	return nil
}

func registeerTypesForGob() {
	gob.Register(map[string]interface{}{})
	gob.Register([]interface{}{})

	gob.Register(compiler.ChoiceState{})
	gob.Register(compiler.CommonState5{})
	gob.Register(compiler.FailState{})
	gob.Register(compiler.MapState{})
	gob.Register(compiler.ParallelState{})
	gob.Register(compiler.PassState{})
	gob.Register(compiler.SucceedState{})
	gob.Register(compiler.TaskState{})
	gob.Register(compiler.WaitState{})

	gob.Register(compiler.AndRule{})
	gob.Register(compiler.OrRule{})
	gob.Register(compiler.NotRule{})
	gob.Register(compiler.StringEqualsRule{})
	gob.Register(compiler.StringEqualsPathRule{})
	gob.Register(compiler.StringLessThanRule{})
	gob.Register(compiler.StringLessThanPathRule{})
	gob.Register(compiler.StringLessThanEqualsRule{})
	gob.Register(compiler.StringLessThanEqualsPathRule{})
	gob.Register(compiler.StringGreaterThanRule{})
	gob.Register(compiler.StringGreaterThanPathRule{})
	gob.Register(compiler.StringGreaterThanEqualsRule{})
	gob.Register(compiler.StringGreaterThanEqualsPathRule{})
	gob.Register(compiler.StringMatchesRule{})
	gob.Register(compiler.NumericEqualsRule{})
	gob.Register(compiler.NumericEqualsPathRule{})
	gob.Register(compiler.NumericLessThanRule{})
	gob.Register(compiler.NumericLessThanPathRule{})
	gob.Register(compiler.NumericLessThanEqualsRule{})
	gob.Register(compiler.NumericLessThanEqualsPathRule{})
	gob.Register(compiler.NumericGreaterThanRule{})
	gob.Register(compiler.NumericGreaterThanPathRule{})
	gob.Register(compiler.NumericGreaterThanEqualsRule{})
	gob.Register(compiler.NumericGreaterThanEqualsPathRule{})
	gob.Register(compiler.BooleanEqualsRule{})
	gob.Register(compiler.BooleanEqualsPathRule{})
	gob.Register(compiler.TimestampEqualsRule{})
	gob.Register(compiler.TimestampEqualsPathRule{})
	gob.Register(compiler.TimestampLessThanRule{})
	gob.Register(compiler.TimestampLessThanPathRule{})
	gob.Register(compiler.TimestampLessThanEqualsRule{})
	gob.Register(compiler.TimestampLessThanEqualsPathRule{})
	gob.Register(compiler.TimestampGreaterThanRule{})
	gob.Register(compiler.TimestampGreaterThanPathRule{})
	gob.Register(compiler.TimestampGreaterThanEqualsRule{})
	gob.Register(compiler.TimestampGreaterThanEqualsPathRule{})
	gob.Register(compiler.IsNullRule{})
	gob.Register(compiler.IsPresentRule{})
	gob.Register(compiler.IsNumericRule{})
	gob.Register(compiler.IsStringRule{})
	gob.Register(compiler.IsBooleanRule{})
	gob.Register(compiler.IsTimestampRule{})

	gob.Register(jp.Root(0))
	gob.Register(jp.Child(""))
}
