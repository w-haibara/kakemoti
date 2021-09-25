package main

import (
	"bytes"
	"log"
)

const input = `
{
	"key": "value"
}
`

func main() {
	sm, err := NewStateMachine("./workflow.json")
	if err != nil {
		log.Panic("error:", err)
	}

	//sm.PrintInfo()
	//sm.PrintStates()

	r := new(bytes.Buffer)
	w := new(bytes.Buffer)
	if _, err := r.WriteString(input); err != nil {
		log.Panic("error:", err)
	}

	log.Println("===  First input  ===\n", r.String())

	if err := sm.Start(r, w); err != nil {
		log.Panic("error:", err)
	}

	log.Println("=== Finaly output ===\n", w.String())
}
