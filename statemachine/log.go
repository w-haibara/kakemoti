package statemachine

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

func SetLogWriter() (close func()) {
	if err := os.Mkdir("logs", os.ModePerm); err != nil {
		logrus.Fatal(err)
	}

	f, err := os.Create(filepath.Join("logs", time.Now().Format("2006010215040507")+".log"))
	if err != nil {
		logrus.Fatal(err)
	}

	w := io.MultiWriter(os.Stderr, f)
	logrus.SetOutput(w)

	return func() {
		if err := f.Close(); err != nil {
			logrus.Fatal(err)
		}
	}
}
