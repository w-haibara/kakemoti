package main

import ()

type ChoiceState struct {
	CommonState
	Choices []map[string]interface{} `json:"Choices"`
	Default string                   `json:"Default"`
}
