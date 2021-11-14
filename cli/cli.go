package cli

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"karage/log"
	"karage/statemachine"
)

type Options struct {
	Logfile string
	Input   string
	ASL     string
	Timeout int64
}

func StartExecution(ctx context.Context, opt Options) ([]byte, error) {
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

func setLogOutput(l *log.Logger, path string) (close func() error) {
	if _, err := os.Stat(path); err != nil {
		if err := os.Mkdir(path, os.ModePerm); err != nil {
			l.Fatal(err)
		}
	}

	f, err := os.Create(filepath.Join(path, time.Now().Format("2006010215040507")+".log"))
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