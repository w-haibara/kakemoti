package log

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Entry
}

func NewLogger() *Logger {
	l := logrus.NewEntry(logrus.New())

	l.Logger.SetLevel(logrus.DebugLevel)
	l.Logger.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint: true,
	})

	return &Logger{l}
}

func Line() string {
	_, path, line, ok := runtime.Caller(2)
	if !ok {
		return "---"
	}

	_, file := filepath.Split(path)

	return fmt.Sprintf("%s:%d", file, line)
}
