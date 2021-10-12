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

func init() {
	pp.ColoringEnabled = false
}

type RootLogger struct {
}

func NewRootLogger() *RootLogger {
	return &RootLogger{}
}

func (l *RootLogger) Println(id string, v ...interface{}) {
	log.Println(id+":", fmt.Sprint(v...))
}

type Logger struct {
	CH                 chan Message
	LogDirPrefix       string
	RootStateMachineID string
	files              map[string]*os.File
}

type Message struct {
	StateMachineID string
	StateName      string
	StateType      string
	Body           string
	close          bool
}

func NewLogger(rootStateMachineID string) *Logger {
	ch := make(chan Message, 100)
	return &Logger{
		CH:                 ch,
		LogDirPrefix:       "logfiles",
		RootStateMachineID: rootStateMachineID,
		files:              make(map[string]*os.File),
	}
}

func (l *Logger) Listen() {
	dir := filepath.Join(l.LogDirPrefix, time.Now().Format("2006-01-02-03-04-")+l.RootStateMachineID+".log")
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Panic("mkdir failed: ", err.Error())
	}

	defer func() {
		for _, f := range l.files {
			if err := f.Close(); err != nil {
				log.Panic("file close failed: ", err.Error())
			}
		}
	}()

	for {
		v, ok := <-l.CH
		if !ok {
			break
		}

		f, ok := l.files[v.StateMachineID]
		if ok && v.close {
			if err := l.files[v.StateMachineID].Close(); err != nil {
				log.Panic("file close failed:", err.Error())
			}
			delete(l.files, v.StateMachineID)
			continue
		}
		if !ok {
			if strings.TrimSpace(v.StateMachineID) == "" {
				log.Panic("statemachine ID is brank")
			}

			name := filepath.Join(dir, v.StateMachineID)

			if _, err := os.Stat(name); err == nil { // if file exists
				f, err = os.Open(name) // #nosec G304
				if err != nil {
					log.Panic("open file error:", err.Error())
				}
			} else {
				f, err = os.Create(name)
				if err != nil {
					log.Panic("create file error:", err.Error())
				}
			}

			l.files[v.StateMachineID] = f
		}

		//TODO: format with json
		_, err := pp.Fprintln(f, v)
		if err != nil {
			log.Panic("can not write logs:", err.Error())
		}
	}
}

func (l *Logger) Close(id string) {
	l.CH <- Message{
		StateMachineID: id,
		close:          true,
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
