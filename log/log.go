package log

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/k0kubun/pp"
)

type Logger struct {
	CH                 chan Message
	LogDirPrefix       string
	RootStateMachineID string
}

type Message struct {
	StateMachineID string
	StateName      string
	StateType      string
	Body           string
}

func init() {
	pp.ColoringEnabled = false
}

func NewLogger(rootStateMachineID string) *Logger {
	ch := make(chan Message, 100)
	return &Logger{
		CH:                 ch,
		LogDirPrefix:       "logfiles",
		RootStateMachineID: rootStateMachineID,
	}
}

func (l *Logger) Listen() {
	dir := filepath.Join(l.LogDirPrefix, time.Now().Format("2006-01-02-03-04-")+l.RootStateMachineID+".log")
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		panic("mkdir failed: " + err.Error())
	}

	files := make(map[string]*os.File)
	defer func() {
		for _, f := range files {
			if err := f.Close(); err != nil {
				panic("file close failed: " + err.Error())
			}
		}
	}()

	for {
		v, ok := <-l.CH
		if !ok {
			break
		}

		f, ok := files[v.StateMachineID]
		if !ok {
			if strings.TrimSpace(v.StateMachineID) == "" {
				panic("statemachine ID is brank")
			}

			var err error
			f, err = os.Create(filepath.Join(dir, v.StateMachineID))
			if err != nil {
				panic("create file error: " + err.Error())
			}

			files[v.StateMachineID] = f
		}

		//TODO: format with json
		_, err := pp.Fprintln(f, v)
		if err != nil {
			log.Panic("can not write logs: ", err.Error())
		}
	}
}

func (l *Logger) Println(id, name, typ string, v ...interface{}) {
	l.CH <- Message{
		StateMachineID: id,
		StateName:      name,
		StateType:      typ,
		Body:           fmt.Sprint(v...),
	}
}

func Println(id string, v ...interface{}) {
	log.Println(id+":", fmt.Sprint(v...))
}
