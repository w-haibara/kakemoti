package main

import (
	"log"
)

func main() {
	sm, err := NewStateMachine("./workflow.json")
	if err != nil {
		log.Panic("error:", err)
	}

	//sm.PrintInfo()
	//sm.PrintStates()

	sm.Start()
}
