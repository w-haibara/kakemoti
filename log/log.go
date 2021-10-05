package log

import (
	"fmt"
	"log"
)

type Logger struct {
	CH  chan Message
	Que []Message
}

type Message struct {
	StateMachineID string
	StateName      string
	StateType      string
	Body           string
}

func NewLogger() *Logger {
	ch := make(chan Message, 100)
	return &Logger{
		CH:  ch,
		Que: []Message{},
	}
}

func (l *Logger) Listen() {
	for {
		v, ok := <-l.CH
		if !ok {
			break
		}

		l.Que = append(l.Que, v)
	}
}

func (l *Logger) FLush() {

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
