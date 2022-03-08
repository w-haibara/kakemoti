package cli

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/w-haibara/kakemoti/log"
)

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
