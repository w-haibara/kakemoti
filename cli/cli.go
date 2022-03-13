package cli

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
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

func confirm(msg string) bool {
	fmt.Print(msg + "\nAre you sure you want to continue? [y/N]: ")

	str, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		panic(err.Error())
	}

	if strings.TrimSpace(strings.TrimSuffix(str, "\n")) == "y" {
		return true
	}

	return false
}
