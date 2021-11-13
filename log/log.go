package log

import (
	"io"
	"os"
	"path/filepath"
	"time"

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

func (l *Logger) SetWriter() (close func()) {
	if _, err := os.Stat("logs"); err != nil {
		if err := os.Mkdir("logs", os.ModePerm); err != nil {
			l.Fatal(err)
		}
	}

	f, err := os.Create(filepath.Join("logs", time.Now().Format("2006010215040507")+".log"))
	if err != nil {
		l.Fatal(err)
	}

	w := io.MultiWriter(os.Stderr, f)
	l.Logger.SetOutput(w)

	return func() {
		if err := f.Close(); err != nil {
			logrus.Fatal(err)
		}
	}
}
