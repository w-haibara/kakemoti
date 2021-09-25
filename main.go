package main

import (
	"bytes"
	"log"
	"os"
)

const (
	input1 = `{
	"IsHelloWorldExample": true
}`

	input2 = `{
	"IsHelloWorldExample": false
}`
)

func main() {
	sm, err := NewStateMachine("./workflow.asl.json", log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile))
	if err != nil {
		log.Panic("error:", err)
	}

	//sm.PrintInfo()
	//sm.PrintStates()

	r := new(bytes.Buffer)
	w := new(bytes.Buffer)
	if _, err := r.WriteString(input1); err != nil {
		log.Panic("error:", err)
	}

	log.Println("===  First input  ===", "\n"+r.String())

	if err := sm.Start(r, w); err != nil {
		log.Panic("error:", err)
	}

	log.Println("=== Finaly output ===", "\n"+w.String())
}
