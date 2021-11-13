package cli

import (
	"bytes"
	"context"
	"os"
	"strings"

	"karage/log"
	"karage/statemachine"
)

type Options struct {
	Input   string
	ASL     string
	Timeout int64
}

func StartExecution(ctx context.Context, opt Options) ([]byte, error) {
	logger := log.NewLogger()
	logger.SetWriter()

	if strings.TrimSpace(opt.Input) == "" {
		logger.Fatalln("input option value is empty")
	}

	if strings.TrimSpace(opt.ASL) == "" {
		logger.Fatalln("ASL option value is empty")
	}

	f1, input, err := readFile(opt.Input)
	if err != nil {
		logger.Fatalln(err)
	}
	defer func() {
		if err := f1.Close(); err != nil {
			logger.Fatalln(err)
		}
	}()

	f2, asl, err := readFile(opt.ASL)
	if err != nil {
		logger.Fatalln(err)
	}
	defer func() {
		if err := f2.Close(); err != nil {
			logger.Fatalln(err)
		}
	}()

	return statemachine.Start(ctx, asl, input, opt.Timeout, logger)
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
