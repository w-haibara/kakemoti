package main

import (
	"fmt"
	"log"
	"strings"
)

type Choice map[string]interface{}

func (c Choice) next() (string, bool) {
	v, ok := c["Next"].(string)
	if !ok || strings.TrimSpace(v) == "" {
		return "", false
	}

	return v, true
}

func (c Choice) variable() (string, bool) {
	v, ok := c["Variable"].(string)
	if !ok || strings.TrimSpace(v) == "" {
		return "", false
	}

	return v, true
}

type ChoiceState struct {
	CommonState
	Choices []Choice `json:"Choices"`
	Default string   `json:"Default"`
}

func (s ChoiceState) Transition() (next string, err error) {
	exist := func(m map[string]interface{}, k string) bool {
		_, ok := m[k]
		return ok
	}

	for _, choice := range s.Choices {
		switch {
		case exist(choice, "And"):
		case exist(choice, "BooleanEquals"):
			//TODO: tempolary implement
			log.Println("Choice: BooleanEquals")

			v, ok := choice.variable()
			if !ok {
				return "", fmt.Errorf("choice rule error: Variable is blank")
			}

			n, ok := choice.next()
			if !ok {
				return "", fmt.Errorf("choice rule error: Next is blank")
			}

			log.Println("	Variable", v)
			log.Println("	Next", n)

			return n, nil

		case exist(choice, "BooleanEqualsPath"):
		case exist(choice, "IsBoolean"):
		case exist(choice, "IsNull"):
		case exist(choice, "IsNumeric"):
		case exist(choice, "IsPresent"):
		case exist(choice, "IsString"):
		case exist(choice, "IsTimestamp"):
		case exist(choice, "Not"):
		case exist(choice, "NumericEquals"):
		case exist(choice, "NumericEqualsPath"):
		case exist(choice, "NumericGreaterThan"):
		case exist(choice, "NumericGreaterThanPath"):
		case exist(choice, "NumericGreaterThanEquals"):
		case exist(choice, "NumericGreaterThanEqualsPath"):
		case exist(choice, "NumericLessThan"):
		case exist(choice, "NumericLessThanPath"):
		case exist(choice, "NumericLessThanEquals"):
		case exist(choice, "NumericLessThanEqualsPath"):
		case exist(choice, "Or"):
		case exist(choice, "StringEquals"):
		case exist(choice, "StringEqualsPath"):
		case exist(choice, "StringGreaterThan"):
		case exist(choice, "StringGreaterThanPath"):
		case exist(choice, "StringGreaterThanEquals"):
		case exist(choice, "StringGreaterThanEqualsPath"):
		case exist(choice, "StringLessThanStringLessThanPath"):
		case exist(choice, "StringLessThanEqualsStringLessThanEqualsPath"):
		case exist(choice, "StringMatches"):
		case exist(choice, "TimestampEquals"):
		case exist(choice, "TimestampEqualsPath"):
		case exist(choice, "TimestampGreaterThan"):
		case exist(choice, "TimestampGreaterThanPath"):
		case exist(choice, "TimestampGreaterThanEquals"):
		case exist(choice, "TimestampGreaterThanEqualsPath"):
		case exist(choice, "TimestampLessThan"):
		case exist(choice, "TimestampLessThanPath"):
		case exist(choice, "TimestampLessThanEquals"):
		case exist(choice, "TimestampLessThanEqualsPath"):
		}

		println()
	}

	return "", EndStateMachine
}
